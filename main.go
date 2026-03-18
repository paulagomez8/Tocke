package main

import (
	"log"
	"tockesanfelipe/modules/admin"
	"tockesanfelipe/modules/auth"
	"tockesanfelipe/modules/bases"
	"tockesanfelipe/modules/pedidos"
	"tockesanfelipe/modules/reportes"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func proteger(siguiente fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if !auth.VerificarSesion(ctx) {
			ctx.Redirect("/login", 302)
			return
		}
		siguiente(ctx)
	}
}

func main() {
	bases.ConectarDB()

	r := router.New()
	r.GET("/login", auth.Login)
	r.POST("/login", auth.Login)

	r.GET("/", proteger(pedidos.Inicio))
	r.GET("/nuevo-pedido", proteger(pedidos.NuevoPedido))
	r.POST("/confirmar-pedido", proteger(pedidos.ConfirmarPedido))

	r.GET("/admin", proteger(admin.VerAdmin))
	r.POST("/admin/producto/agregar", proteger(admin.AgregarProducto))
	r.POST("/admin/producto/editar/{id}", proteger(admin.EditarProducto))
	r.GET("/admin/producto/eliminar/{id}", proteger(admin.EliminarProducto))
	r.POST("/admin/categoria/agregar", proteger(admin.AgregarCategoria))
	r.GET("/admin/categoria/eliminar/{id}", proteger(admin.EliminarCategoria))
	r.POST("/admin/modificador/agregar/{id}", proteger(admin.AgregarModificador))
	r.GET("/admin/modificador/eliminar/{id}", proteger(admin.EliminarModificador))
	r.GET("/reportes", proteger(reportes.VerReportes))
	r.GET("/reportes/exportar", proteger(reportes.ExportarExcel))
	r.POST("/admin/credenciales", proteger(admin.CambiarCredenciales))
	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(fasthttp.ListenAndServe(":80", r.Handler))
}
