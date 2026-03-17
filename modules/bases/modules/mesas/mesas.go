package mesas

import (
	"html/template"
	"log"
	"tockesanfelipe/modules/bases"

	"github.com/valyala/fasthttp"
)

type Mesa struct {
	ID     int
	Numero int
	Estado string
}

func VerMesas(ctx *fasthttp.RequestCtx) {
	rows, err := bases.DB.Query("SELECT id, numero, estado FROM mesas")
	if err != nil {
		log.Println("Error al obtener mesas:", err)
		ctx.Error("Error al obtener mesas", 500)
		return
	}
	defer rows.Close()

	var listaMesas []Mesa
	for rows.Next() {
		var m Mesa
		rows.Scan(&m.ID, &m.Numero, &m.Estado)
		listaMesas = append(listaMesas, m)
	}

	tmpl, err := template.ParseFiles("templates/mesas.html")
	if err != nil {
		log.Println("Error al cargar template:", err)
		ctx.Error("Error al cargar template", 500)
		return
	}

	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, listaMesas)
}
