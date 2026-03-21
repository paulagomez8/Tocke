package pedidos

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"strconv"
	"tockesanfelipe/modules/bases"
	"tockesanfelipe/modules/impresora"

	"tockesanfelipe/modules/inventario"

	"github.com/valyala/fasthttp"
)

type Producto struct {
	ID            int
	Nombre        string
	Precio        int
	Modificadores []Modificador
}

type Categoria struct {
	Nombre    string
	Productos []Producto
}

type Modificador struct {
	ID     int
	Nombre string
}

type Mesa struct {
	ID      int
	Nombre  string
	Ocupada bool
}

func Inicio(ctx *fasthttp.RequestCtx) {
	type TurnoActivo struct {
		ID     int
		Nombre string
		Inicio string
	}
	type ProductoTurno struct {
		Nombre   string
		Cantidad int
	}
	type InicioData struct {
		TurnoActivo    *TurnoActivo
		ProductosTurno []ProductoTurno
		Mesas          []struct {
			ID      int
			Nombre  string
			Ocupada bool
		}
	}

	var data InicioData
	row := bases.DB.QueryRow("SELECT id_turno, nombre, DATE_FORMAT(inicio, '%d/%m %H:%i'), inicio FROM turnos WHERE fin IS NULL ORDER BY id_turno DESC LIMIT 1")
	var id int
	var nombre, inicioFormateado, inicioRaw string
	if err := row.Scan(&id, &nombre, &inicioFormateado, &inicioRaw); err != nil {
		log.Println("Error scan turno:", err)
	} else {
		data.TurnoActivo = &TurnoActivo{id, nombre, inicioFormateado}

		rows, err := bases.DB.Query(`
    SELECT pr.nombre, SUM(pd.cantidad) as cantidad
    FROM pedidos p
    JOIN pedidos_detalle pd ON p.id_ped = pd.id_ped
    JOIN productos pr ON pd.id_pro = pr.id_pro
    WHERE p.fecha >= ?
    GROUP BY pr.id_pro, pr.nombre
    ORDER BY cantidad DESC
`, inicioRaw)
		if err != nil {
			log.Println("Error query productos turno:", err)
		} else {
			defer rows.Close()
			for rows.Next() {
				var p ProductoTurno
				rows.Scan(&p.Nombre, &p.Cantidad)
				data.ProductosTurno = append(data.ProductosTurno, p)
			}
			log.Printf("Turno ID: %d inicio: %s productos: %d\n", id, inicioRaw, len(data.ProductosTurno))
		}
	}

	rowsMesas, err := bases.DB.Query("SELECT id_mesa, nombre, ocupada FROM mesas ORDER BY id_mesa")
	if err == nil {
		defer rowsMesas.Close()
		for rowsMesas.Next() {
			var m struct {
				ID      int
				Nombre  string
				Ocupada bool
			}
			rowsMesas.Scan(&m.ID, &m.Nombre, &m.Ocupada)
			data.Mesas = append(data.Mesas, m)
		}
	}

	tmpl, err := template.ParseFiles("templates/inicio.html")
	if err != nil {
		log.Println("Error al cargar template:", err)
		ctx.Error("Error al cargar template", 500)
		return
	}
	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, data)

}

func LiberarMesa(ctx *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	bases.DB.Exec("UPDATE mesas SET ocupada=0 WHERE id_mesa=?", id)
	ctx.Redirect("/", 302)
}

func NuevoPedido(ctx *fasthttp.RequestCtx) {
	rows, err := bases.DB.Query(`
        SELECT p.id_pro, p.nombre, p.precio, c.nombre 
        FROM categorias c
        LEFT JOIN productos p ON p.id_cat = c.id_cat
        ORDER BY c.nombre, p.nombre
    `)
	if err != nil {
		log.Println("Error al obtener productos:", err)
		ctx.Error("Error al obtener productos", 500)
		return
	}
	defer rows.Close()

	categoriasMap := map[string]*Categoria{}
	var orden []string
	productosMap := map[int]*Producto{}

	for rows.Next() {
		var idPro, precio int
		var nombrePro, nombreCat string
		rows.Scan(&idPro, &nombrePro, &precio, &nombreCat)

		if _, existe := categoriasMap[nombreCat]; !existe {
			categoriasMap[nombreCat] = &Categoria{Nombre: nombreCat}
			orden = append(orden, nombreCat)
		}

		if idPro == 0 {
			continue
		}

		p := Producto{ID: idPro, Nombre: nombrePro, Precio: precio}
		productosMap[idPro] = &p
		categoriasMap[nombreCat].Productos = append(categoriasMap[nombreCat].Productos, p)
	}

	rowsMod, err := bases.DB.Query("SELECT id_mod, id_pro, nombre FROM modificadores")
	if err == nil {
		defer rowsMod.Close()
		for rowsMod.Next() {
			var idMod, idPro int
			var nombre string
			rowsMod.Scan(&idMod, &idPro, &nombre)
			if p, existe := productosMap[idPro]; existe {
				p.Modificadores = append(p.Modificadores, Modificador{ID: idMod, Nombre: nombre})
			}
		}
	}

	for _, cat := range categoriasMap {
		for i, p := range cat.Productos {
			if prod, existe := productosMap[p.ID]; existe {
				cat.Productos[i].Modificadores = prod.Modificadores
			}
		}
	}

	var categorias []Categoria
	for _, nombre := range orden {
		categorias = append(categorias, *categoriasMap[nombre])
	}
	var mesas []Mesa
	rowsMesas, err := bases.DB.Query("SELECT id_mesa, nombre, ocupada FROM mesas ORDER BY id_mesa")
	if err == nil {
		defer rowsMesas.Close()
		for rowsMesas.Next() {
			var m Mesa
			rowsMesas.Scan(&m.ID, &m.Nombre, &m.Ocupada)
			mesas = append(mesas, m)
		}
	}
	type NuevoPedidoData struct {
		Categorias []Categoria
		Mesas      []Mesa
	}

	tmpl, err := template.ParseFiles("templates/pedido.html")
	if err != nil {
		log.Println("Error al cargar template:", err)
		ctx.Error("Error al cargar template", 500)
		return
	}
	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, NuevoPedidoData{Categorias: categorias, Mesas: mesas})

}

