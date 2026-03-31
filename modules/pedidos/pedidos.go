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
	type PedidoOnline struct {
		ID         int
		Cliente    string
		Total      int
		Fecha      string
		TipoPedido string
	}
	type PedidoAbierto struct {
		ID         int
		Cliente    string
		Total      int
		Fecha      string
		TipoPedido string
		Mesa       string
	}
	type InicioData struct {
		TurnoActivo    *TurnoActivo
		ProductosTurno []ProductoTurno
		Mesas          []struct {
			ID      int
			Nombre  string
			Ocupada bool
		}
		PedidosOnline   []PedidoOnline
		PedidosAbiertos bool
		PedidosLocales  []PedidoAbierto
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

	rowsOnline, err := bases.DB.Query(`
    SELECT id_online, cliente, total, DATE_FORMAT(fecha, '%H:%i'), tipo_pedido
    FROM pedidos_online 
    WHERE estado = 'pendiente'
    ORDER BY fecha ASC
`)
	if err == nil {
		defer rowsOnline.Close()
		for rowsOnline.Next() {
			var p PedidoOnline
			rowsOnline.Scan(&p.ID, &p.Cliente, &p.Total, &p.Fecha, &p.TipoPedido)
			data.PedidosOnline = append(data.PedidosOnline, p)
		}
	}
	var estadoPedidos string
	bases.DB.QueryRow("SELECT valor FROM configuracion WHERE clave='pedidos_online'").Scan(&estadoPedidos)
	data.PedidosAbiertos = estadoPedidos == "abierto"
	rowsLocales, err := bases.DB.Query(`
    SELECT p.id_ped, p.cliente, p.total, DATE_FORMAT(p.fecha, '%H:%i'), p.tipo_pedido, COALESCE(m.nombre, '-')
    FROM pedidos p
    LEFT JOIN mesas m ON p.id_mesa = m.id_mesa
    WHERE p.estado = 'abierto'
    ORDER BY p.fecha ASC
`)
	if err == nil {
		defer rowsLocales.Close()
		for rowsLocales.Next() {
			var p PedidoAbierto
			rowsLocales.Scan(&p.ID, &p.Cliente, &p.Total, &p.Fecha, &p.TipoPedido, &p.Mesa)
			data.PedidosLocales = append(data.PedidosLocales, p)
		}
	}

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
			"INSERT INTO pedidos (fecha, total, cliente, tipo_pedido, id_mesa, estado) VALUES (NOW(), 0, ?, ?, ?, 'abierto')",
			cliente, tipoPedido, idMesa,
		)
	} else {
		result, err = bases.DB.Exec(
			"INSERT INTO pedidos (fecha, total, cliente, tipo_pedido, estado) VALUES (NOW(), 0, ?, ?, 'abierto')",
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

	ctx.Redirect("/inicio", 302)
}
func CerrarPedido(ctx *fasthttp.RequestCtx) {
	idPedido, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil {
		ctx.Error("ID inválido", 400)
		return
	}

	var idMesa sql.NullInt64
	bases.DB.QueryRow("SELECT id_mesa FROM pedidos WHERE id_ped = ?", idPedido).Scan(&idMesa)
	bases.DB.Exec("UPDATE pedidos SET estado = 'cerrado' WHERE id_ped = ?", idPedido)

	if idMesa.Valid {
		bases.DB.Exec("UPDATE mesas SET ocupada = 0 WHERE id_mesa = ?", idMesa.Int64)
	}

	log.Printf(">>> Llamando DescontarStock para pedido %d\n", idPedido)
	inventario.DescontarStock(int64(idPedido))
	log.Printf(">>> DescontarStock finalizado para pedido %d\n", idPedido)

	ctx.Redirect("/inicio", 302)
}
func MarcarListo(ctx *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	bases.DB.Exec("UPDATE pedidos_online SET estado='listo' WHERE id_online=?", id)

	// ✅ AGREGAR ESTA LÍNEA:
	inventario.DescontarStockOnline(int64(id))

	ctx.Redirect("/inicio", 302)
}
func DetalleOnline(ctx *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	log.Printf(">>> DetalleOnline id: %d\n", id)
	type ItemDetalle struct {
		Nombre        string
		Cantidad      int
		Precio        int
		Modificadores []string
	}
	type DetalleData struct {
		ID         int
		Cliente    string
		Total      int
		Fecha      string
		TipoPedido string
		Notas      string
		Telefono   string
		Direccion  string
		Items      []ItemDetalle
	}

	var data DetalleData
	data.ID = id

	// 1. Traer lo básico (Nombre, Total, Fecha, Tipo, Notas)
	// Quitamos el pedido_json para que no de error si no existe la columna
	err := bases.DB.QueryRow(`
    SELECT cliente, telefono, total, DATE_FORMAT(fecha, '%d/%m %Y %H:%i'), tipo_pedido, IFNULL(notas, ''), IFNULL(direccion, '')
    FROM pedidos_online WHERE id_online = ?`, id).Scan(
		&data.Cliente, &data.Telefono, &data.Total, &data.Fecha, &data.TipoPedido, &data.Notas, &data.Direccion,
	)

	if err != nil {
		log.Println("Error en consulta principal:", err)
	}

	// 2. Traer productos
	rows, err := bases.DB.Query(`
    SELECT pr.nombre, pod.cantidad, pod.precio, pod.id_ped, IFNULL(pod.notas_producto, '')
    FROM pedidos_online_detalle pod
    JOIN productos pr ON pod.id_pro = pr.id_pro
    WHERE pod.id_online = ?
`, id)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var item ItemDetalle
			var idDetalle int
			var notaProd string
			rows.Scan(&item.Nombre, &item.Cantidad, &item.Precio, &idDetalle, &notaProd)

			// 👇 Buscar mods por id_detalle, no por id_pro
			rowsMod, _ := bases.DB.Query(`
            SELECT m.nombre 
            FROM pedidos_online_modificadores pom
            JOIN modificadores m ON pom.id_mod = m.id_mod
            WHERE pom.id_detalle = ?
        `, idDetalle)

			if rowsMod != nil {
				for rowsMod.Next() {
					var mod string
					if errM := rowsMod.Scan(&mod); errM == nil && mod != "" {
						item.Modificadores = append(item.Modificadores, mod)
					}
				}
				rowsMod.Close()
			}

			// Nota manual una sola vez al final
			if notaProd != "" {
				item.Modificadores = append(item.Modificadores, notaProd)
			}

			data.Items = append(data.Items, item)
		}
	}

	tmpl, err := template.ParseFiles("templates/online_detalle.html")
	if err != nil {
		ctx.Error("Error template", 500)
		return
	}
	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, data)
}
func ImprimirOnline(ctx *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))

	var cliente string
	bases.DB.QueryRow("SELECT cliente FROM pedidos_online WHERE id_online = ?", id).Scan(&cliente)

	rows, err := bases.DB.Query(`
        SELECT pr.nombre, pod.cantidad, pod.id_ped, IFNULL(pod.notas_producto, '')
        FROM pedidos_online_detalle pod
        JOIN productos pr ON pod.id_pro = pr.id_pro
        WHERE pod.id_online = ?
    `, id)

	var items []impresora.ItemTicket
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var nombre, notaProd string
			var cantidad, idDetalle int
			rows.Scan(&nombre, &cantidad, &idDetalle, &notaProd)

			// Traer modificadores por id_detalle
			var mods []string
			rowsMod, _ := bases.DB.Query(`
                SELECT m.nombre
                FROM pedidos_online_modificadores pom
                JOIN modificadores m ON pom.id_mod = m.id_mod
                WHERE pom.id_detalle = ?
            `, idDetalle)
			if rowsMod != nil {
				for rowsMod.Next() {
					var mod string
					if errM := rowsMod.Scan(&mod); errM == nil && mod != "" {
						mods = append(mods, mod)
					}
				}
				rowsMod.Close()
			}
			if notaProd != "" {
				mods = append(mods, notaProd)
			}

			items = append(items, impresora.ItemTicket{
				Cantidad:      cantidad,
				Nombre:        nombre,
				Modificadores: mods,
			})
		}
	}

	err = impresora.ImprimirTicket(id, cliente, items)
	if err != nil {
		log.Println("Advertencia - impresora no disponible:", err)
	}
	ctx.Redirect("/online/detalle/"+strconv.Itoa(id), 302)
}
func Landing(ctx *fasthttp.RequestCtx) {
	tmpl, err := template.ParseFiles("templates/landing.html")
	if err != nil {
		log.Println("Error al cargar template:", err)
		ctx.Error("Error al cargar template", 500)
		return
	}
	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, nil)
}
func AgregarAPedido(ctx *fasthttp.RequestCtx) {
	idPedido, _ := strconv.Atoi(ctx.UserValue("id").(string))
	args := ctx.Request.PostArgs()
	pedidoJSON := string(args.Peek("pedido_json"))
	var cliente string
	bases.DB.QueryRow("SELECT cliente FROM pedidos WHERE id_ped = ?", idPedido).Scan(&cliente)

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
		ctx.Error("Error al procesar pedido", 500)
		return
	}

	type ItemAgrupado struct {
		Nombre   string
		Cantidad int
		Mods     []string
		IDPro    int
		Precio   int
	}
	agrupados := map[string]*ItemAgrupado{}
	var ordenAgrupado []string

	for _, item := range itemsJSON {
		var precio int
		bases.DB.QueryRow("SELECT precio FROM productos WHERE id_pro = ?", item.IDPro).Scan(&precio)

		var modNombres []string
		for _, mod := range item.Mods {
			modNombres = append(modNombres, mod.Nombre)
		}

		modsKey := item.Nombre
		for _, m := range modNombres {
			modsKey += "_" + m
		}

		if _, existe := agrupados[modsKey]; !existe {
			agrupados[modsKey] = &ItemAgrupado{
				Nombre:   item.Nombre,
				Cantidad: 0,
				Mods:     modNombres,
				IDPro:    item.IDPro,
				Precio:   precio,
			}
			ordenAgrupado = append(ordenAgrupado, modsKey)
		}
		agrupados[modsKey].Cantidad++
	}

	// Insertar agrupado en BD y modificadores
	for _, key := range ordenAgrupado {
		a := agrupados[key]

		bases.DB.Exec(
			"INSERT INTO pedidos_detalle (id_ped, id_pro, cantidad, precio) VALUES (?, ?, ?, ?)",
			idPedido, a.IDPro, a.Cantidad, a.Precio,
		)
		bases.DB.Exec(
			"UPDATE pedidos SET total = total + ? WHERE id_ped = ?",
			a.Precio*a.Cantidad, idPedido,
		)

		// Insertar modificadores una sola vez por grupo
		for _, item := range itemsJSON {
			modsKeyItem := item.Nombre
			for _, m := range item.Mods {
				modsKeyItem += "_" + m.Nombre
			}
			if modsKeyItem == key {
				for _, mod := range item.Mods {
					idMod, _ := strconv.Atoi(mod.ID)
					bases.DB.Exec(
						"INSERT INTO pedidos_modificadores (id_ped, id_pro, id_mod) VALUES (?, ?, ?)",
						idPedido, a.IDPro, idMod,
					)
				}
				break
			}
		}
	}

	// Imprimir ticket agrupado
	var items []impresora.ItemTicket
	for _, key := range ordenAgrupado {
		a := agrupados[key]
		items = append(items, impresora.ItemTicket{
			Cantidad:      a.Cantidad,
			Nombre:        a.Nombre,
			Modificadores: a.Mods,
		})
	}
	err := impresora.ImprimirTicket(idPedido, cliente, items)
	if err != nil {
		log.Println("Advertencia - impresora no disponible:", err)
	}

	ctx.Redirect("/inicio", 302)
}
func EditarPedido(ctx *fasthttp.RequestCtx) {
	idPedido, _ := strconv.Atoi(ctx.UserValue("id").(string))

	// Cargar datos del pedido
	type ItemExistente struct {
		Nombre   string
		Cantidad int
		Mods     []string
	}
	type EditarData struct {
		IDPedido   int
		Cliente    string
		Categorias []Categoria
		Items      []ItemExistente
	}

	var data EditarData
	data.IDPedido = idPedido
	bases.DB.QueryRow("SELECT cliente FROM pedidos WHERE id_ped = ?", idPedido).Scan(&data.Cliente)

	// Cargar items existentes del pedido
	rows, _ := bases.DB.Query(`
        SELECT pr.nombre, pd.cantidad
        FROM pedidos_detalle pd
        JOIN productos pr ON pd.id_pro = pr.id_pro
        WHERE pd.id_ped = ?
    `, idPedido)
	defer rows.Close()
	for rows.Next() {
		var item ItemExistente
		rows.Scan(&item.Nombre, &item.Cantidad)
		data.Items = append(data.Items, item)
	}

	// Cargar categorias y productos (igual que NuevoPedido)
	rowsPro, err := bases.DB.Query(`
        SELECT p.id_pro, p.nombre, p.precio, c.nombre 
        FROM categorias c
        LEFT JOIN productos p ON p.id_cat = c.id_cat
        ORDER BY c.nombre, p.nombre
    `)
	if err == nil {
		defer rowsPro.Close()
		categoriasMap := map[string]*Categoria{}
		var orden []string
		productosMap := map[int]*Producto{}
		for rowsPro.Next() {
			var idPro, precio int
			var nombrePro, nombreCat string
			rowsPro.Scan(&idPro, &nombrePro, &precio, &nombreCat)
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
		for _, nombre := range orden {
			data.Categorias = append(data.Categorias, *categoriasMap[nombre])
		}
	}

	tmpl, err := template.ParseFiles("templates/pedido_editar.html")
	if err != nil {
		log.Println("Error al cargar template:", err)
		ctx.Error("Error al cargar template", 500)
		return
	}
	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, data)
}
func AbrirPedidos(ctx *fasthttp.RequestCtx) {
	bases.DB.Exec("UPDATE configuracion SET valor='abierto' WHERE clave='pedidos_online'")
	ctx.Redirect("/inicio", 302)
}

func CerrarPedidos(ctx *fasthttp.RequestCtx) {
	bases.DB.Exec("UPDATE configuracion SET valor='cerrado' WHERE clave='pedidos_online'")
	ctx.Redirect("/inicio", 302)
}
