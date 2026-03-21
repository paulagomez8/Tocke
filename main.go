package main

import (
	"log"
	"tockesanfelipe/modules/admin"
	"tockesanfelipe/modules/auth"
	"tockesanfelipe/modules/bases"
	"tockesanfelipe/modules/inventario"
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
	r.GET("/inventario", proteger(inventario.VerInventario))
	r.POST("/inventario/ingrediente/agregar", proteger(inventario.AgregarIngrediente))
	r.POST("/inventario/ingrediente/editar/{id}", proteger(inventario.EditarStock))
	r.GET("/inventario/ingrediente/eliminar/{id}", proteger(inventario.EliminarIngrediente))
	r.POST("/inventario/receta/agregar/{id}", proteger(inventario.AgregarReceta))
	r.GET("/inventario/receta/eliminar/{idpro}/{idping}", proteger(inventario.EliminarReceta))
	r.POST("/admin/categoria/editar/{id}", proteger(admin.EditarCategoria))
	r.POST("/admin/modificador/editar/{id}", proteger(admin.EditarModificador))
	r.POST("/mesa/liberar/{id}", proteger(pedidos.LiberarMesa))
	r.POST("/turno/abrir", proteger(reportes.AbrirTurno))
	r.POST("/turno/cerrar/{id}", proteger(reportes.CerrarTurno))

	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
