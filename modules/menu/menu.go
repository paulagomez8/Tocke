package menu

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"strconv"
	"time"
	"tockesanfelipe/modules/bases"

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

func VerMenu(ctx *fasthttp.RequestCtx) {
	var estado string
	bases.DB.QueryRow("SELECT valor FROM configuracion WHERE clave='pedidos_online'").Scan(&estado)
	if estado == "cerrado" {
		ctx.SetContentType("text/html")
		ctx.WriteString(`<!DOCTYPE html>
<html lang="es">
<head><meta charset="UTF-8"><title>Al Tocke</title>
<style>body{font-family:Arial,sans-serif;display:flex;justify-content:center;align-items:center;height:100vh;margin:0;background:#f4f4f4;text-align:center;}
.box{background:white;padding:40px;border-radius:8px;box-shadow:0 2px 10px rgba(0,0,0,0.1);}
h2{color:#cc0000;margin-bottom:15px;}p{color:#666;font-size:16px;}</style>
</head>
<body><div class="box">
<h2>😴 No estamos recibiendo pedidos</h2>
<p>Estamos fuera de horario o temporalmente cerrados.</p>
<p>Vuelve pronto!</p>
</div></body></html>`)
		return
	}
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

	tmpl, err := template.ParseFiles("templates/menu.html")
	if err != nil {
		log.Println("Error al cargar template:", err)
		ctx.Error("Error al cargar template", 500)
		return
	}
	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, categorias)
}

func ConfirmarPedidoOnline(ctx *fasthttp.RequestCtx) {
	args := ctx.Request.PostArgs()
	cliente := string(args.Peek("cliente"))
	tipoPedido := string(args.Peek("tipo_pedido"))
	pedidoJSON := string(args.Peek("pedido_json"))
	notasGenerales := string(args.Peek("notas"))

	type Mod struct {
		ID     string `json:"id"`
		Nombre string `json:"nombre"`
	}

	type ItemJSON struct {
		IDPro  int    `json:"idPro"`
		Nombre string `json:"nombre"`
		Mods   []Mod  `json:"mods"`
	}

	type ItemAgrupado struct {
		IDPro    int
		Cantidad int
		Mods     []Mod
	}

	// Parsear el pedido JSON
	var itemsJSON []ItemJSON
	if err := json.Unmarshal([]byte(pedidoJSON), &itemsJSON); err != nil {
		log.Println("Error al parsear pedido online:", err)
		ctx.Error("Error al procesar pedido", 500)
		return
	}

	// Obtener el turno activo
	var turnoID int
	err := bases.DB.QueryRow("SELECT id_turno FROM turnos WHERE fin IS NULL LIMIT 1").Scan(&turnoID)
	if err == sql.ErrNoRows {
		log.Println("No hay turno activo")
		ctx.Error("No hay turno activo, no se puede registrar el pedido", 500)
		return
	} else if err != nil {
		log.Println("Error al obtener turno:", err)
		ctx.Error("Error interno", 500)
		return
	}

	// Insertar el pedido principal sin columna 'fecha'
	result, err := bases.DB.Exec(
		"INSERT INTO pedidos_online (cliente, total, estado, tipo_pedido, notas, pedido_json, turno_id) VALUES (?, 0, 'pendiente', ?, ?, ?, ?)",
		cliente, tipoPedido, notasGenerales, pedidoJSON, turnoID,
	)
	if err != nil {
		log.Println("Error al crear pedido online:", err)
		ctx.Error("Error al crear pedido", 500)
		return
	}

	idOnline, _ := result.LastInsertId()
	total := 0
	agrupados := []ItemAgrupado{}

	// Agrupar productos idénticos con los mismos modificadores
	for _, item := range itemsJSON {
		modsKey := ""
		for _, m := range item.Mods {
			modsKey += m.ID + ","
		}
		encontrado := false
		for i, a := range agrupados {
			aModsKey := ""
			for _, m := range a.Mods {
				aModsKey += m.ID + ","
			}
			if a.IDPro == item.IDPro && aModsKey == modsKey {
				agrupados[i].Cantidad++
				encontrado = true
				break
			}
		}
		if !encontrado {
			agrupados = append(agrupados, ItemAgrupado{
				IDPro:    item.IDPro,
				Cantidad: 1,
				Mods:     item.Mods,
			})
		}
	}

	// Insertar detalles y modificadores
	for _, item := range agrupados {
		var precio int
		bases.DB.QueryRow("SELECT precio FROM productos WHERE id_pro = ?", item.IDPro).Scan(&precio)

		notaManual := ""
		for _, m := range item.Mods {
			if m.ID == "0" {
				notaManual = m.Nombre
			}
		}

		_, err := bases.DB.Exec(
			"INSERT INTO pedidos_online_detalle (id_online, id_pro, cantidad, precio, notas_producto) VALUES (?, ?, ?, ?, ?)",
			idOnline, item.IDPro, item.Cantidad, precio*item.Cantidad, notaManual,
		)
		if err != nil {
			log.Println("Error al insertar detalle:", err)
		}

		total += precio * item.Cantidad

		for _, mod := range item.Mods {
			if mod.ID != "0" {
				idMod, _ := strconv.Atoi(mod.ID)
				_, err := bases.DB.Exec(
					"INSERT INTO pedidos_online_modificadores (id_online, id_pro, id_mod) VALUES (?, ?, ?)",
					idOnline, item.IDPro, idMod,
				)
				if err != nil {
					log.Println("Error al insertar modificador:", err)
				}
			}
		}
	}

	// Actualizar total
	_, err = bases.DB.Exec("UPDATE pedidos_online SET total = ? WHERE id_online = ?", total, idOnline)
	if err != nil {
		log.Println("Error al actualizar total:", err)
	}

	// Redirigir a página de confirmación
	ctx.Redirect("/menu/confirmado/"+strconv.FormatInt(idOnline, 10), 302)
}

func PedidoConfirmado(ctx *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	type ConfirmadoData struct {
		IDOnline int
		Cliente  string
		Fecha    string
		Total    int
	}
	var data ConfirmadoData
	data.IDOnline = id
	data.Fecha = time.Now().Format("02/01/2006 15:04")
	bases.DB.QueryRow("SELECT cliente, total FROM pedidos_online WHERE id_online = ?", id).Scan(&data.Cliente, &data.Total)

	tmpl, err := template.ParseFiles("templates/menu_confirmado.html")
	if err != nil {
		log.Println("Error al cargar template:", err)
		ctx.Error("Error al cargar template", 500)
		return
	}
	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, data)
}
