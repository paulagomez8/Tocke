package inventario

import (
	"html/template"
	"log"
	"strconv"
	"tockesanfelipe/modules/bases"

	"github.com/valyala/fasthttp"
)

type Ingrediente struct {
	ID          int
	Nombre      string
	Stock       int
	StockMinimo int
	Unidad      string
	Alerta      bool
}

type RecetaItem struct {
	IDIng     int
	NombreIng string
	Cantidad  int
}

type ProductoReceta struct {
	ID     int
	Nombre string
	Receta []RecetaItem
}

type InventarioData struct {
	Ingredientes []Ingrediente
	Productos    []ProductoReceta
}

func VerInventario(ctx *fasthttp.RequestCtx) {
	var data InventarioData

	// Obtener ingredientes
	rows, err := bases.DB.Query("SELECT id_ing, nombre, stock, stock_minimo, unidad FROM ingredientes ORDER BY nombre")
	if err != nil {
		log.Println("Error al obtener ingredientes:", err)
		ctx.Error("Error al obtener ingredientes", 500)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var i Ingrediente
		rows.Scan(&i.ID, &i.Nombre, &i.Stock, &i.StockMinimo, &i.Unidad)
		i.Alerta = i.Stock <= i.StockMinimo
		data.Ingredientes = append(data.Ingredientes, i)
	}

	// Obtener todos los productos
	rowsPro, err := bases.DB.Query("SELECT id_pro, nombre FROM productos ORDER BY nombre")
	if err != nil {
		log.Println("Error al obtener productos:", err)
		ctx.Error("Error al obtener productos", 500)
		return
	}
	defer rowsPro.Close()
	for rowsPro.Next() {
		var p ProductoReceta
		rowsPro.Scan(&p.ID, &p.Nombre)

		rowsRec, err := bases.DB.Query(`
			SELECT r.id_ing, i.nombre, r.cantidad
			FROM recetas r
			JOIN ingredientes i ON r.id_ing = i.id_ing
			WHERE r.id_pro = ?
		`, p.ID)
		if err == nil {
			defer rowsRec.Close()
			for rowsRec.Next() {
				var r RecetaItem
				rowsRec.Scan(&r.IDIng, &r.NombreIng, &r.Cantidad)
				p.Receta = append(p.Receta, r)
			}
		}
		data.Productos = append(data.Productos, p)
	}

	tmpl, err := template.ParseFiles("templates/inventario.html")
	if err != nil {
		log.Println("Error al cargar template:", err)
		ctx.Error("Error al cargar template", 500)
		return
	}
	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, data)
}

func AgregarIngrediente(ctx *fasthttp.RequestCtx) {
	nombre := string(ctx.FormValue("nombre"))
	stock, _ := strconv.Atoi(string(ctx.FormValue("stock")))
	stockMinimo, _ := strconv.Atoi(string(ctx.FormValue("stock_minimo")))
	unidad := string(ctx.FormValue("unidad"))
	bases.DB.Exec("INSERT INTO ingredientes (nombre, stock, stock_minimo, unidad) VALUES (?, ?, ?, ?)", nombre, stock, stockMinimo, unidad)
	ctx.Redirect("/inventario", 302)
}

func EditarStock(ctx *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	nombre := string(ctx.FormValue("nombre"))
	stock, _ := strconv.Atoi(string(ctx.FormValue("stock")))
	stockMinimo, _ := strconv.Atoi(string(ctx.FormValue("stock_minimo")))
	bases.DB.Exec("UPDATE ingredientes SET nombre=?, stock=?, stock_minimo=? WHERE id_ing=?", nombre, stock, stockMinimo, id)
	ctx.Redirect("/inventario", 302)
}

func EliminarIngrediente(ctx *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	bases.DB.Exec("DELETE FROM recetas WHERE id_ing=?", id)
	bases.DB.Exec("DELETE FROM ingredientes WHERE id_ing=?", id)
	ctx.Redirect("/inventario", 302)
}

func AgregarReceta(ctx *fasthttp.RequestCtx) {
	idPro, _ := strconv.Atoi(ctx.UserValue("id").(string))
	idIng, _ := strconv.Atoi(string(ctx.FormValue("id_ing")))
	cantidad, _ := strconv.Atoi(string(ctx.FormValue("cantidad")))
	bases.DB.Exec("INSERT INTO recetas (id_pro, id_ing, cantidad) VALUES (?, ?, ?)", idPro, idIng, cantidad)
	ctx.Redirect("/inventario", 302)
}

func EliminarReceta(ctx *fasthttp.RequestCtx) {
	idPro, _ := strconv.Atoi(ctx.UserValue("idpro").(string))
	idIng, _ := strconv.Atoi(ctx.UserValue("idping").(string))
	bases.DB.Exec("DELETE FROM recetas WHERE id_pro=? AND id_ing=?", idPro, idIng)
	ctx.Redirect("/inventario", 302)
}

