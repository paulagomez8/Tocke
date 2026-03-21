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

type InventarioData struct {
	Ingredientes []Ingrediente
}

func VerInventario(ctx *fasthttp.RequestCtx) {
	var data InventarioData

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

	for _, d := range detalles {
		res, err := bases.DB.Exec(`
			UPDATE ingredientes i
			JOIN recetas r ON i.id_ing = r.id_ing
			SET i.stock = i.stock - (r.cantidad * ?)
			WHERE r.id_pro = ?
		`, d.cantidad, d.idPro)
		if err != nil {
			log.Println("Error al descontar stock:", err)
		} else {
			filas, _ := res.RowsAffected()
			log.Printf("Producto %d: %d filas afectadas\n", d.idPro, filas)
		}
	}
}
