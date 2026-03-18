package pedidos

import (
	"encoding/json"
	"html/template"
	"log"
	"strconv"
	"time"
	"tockesanfelipe/modules/bases"
	"tockesanfelipe/modules/impresora"

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

func Inicio(ctx *fasthttp.RequestCtx) {
	tmpl, err := template.ParseFiles("templates/inicio.html")
	if err != nil {
		log.Println("Error al cargar template:", err)
		ctx.Error("Error al cargar template", 500)
		return
	}
	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, nil)
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

	tmpl, err := template.ParseFiles("templates/pedido.html")
	if err != nil {
		log.Println("Error al cargar template:", err)
		ctx.Error("Error al cargar template", 500)
		return
	}
	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, categorias)
}

func ConfirmarPedido(ctx *fasthttp.RequestCtx) {
	args := ctx.Request.PostArgs()
	cliente := string(args.Peek("cliente"))
	pedidoJSON := string(args.Peek("pedido_json"))

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

	result, err := bases.DB.Exec(
		"INSERT INTO pedidos (fecha, total, cliente) VALUES (?, 0, ?)",
		time.Now(), cliente,
	)
	if err != nil {
		log.Println("Error al crear pedido:", err)
		ctx.Error("Error al crear pedido", 500)
		return
	}

	idPedido, _ := result.LastInsertId()
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

	type TicketData struct {
		IDPedido int64
		Fecha    string
		Cliente  string
		Items    []impresora.ItemTicket
	}

	data := TicketData{
		IDPedido: idPedido,
		Fecha:    time.Now().Format("02/01/2006 15:04:05"),
		Cliente:  cliente,
		Items:    items,
	}

	tmpl, err := template.ParseFiles("templates/ticket.html")
	if err != nil {
		log.Println("Error al cargar template:", err)
		ctx.Redirect("/", 302)
		return
	}
	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, data)
}