func ObtenerAlertas() []Ingrediente {
	var alertas []Ingrediente
	rows, err := bases.DB.Query("SELECT id_ing, nombre, stock, stock_minimo, unidad FROM ingredientes WHERE stock <= stock_minimo")
	if err != nil {
		return alertas
	}
	defer rows.Close()
	for rows.Next() {
		var i Ingrediente
		rows.Scan(&i.ID, &i.Nombre, &i.Stock, &i.StockMinimo, &i.Unidad)
		i.Alerta = true
		alertas = append(alertas, i)
	}
	return alertas
}
func DescontarStock(idPedido int64) {
	var tipoPedido string
	bases.DB.QueryRow("SELECT tipo_pedido FROM pedidos WHERE id_ped=?", idPedido).Scan(&tipoPedido)
	log.Printf(">>> tipoPedido: '%s'\n", tipoPedido)

	if tipoPedido == "Llevar" || tipoPedido == "Retiro" || tipoPedido == "Delivery" {
		bases.DB.Exec("UPDATE ingredientes SET stock = stock - 1 WHERE nombre = 'Plumavit'")
	}

	rows, err := bases.DB.Query("SELECT id_pro, cantidad FROM pedidos_detalle WHERE id_ped=?", idPedido)
	if err != nil {
		log.Println("Error DescontarStock:", err)
		return
	}
	defer rows.Close()

	type detalle struct{ idPro, cantidad int }
	var detalles []detalle
	for rows.Next() {
		var d detalle
		rows.Scan(&d.idPro, &d.cantidad)
		detalles = append(detalles, d)
	}
	log.Printf(">>> detalles encontrados: %d\n", len(detalles))

	for _, d := range detalles {
		var tieneReceta int
		bases.DB.QueryRow("SELECT COUNT(*) FROM recetas WHERE id_pro=?", d.idPro).Scan(&tieneReceta)
		log.Printf(">>> producto %d cantidad %d tieneReceta %d\n", d.idPro, d.cantidad, tieneReceta)

		if tieneReceta > 0 {
			recetas, err := bases.DB.Query("SELECT id_ing, cantidad FROM recetas WHERE id_pro=?", d.idPro)
			if err != nil {
				continue
			}
			for recetas.Next() {
				var idIng, cantRec int
				recetas.Scan(&idIng, &cantRec)
				log.Printf(">>> descontando ingrediente %d cantidad %d\n", idIng, cantRec*d.cantidad)
				bases.DB.Exec("UPDATE ingredientes SET stock = stock - ? WHERE id_ing=?", cantRec*d.cantidad, idIng)
			}
			recetas.Close()
		} else {
			var nombrePro string
			bases.DB.QueryRow("SELECT nombre FROM productos WHERE id_pro=?", d.idPro).Scan(&nombrePro)
			log.Printf(">>> sin receta, descontando por nombre: '%s'\n", nombrePro)
			bases.DB.Exec("UPDATE ingredientes SET stock = stock - ? WHERE nombre = ?", d.cantidad, nombrePro)
		}
	}
}
func DescontarStockOnline(idOnline int64) {
	// Obtener tipo de pedido
	var tipoPedido string
	bases.DB.QueryRow("SELECT tipo_pedido FROM pedidos_online WHERE id_online=?", idOnline).Scan(&tipoPedido)

	// Descontar envase si corresponde
	if tipoPedido == "Llevar" || tipoPedido == "Retiro" || tipoPedido == "Delivery" {
		bases.DB.Exec("UPDATE ingredientes SET stock = stock - 1 WHERE nombre = 'Plumavit'")
	}

	// Obtener detalle del pedido online
	rows, err := bases.DB.Query("SELECT id_pro, cantidad FROM pedidos_online_detalle WHERE id_online=?", idOnline)
	if err != nil {
		log.Println("Error DescontarStockOnline al obtener detalle:", err)
		return
	}
	defer rows.Close()

	type detalle struct{ idPro, cantidad int }
	var detalles []detalle
	for rows.Next() {
		var d detalle
		rows.Scan(&d.idPro, &d.cantidad)
		detalles = append(detalles, d)
	}

	for _, d := range detalles {
		// Verificar si el producto tiene receta
		var tieneReceta int
		bases.DB.QueryRow("SELECT COUNT(*) FROM recetas WHERE id_pro=?", d.idPro).Scan(&tieneReceta)

		if tieneReceta > 0 {
			// Tiene receta: descontar ingredientes
			recetas, err := bases.DB.Query("SELECT id_ing, cantidad FROM recetas WHERE id_pro=?", d.idPro)
			if err != nil {
				log.Println("Error al obtener receta del producto:", d.idPro, err)
				continue
			}
			for recetas.Next() {
				var idIng, cantRec int
				recetas.Scan(&idIng, &cantRec)
				_, err := bases.DB.Exec(
					"UPDATE ingredientes SET stock = stock - ? WHERE id_ing=?",
					cantRec*d.cantidad, idIng,
				)
				if err != nil {
					log.Println("Error al descontar ingrediente:", idIng, err)
				}
			}
			recetas.Close()
		} else {
			// Sin receta: descontar por nombre del producto
			var nombrePro string
			bases.DB.QueryRow("SELECT nombre FROM productos WHERE id_pro=?", d.idPro).Scan(&nombrePro)
			res, err := bases.DB.Exec(
				"UPDATE ingredientes SET stock = stock - ? WHERE nombre = ?",
				d.cantidad, nombrePro,
			)
			if err != nil {
				log.Println("Error al descontar producto sin receta:", nombrePro, err)
			} else {
				filas, _ := res.RowsAffected()
				if filas == 0 {
					log.Printf("Sin receta y sin ingrediente matching: %s (id:%d)\n", nombrePro, d.idPro)
				}
			}
		}
	}
}
