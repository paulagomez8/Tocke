package reportes

import (
	"fmt"
	"html/template"
	"log"
	"tockesanfelipe/modules/bases"

	"github.com/valyala/fasthttp"
	"github.com/xuri/excelize/v2"
)

type VentaDia struct {
	Fecha   string
	Total   int
	Pedidos int
}

type VentaMes struct {
	Mes     string
	Total   int
	Pedidos int
}

type ProductoVendido struct {
	Nombre   string
	Cantidad int
	Total    int
}

type ReporteData struct {
	VentasDia    []VentaDia
	VentasMes    []VentaMes
	ProductosTop []ProductoVendido
}

func VerReportes(ctx *fasthttp.RequestCtx) {
	var data ReporteData

	// Ventas por dia (ultimos 30 dias)
	rowsDia, err := bases.DB.Query(`
    SELECT DATE_FORMAT(fecha, '%d-%m-%Y') as dia, SUM(total), COUNT(*)
    FROM pedidos
    WHERE fecha >= DATE_SUB(NOW(), INTERVAL 30 DAY)
    GROUP BY DATE(fecha)
    ORDER BY DATE(fecha) DESC
    `)
	if err != nil {
		log.Println("Error ventas dia:", err)
	} else {
		defer rowsDia.Close()
		for rowsDia.Next() {
			var v VentaDia
			rowsDia.Scan(&v.Fecha, &v.Total, &v.Pedidos)
			data.VentasDia = append(data.VentasDia, v)
		}
	}

	// Ventas por mes (ultimos 12 meses)
	rowsMes, err := bases.DB.Query(`
        SELECT DATE_FORMAT(fecha, '%Y-%m') as mes, SUM(total), COUNT(*)
        FROM pedidos
        GROUP BY mes
        ORDER BY mes DESC
        LIMIT 12
    `)
	if err != nil {
		log.Println("Error ventas mes:", err)
	} else {
		defer rowsMes.Close()
		for rowsMes.Next() {
			var v VentaMes
			rowsMes.Scan(&v.Mes, &v.Total, &v.Pedidos)
			data.VentasMes = append(data.VentasMes, v)
		}
	}

	// Productos mas vendidos
	rowsPro, err := bases.DB.Query(`
        SELECT p.nombre, SUM(pd.cantidad) as cantidad, SUM(pd.cantidad * pd.precio) as total
        FROM pedidos_detalle pd
        JOIN productos p ON pd.id_pro = p.id_pro
        GROUP BY p.id_pro, p.nombre
        ORDER BY cantidad DESC
        LIMIT 10
    `)
	if err != nil {
		log.Println("Error productos top:", err)
	} else {
		defer rowsPro.Close()
		for rowsPro.Next() {
			var p ProductoVendido
			rowsPro.Scan(&p.Nombre, &p.Cantidad, &p.Total)
			data.ProductosTop = append(data.ProductosTop, p)
		}
	}

	tmpl, err := template.ParseFiles("templates/reportes.html")
	if err != nil {
		log.Println("Error al cargar template:", err)
		ctx.Error("Error al cargar template", 500)
		return
	}
	ctx.SetContentType("text/html")
	tmpl.Execute(ctx, data)
}