func ConfirmarPedido(ctx *fasthttp.RequestCtx) {
	args := ctx.Request.PostArgs()
	cliente := string(args.Peek("cliente"))
	tipoPedido := string(args.Peek("tipo_pedido"))
	pedidoJSON := string(args.Peek("pedido_json"))
	idMesa, _ := strconv.Atoi(string(args.Peek("id_mesa")))

	type ItemJSON struct {
		IDPro  int    `json:"idPro"`
		Nombre string `json:"nombre"`
		Mods   []struct {
			ID     string `json:"id"`
			Nombre string `json:"nombre"`
		} `json:"mods"`
	}

	var itemsJSON []ItemJSON
	if err := json.Unmarshal([]byte(pedidoJSON), &itemsJSON); err != nil {
		log.Println("Error al parsear pedido:", err)
		ctx.Error("Error al procesar pedido", 500)
		return
	}

	var result sql.Result
	var err error
	if idMesa > 0 {
		result, err = bases.DB.Exec(
			"INSERT INTO pedidos (fecha, total, cliente, tipo_pedido, id_mesa) VALUES (NOW(), 0, ?, ?, ?)",
			cliente, tipoPedido, idMesa,
		)
	} else {
		result, err = bases.DB.Exec(
			"INSERT INTO pedidos (fecha, total, cliente, tipo_pedido) VALUES (NOW(), 0, ?, ?)",
			cliente, tipoPedido,
		)
	}
	if err != nil {
		log.Println("Error al crear pedido:", err)
		ctx.Error("Error al crear pedido", 500)
		return
	}
	idPedido, _ := result.LastInsertId()
	log.Printf("Pedido creado: %d cliente: %s items: %d\n", idPedido, cliente, len(itemsJSON))
	total := 0

	type ItemAgrupado struct {
		Nombre   string
		Cantidad int
		Mods     []string
	}
	agrupados := map[string]*ItemAgrupado{}
	var ordenAgrupado []string

	for _, item := range itemsJSON {
		var precio int
		bases.DB.QueryRow("SELECT precio FROM productos WHERE id_pro = ?", item.IDPro).Scan(&precio)
		bases.DB.Exec(
			"INSERT INTO pedidos_detalle (id_ped, id_pro, cantidad, precio) VALUES (?, ?, 1, ?)",
			idPedido, item.IDPro, precio,
		)
		total += precio

		var modNombres []string
		for _, mod := range item.Mods {
			idMod, _ := strconv.Atoi(mod.ID)
			bases.DB.Exec(
				"INSERT INTO pedidos_modificadores (id_ped, id_pro, id_mod) VALUES (?, ?, ?)",
				idPedido, item.IDPro, idMod,
			)
			modNombres = append(modNombres, mod.Nombre)
		}

		modsKey := item.Nombre
		for _, m := range modNombres {
			modsKey += "_" + m
		}

		if _, existe := agrupados[modsKey]; !existe {
			agrupados[modsKey] = &ItemAgrupado{Nombre: item.Nombre, Cantidad: 0, Mods: modNombres}
			ordenAgrupado = append(ordenAgrupado, modsKey)
		}
		agrupados[modsKey].Cantidad++
	}

	bases.DB.Exec("UPDATE pedidos SET total = ? WHERE id_ped = ?", total, idPedido)

	if idMesa > 0 {
		bases.DB.Exec("UPDATE mesas SET ocupada=1 WHERE id_mesa=?", idMesa)
	}

	var items []impresora.ItemTicket
	for _, key := range ordenAgrupado {
		a := agrupados[key]
		items = append(items, impresora.ItemTicket{
			Cantidad:      a.Cantidad,
			Nombre:        a.Nombre,
			Modificadores: a.Mods,
		})
	}

	err = impresora.ImprimirTicket(int(idPedido), cliente, items)
	if err != nil {
		log.Println("Advertencia - impresora no disponible:", err)
	}

	inventario.DescontarStock(idPedido)
	ctx.Redirect("/", 302)
}
