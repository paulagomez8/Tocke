package admin

import (
	"html/template"
	"log"
	"strconv"
	"tockesanfelipe/modules/bases"

	"github.com/valyala/fasthttp"
)

type Producto struct {
	ID        int
	Nombre    string
	Precio    int
	IDCat     int
	Categoria string
}

type Categoria struct {
	ID     int
	Nombre string
}

type AdminData struct {
	Productos     []Producto
	Categorias    []Categoria
	Modificadores []Modificador
}

func VerAdmin(ctx *fasthttp.RequestCtx) {
	var data AdminData

	// Obtener categorias
	rowsCat, err := bases.DB.Query("SELECT id_cat, nombre FROM categorias ORDER BY nombre")
	if err != nil {
		log.Println("Error al obtener categorias:", err)
		ctx.Error("Error al obtener categorias", 500)
		return
	}
	defer rowsCat.Close()
	for rowsCat.Next() {
		var c Categoria
		rowsCat.Scan(&c.ID, &c.Nombre)
		data.Categorias = append(data.Categorias, c)
	}

	// Obtener productos
	rowsPro, err := bases.DB.Query(`
        SELECT p.id_pro, p.nombre, p.precio, p.id_cat, c.nombre
        FROM productos p
        JOIN categorias c ON p.id_cat = c.id_cat
        ORDER BY c.nombre, p.nombre
    `)
	if err != nil {
		log.Println("Error al obtener productos:", err)
		ctx.Error("Error al obtener productos", 500)
		return
	}
	defer rowsPro.Close()
	for rowsPro.Next() {
		var p Producto
		rowsPro.Scan(&p.ID, &p.Nombre, &p.Precio, &p.IDCat, &p.Categoria)
		data.Productos = append(data.Productos, p)
	}

	tmpl, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		log.Println("Error al cargar template:", err)
		ctx.Error("Error al cargar template", 500)
		return
	}
	ctx.SetContentType("text/html")
	rowsMod, err := bases.DB.Query(`
        SELECT m.id_mod, m.id_pro, m.nombre
        FROM modificadores m
        ORDER BY m.id_pro
    `)
	if err != nil {
		log.Println("Error al obtener modificadores:", err)
	} else {
		defer rowsMod.Close()
		for rowsMod.Next() {
			var m Modificador
			rowsMod.Scan(&m.ID, &m.IDPro, &m.Nombre)
			data.Modificadores = append(data.Modificadores, m)
		}
	}
	tmpl.Execute(ctx, data)
}

// Productos
func AgregarProducto(ctx *fasthttp.RequestCtx) {
	nombre := string(ctx.FormValue("nombre"))
	precio, _ := strconv.Atoi(string(ctx.FormValue("precio")))
	idCat, _ := strconv.Atoi(string(ctx.FormValue("id_cat")))
	bases.DB.Exec("INSERT INTO productos (nombre, precio, id_cat) VALUES (?, ?, ?)", nombre, precio, idCat)
	ctx.Redirect("/admin", 302)
}

func EditarProducto(ctx *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	nombre := string(ctx.FormValue("nombre"))
	precio, _ := strconv.Atoi(string(ctx.FormValue("precio")))
	idCat, _ := strconv.Atoi(string(ctx.FormValue("id_cat")))
	bases.DB.Exec("UPDATE productos SET nombre=?, precio=?, id_cat=? WHERE id_pro=?", nombre, precio, idCat, id)
	ctx.Redirect("/admin", 302)
}

func EliminarProducto(ctx *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	bases.DB.Exec("DELETE FROM productos WHERE id_pro=?", id)
	ctx.Redirect("/admin", 302)
}

// Categorias
func AgregarCategoria(ctx *fasthttp.RequestCtx) {
	nombre := string(ctx.FormValue("nombre"))
	bases.DB.Exec("INSERT INTO categorias (nombre) VALUES (?)", nombre)
	ctx.Redirect("/admin", 302)
}

func EliminarCategoria(ctx *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	bases.DB.Exec("DELETE FROM categorias WHERE id_cat=?", id)
	ctx.Redirect("/admin", 302)
}

func AgregarModificador(ctx *fasthttp.RequestCtx) {
	idPro, _ := strconv.Atoi(ctx.UserValue("id").(string))
	nombre := string(ctx.FormValue("nombre"))
	bases.DB.Exec("INSERT INTO modificadores (id_pro, nombre) VALUES (?, ?)", idPro, nombre)
	ctx.Redirect("/admin", 302)
}

func EliminarModificador(ctx *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	bases.DB.Exec("DELETE FROM modificadores WHERE id_mod=?", id)
	ctx.Redirect("/admin", 302)
}

type Modificador struct {
	ID     int
	IDPro  int
	Nombre string
}