func ExportarExcel(ctx *fasthttp.RequestCtx) {
	f := excelize.NewFile()

	// Hoja 1 - Ventas por mes
	f.NewSheet("Ventas por Mes")
	f.SetCellValue("Ventas por Mes", "A1", "Mes")
	f.SetCellValue("Ventas por Mes", "B1", "Pedidos")
	f.SetCellValue("Ventas por Mes", "C1", "Total")

	rows, _ := bases.DB.Query(`
        SELECT DATE_FORMAT(fecha, '%Y-%m') as mes, COUNT(*), SUM(total)
        FROM pedidos
        GROUP BY mes
        ORDER BY mes DESC
    `)
	defer rows.Close()
	i := 2
	for rows.Next() {
		var mes string
		var pedidos, total int
		rows.Scan(&mes, &pedidos, &total)
		f.SetCellValue("Ventas por Mes", fmt.Sprintf("A%d", i), mes)
		f.SetCellValue("Ventas por Mes", fmt.Sprintf("B%d", i), pedidos)
		f.SetCellValue("Ventas por Mes", fmt.Sprintf("C%d", i), total)
		i++
	}

	// Hoja 2 - Ventas por dia
	f.NewSheet("Ventas por Día")
	f.SetCellValue("Ventas por Día", "A1", "Fecha")
	f.SetCellValue("Ventas por Día", "B1", "Pedidos")
	f.SetCellValue("Ventas por Día", "C1", "Total")

	rows2, _ := bases.DB.Query(`
        SELECT DATE(fecha), COUNT(*), SUM(total)
        FROM pedidos
        GROUP BY DATE(fecha)
        ORDER BY DATE(fecha) DESC
    `)
	defer rows2.Close()
	i = 2
	for rows2.Next() {
		var fecha string
		var pedidos, total int
		rows2.Scan(&fecha, &pedidos, &total)
		f.SetCellValue("Ventas por Día", fmt.Sprintf("A%d", i), fecha)
		f.SetCellValue("Ventas por Día", fmt.Sprintf("B%d", i), pedidos)
		f.SetCellValue("Ventas por Día", fmt.Sprintf("C%d", i), total)
		i++
	}

	// Hoja 3 - Productos mas vendidos
	f.NewSheet("Productos")
	f.SetCellValue("Productos", "A1", "Producto")
	f.SetCellValue("Productos", "B1", "Cantidad")
	f.SetCellValue("Productos", "C1", "Total")

	rows3, _ := bases.DB.Query(`
        SELECT p.nombre, SUM(pd.cantidad), SUM(pd.cantidad * pd.precio)
        FROM pedidos_detalle pd
        JOIN productos p ON pd.id_pro = p.id_pro
        GROUP BY p.id_pro, p.nombre
        ORDER BY SUM(pd.cantidad) DESC
    `)
	defer rows3.Close()
	i = 2
	for rows3.Next() {
		var nombre string
		var cantidad, total int
		rows3.Scan(&nombre, &cantidad, &total)
		f.SetCellValue("Productos", fmt.Sprintf("A%d", i), nombre)
		f.SetCellValue("Productos", fmt.Sprintf("B%d", i), cantidad)
		f.SetCellValue("Productos", fmt.Sprintf("C%d", i), total)
		i++
	}

	// Hoja 4 - Resumen general
	f.NewSheet("Resumen")
	f.SetCellValue("Resumen", "A1", "Total Pedidos")
	f.SetCellValue("Resumen", "B1", "Total Ventas")
	f.SetCellValue("Resumen", "C1", "Ticket Promedio")

	var totalPedidos, totalVentas int
	bases.DB.QueryRow("SELECT COUNT(*), COALESCE(SUM(total),0) FROM pedidos").Scan(&totalPedidos, &totalVentas)
	promedio := 0
	if totalPedidos > 0 {
		promedio = totalVentas / totalPedidos
	}
	f.SetCellValue("Resumen", "A2", totalPedidos)
	f.SetCellValue("Resumen", "B2", totalVentas)
	f.SetCellValue("Resumen", "C2", promedio)

	// Hoja 5 - Detalle de pedidos
	f.NewSheet("Detalle Pedidos")
	f.SetCellValue("Detalle Pedidos", "A1", "Pedido")
	f.SetCellValue("Detalle Pedidos", "B1", "Fecha")
	f.SetCellValue("Detalle Pedidos", "C1", "Cliente")
	f.SetCellValue("Detalle Pedidos", "D1", "Producto")
	f.SetCellValue("Detalle Pedidos", "E1", "Cantidad")
	f.SetCellValue("Detalle Pedidos", "F1", "Precio")
	f.SetCellValue("Detalle Pedidos", "G1", "Subtotal")
	f.SetCellValue("Detalle Pedidos", "H1", "Total Pedido")

	rows4, _ := bases.DB.Query(`
        SELECT p.id_ped, p.fecha, p.cliente, pr.nombre, pd.cantidad, pd.precio, 
               (pd.cantidad * pd.precio), p.total
        FROM pedidos p
        JOIN pedidos_detalle pd ON p.id_ped = pd.id_ped
        JOIN productos pr ON pd.id_pro = pr.id_pro
        ORDER BY p.fecha DESC
    `)
	defer rows4.Close()
	i = 2
	for rows4.Next() {
		var idPed, cantidad, precio, subtotal, total int
		var fecha, cliente, producto string
		rows4.Scan(&idPed, &fecha, &cliente, &producto, &cantidad, &precio, &subtotal, &total)
		f.SetCellValue("Detalle Pedidos", fmt.Sprintf("A%d", i), idPed)
		f.SetCellValue("Detalle Pedidos", fmt.Sprintf("B%d", i), fecha)
		f.SetCellValue("Detalle Pedidos", fmt.Sprintf("C%d", i), cliente)
		f.SetCellValue("Detalle Pedidos", fmt.Sprintf("D%d", i), producto)
		f.SetCellValue("Detalle Pedidos", fmt.Sprintf("E%d", i), cantidad)
		f.SetCellValue("Detalle Pedidos", fmt.Sprintf("F%d", i), precio)
		f.SetCellValue("Detalle Pedidos", fmt.Sprintf("G%d", i), subtotal)
		f.SetCellValue("Detalle Pedidos", fmt.Sprintf("H%d", i), total)
		i++
	}

	// Eliminar hoja default
	f.DeleteSheet("Sheet1")

	ctx.Response.Header.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Response.Header.Set("Content-Disposition", "attachment; filename=reporte_tocke.xlsx")

	f.Write(ctx)
}
