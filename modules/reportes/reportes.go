package reportes

import (
	"fmt"
	"html/template"
	"log"
	"strconv"
	"tockesanfelipe/modules/bases"

	"github.com/valyala/fasthttp"
	"github.com/xuri/excelize/v2"
)

type VentaDia struct {
	Fecha      string
	Total      int
	Pedidos    int
	Servir     int
	Retiro     int
	Llevar     int
	Delivery   int
	IrComiendo int
	Online     int
}

type VentaMes struct {
	Mes       string
	MesNombre string
	Total     int
	Pedidos   int
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
	FechaFiltro  string
}

func VerReportes(ctx *fasthttp.RequestCtx) {
	var data ReporteData

	fecha := string(ctx.QueryArgs().Peek("fecha"))
	data.FechaFiltro = fecha

	// Ventas por mes (local + online combinados)
	rowsMes, err := bases.DB.Query(`
		SELECT mes,
		       CASE CAST(SUBSTRING(mes, 6, 2) AS UNSIGNED)
		           WHEN 1 THEN 'Enero' WHEN 2 THEN 'Febrero' WHEN 3 THEN 'Marzo'
		           WHEN 4 THEN 'Abril' WHEN 5 THEN 'Mayo' WHEN 6 THEN 'Junio'
		           WHEN 7 THEN 'Julio' WHEN 8 THEN 'Agosto' WHEN 9 THEN 'Septiembre'
		           WHEN 10 THEN 'Octubre' WHEN 11 THEN 'Noviembre' WHEN 12 THEN 'Diciembre'
		       END as mes_nombre,
		       SUM(total), SUM(pedidos)
		FROM (
		    SELECT DATE_FORMAT(fecha, '%Y-%m') as mes, SUM(total) as total, COUNT(*) as pedidos
		    FROM pedidos GROUP BY mes
		    UNION ALL
		    SELECT DATE_FORMAT(fecha, '%Y-%m') as mes, SUM(total) as total, COUNT(*) as pedidos
		    FROM pedidos_online WHERE estado = 'listo' GROUP BY mes
		) t
		GROUP BY mes ORDER BY mes DESC
	`)
	if err != nil {
		log.Println("Error ventas mes:", err)
	} else {
		defer rowsMes.Close()
		for rowsMes.Next() {
			var v VentaMes
			rowsMes.Scan(&v.Mes, &v.MesNombre, &v.Total, &v.Pedidos)
			data.VentasMes = append(data.VentasMes, v)
		}
	}

	// Ventas por dia con desglose de tipos
	var queryDia string
	var argsDia []interface{}

	if fecha != "" {
		queryDia = `
			SELECT dia, SUM(pedidos), SUM(total),
			       SUM(servir), SUM(retiro), SUM(llevar), SUM(delivery), SUM(ircomiendo), SUM(online)
			FROM (
			    SELECT DATE_FORMAT(fecha, '%d-%m-%Y') as dia,
			           COUNT(*) as pedidos, SUM(total) as total,
			           SUM(tipo_pedido = 'Servir') as servir,
			           SUM(tipo_pedido = 'Retiro') as retiro,
			           SUM(tipo_pedido = 'Llevar') as llevar,
			           SUM(tipo_pedido = 'Delivery') as delivery,
			           SUM(tipo_pedido = 'Ir comiendo') as ircomiendo,
			           0 as online
			    FROM pedidos WHERE DATE(fecha) = ?
			    GROUP BY DATE(fecha)
			    UNION ALL
			    SELECT DATE_FORMAT(fecha, '%d-%m-%Y') as dia,
			           COUNT(*) as pedidos, SUM(total) as total,
			           0, 0, 0, 0, 0, COUNT(*) as online
			    FROM pedidos_online WHERE estado = 'listo' AND DATE(fecha) = ?
			    GROUP BY DATE(fecha)
			) t
			GROUP BY dia ORDER BY dia DESC
		`
		argsDia = []interface{}{fecha, fecha}
	} else {
		queryDia = `
			SELECT dia, SUM(pedidos), SUM(total),
			       SUM(servir), SUM(retiro), SUM(llevar), SUM(delivery), SUM(ircomiendo), SUM(online)
			FROM (
			    SELECT DATE_FORMAT(fecha, '%d-%m-%Y') as dia,
			           COUNT(*) as pedidos, SUM(total) as total,
			           SUM(tipo_pedido = 'Servir') as servir,
			           SUM(tipo_pedido = 'Retiro') as retiro,
			           SUM(tipo_pedido = 'Llevar') as llevar,
			           SUM(tipo_pedido = 'Delivery') as delivery,
			           SUM(tipo_pedido = 'Ir comiendo') as ircomiendo,
			           0 as online
			    FROM pedidos
			    WHERE DATE_FORMAT(fecha, '%Y-%m') = DATE_FORMAT(NOW(), '%Y-%m')
			    GROUP BY DATE(fecha)
			    UNION ALL
			    SELECT DATE_FORMAT(fecha, '%d-%m-%Y') as dia,
			           COUNT(*) as pedidos, SUM(total) as total,
			           0, 0, 0, 0, 0, COUNT(*) as online
			    FROM pedidos_online
			    WHERE estado = 'listo' AND DATE_FORMAT(fecha, '%Y-%m') = DATE_FORMAT(NOW(), '%Y-%m')
			    GROUP BY DATE(fecha)
			) t
			GROUP BY dia ORDER BY dia DESC
		`
		argsDia = []interface{}{}
	}

	rowsDia, err := bases.DB.Query(queryDia, argsDia...)
	if err != nil {
		log.Println("Error ventas dia:", err)
	} else {
		defer rowsDia.Close()
		for rowsDia.Next() {
			var v VentaDia
			rowsDia.Scan(&v.Fecha, &v.Pedidos, &v.Total,
				&v.Servir, &v.Retiro, &v.Llevar, &v.Delivery, &v.IrComiendo, &v.Online)
			data.VentasDia = append(data.VentasDia, v)
		}
	}

	// Productos mas vendidos (local + online)
	var queryPro string
	var argsPro []interface{}

	if fecha != "" {
		queryPro = `
			SELECT nombre, SUM(cantidad), SUM(total)
			FROM (
			    SELECT p.nombre, SUM(pd.cantidad) as cantidad, SUM(pd.cantidad * pd.precio) as total
			    FROM pedidos_detalle pd
			    JOIN productos p ON pd.id_pro = p.id_pro
			    JOIN pedidos pe ON pd.id_ped = pe.id_ped
			    WHERE DATE(pe.fecha) = ?
			    GROUP BY p.id_pro, p.nombre
			    UNION ALL
			    SELECT p.nombre, SUM(pod.cantidad) as cantidad, SUM(pod.precio) as total
			    FROM pedidos_online_detalle pod
			    JOIN productos p ON pod.id_pro = p.id_pro
			    JOIN pedidos_online po ON pod.id_online = po.id_online
			    WHERE po.estado = 'listo' AND DATE(po.fecha) = ?
			    GROUP BY p.id_pro, p.nombre
			) t
			GROUP BY nombre ORDER BY SUM(cantidad) DESC LIMIT 10
		`
		argsPro = []interface{}{fecha, fecha}
	} else {
		queryPro = `
			SELECT nombre, SUM(cantidad), SUM(total)
			FROM (
			    SELECT p.nombre, SUM(pd.cantidad) as cantidad, SUM(pd.cantidad * pd.precio) as total
			    FROM pedidos_detalle pd
			    JOIN productos p ON pd.id_pro = p.id_pro
			    JOIN pedidos pe ON pd.id_ped = pe.id_ped
			    WHERE DATE_FORMAT(pe.fecha, '%Y-%m') = DATE_FORMAT(NOW(), '%Y-%m')
			    GROUP BY p.id_pro, p.nombre
			    UNION ALL
			    SELECT p.nombre, SUM(pod.cantidad) as cantidad, SUM(pod.precio) as total
			    FROM pedidos_online_detalle pod
			    JOIN productos p ON pod.id_pro = p.id_pro
			    JOIN pedidos_online po ON pod.id_online = po.id_online
			    WHERE po.estado = 'listo' AND DATE_FORMAT(po.fecha, '%Y-%m') = DATE_FORMAT(NOW(), '%Y-%m')
			    GROUP BY p.id_pro, p.nombre
			) t
			GROUP BY nombre ORDER BY SUM(cantidad) DESC LIMIT 10
		`
		argsPro = []interface{}{}
	}

	rowsPro, err := bases.DB.Query(queryPro, argsPro...)
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

type Turno struct {
	ID     int
	Nombre string
	Inicio string
	Fin    string
}

func AbrirTurno(ctx *fasthttp.RequestCtx) {
	nombre := string(ctx.FormValue("nombre"))
	if nombre == "" {
		nombre = "Turno"
	}
	bases.DB.Exec("INSERT INTO turnos (nombre, inicio, fin) VALUES (?, NOW(), NULL)", nombre)
	ctx.Redirect("/", 302)
}

func CerrarTurno(ctx *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(ctx.UserValue("id").(string))
	bases.DB.Exec("UPDATE turnos SET fin = NOW() WHERE id_turno = ?", id)
	ctx.Redirect("/", 302)
}

func ExportarExcel(ctx *fasthttp.RequestCtx) {
	f := excelize.NewFile()

	f.NewSheet("Ventas por Mes")
	f.SetCellValue("Ventas por Mes", "A1", "Mes")
	f.SetCellValue("Ventas por Mes", "B1", "Pedidos")
	f.SetCellValue("Ventas por Mes", "C1", "Total")

	rows, _ := bases.DB.Query(`
		SELECT mes, SUM(pedidos), SUM(total)
		FROM (
		    SELECT DATE_FORMAT(fecha, '%Y-%m') as mes, COUNT(*) as pedidos, SUM(total) as total
		    FROM pedidos GROUP BY mes
		    UNION ALL
		    SELECT DATE_FORMAT(fecha, '%Y-%m') as mes, COUNT(*) as pedidos, SUM(total) as total
		    FROM pedidos_online WHERE estado = 'listo' GROUP BY mes
		) t GROUP BY mes ORDER BY mes DESC
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

	f.NewSheet("Ventas por Día")
	f.SetCellValue("Ventas por Día", "A1", "Fecha")
	f.SetCellValue("Ventas por Día", "B1", "Pedidos")
	f.SetCellValue("Ventas por Día", "C1", "Total")
	f.SetCellValue("Ventas por Día", "D1", "Servir")
	f.SetCellValue("Ventas por Día", "E1", "Retiro")
	f.SetCellValue("Ventas por Día", "F1", "Llevar")
	f.SetCellValue("Ventas por Día", "G1", "Delivery")
	f.SetCellValue("Ventas por Día", "H1", "Ir comiendo")
	f.SetCellValue("Ventas por Día", "I1", "Online")

	rows2, _ := bases.DB.Query(`
		SELECT dia, SUM(pedidos), SUM(total),
		       SUM(servir), SUM(retiro), SUM(llevar), SUM(delivery), SUM(ircomiendo), SUM(online)
		FROM (
		    SELECT DATE(fecha) as dia, COUNT(*) as pedidos, SUM(total) as total,
		           SUM(tipo_pedido = 'Servir') as servir,
		           SUM(tipo_pedido = 'Retiro') as retiro,
		           SUM(tipo_pedido = 'Llevar') as llevar,
		           SUM(tipo_pedido = 'Delivery') as delivery,
		           SUM(tipo_pedido = 'Ir comiendo') as ircomiendo,
		           0 as online
		    FROM pedidos GROUP BY DATE(fecha)
		    UNION ALL
		    SELECT DATE(fecha) as dia, COUNT(*) as pedidos, SUM(total) as total,
		           0, 0, 0, 0, 0, COUNT(*) as online
		    FROM pedidos_online WHERE estado = 'listo' GROUP BY DATE(fecha)
		) t GROUP BY dia ORDER BY dia DESC
	`)
	defer rows2.Close()
	i = 2
	for rows2.Next() {
		var fecha string
		var pedidos, total, servir, retiro, llevar, delivery, ircomiendo, online int
		rows2.Scan(&fecha, &pedidos, &total, &servir, &retiro, &llevar, &delivery, &ircomiendo, &online)
		f.SetCellValue("Ventas por Día", fmt.Sprintf("A%d", i), fecha)
		f.SetCellValue("Ventas por Día", fmt.Sprintf("B%d", i), pedidos)
		f.SetCellValue("Ventas por Día", fmt.Sprintf("C%d", i), total)
		f.SetCellValue("Ventas por Día", fmt.Sprintf("D%d", i), servir)
		f.SetCellValue("Ventas por Día", fmt.Sprintf("E%d", i), retiro)
		f.SetCellValue("Ventas por Día", fmt.Sprintf("F%d", i), llevar)
		f.SetCellValue("Ventas por Día", fmt.Sprintf("G%d", i), delivery)
		f.SetCellValue("Ventas por Día", fmt.Sprintf("H%d", i), ircomiendo)
		f.SetCellValue("Ventas por Día", fmt.Sprintf("I%d", i), online)
		i++
	}

	f.NewSheet("Productos")
	f.SetCellValue("Productos", "A1", "Producto")
	f.SetCellValue("Productos", "B1", "Cantidad")
	f.SetCellValue("Productos", "C1", "Total")

	rows3, _ := bases.DB.Query(`
		SELECT nombre, SUM(cantidad), SUM(total)
		FROM (
		    SELECT p.nombre, SUM(pd.cantidad) as cantidad, SUM(pd.cantidad * pd.precio) as total
		    FROM pedidos_detalle pd
		    JOIN productos p ON pd.id_pro = p.id_pro
		    GROUP BY p.id_pro, p.nombre
		    UNION ALL
		    SELECT p.nombre, SUM(pod.cantidad) as cantidad, SUM(pod.precio) as total
		    FROM pedidos_online_detalle pod
		    JOIN productos p ON pod.id_pro = p.id_pro
		    JOIN pedidos_online po ON pod.id_online = po.id_online
		    WHERE po.estado = 'listo'
		    GROUP BY p.id_pro, p.nombre
		) t GROUP BY nombre ORDER BY SUM(cantidad) DESC
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

	f.NewSheet("Resumen")
	f.SetCellValue("Resumen", "A1", "Total Pedidos")
	f.SetCellValue("Resumen", "B1", "Total Ventas")
	f.SetCellValue("Resumen", "C1", "Ticket Promedio")

	var totalPedidos, totalVentas int
	bases.DB.QueryRow(`
		SELECT SUM(pedidos), SUM(total)
		FROM (
		    SELECT COUNT(*) as pedidos, COALESCE(SUM(total),0) as total FROM pedidos
		    UNION ALL
		    SELECT COUNT(*) as pedidos, COALESCE(SUM(total),0) as total FROM pedidos_online WHERE estado = 'listo'
		) t
	`).Scan(&totalPedidos, &totalVentas)
	promedio := 0
	if totalPedidos > 0 {
		promedio = totalVentas / totalPedidos
	}
	f.SetCellValue("Resumen", "A2", totalPedidos)
	f.SetCellValue("Resumen", "B2", totalVentas)
	f.SetCellValue("Resumen", "C2", promedio)

	f.NewSheet("Pedidos Local")
	f.SetCellValue("Pedidos Local", "A1", "Pedido")
	f.SetCellValue("Pedidos Local", "B1", "Fecha")
	f.SetCellValue("Pedidos Local", "C1", "Cliente")
	f.SetCellValue("Pedidos Local", "D1", "Tipo")
	f.SetCellValue("Pedidos Local", "E1", "Producto")
	f.SetCellValue("Pedidos Local", "F1", "Cantidad")
	f.SetCellValue("Pedidos Local", "G1", "Precio")
	f.SetCellValue("Pedidos Local", "H1", "Subtotal")
	f.SetCellValue("Pedidos Local", "I1", "Total Pedido")

	rows4, _ := bases.DB.Query(`
		SELECT p.id_ped, p.fecha, p.cliente, p.tipo_pedido, pr.nombre, pd.cantidad, pd.precio,
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
		var fecha, cliente, tipo, producto string
		rows4.Scan(&idPed, &fecha, &cliente, &tipo, &producto, &cantidad, &precio, &subtotal, &total)
		f.SetCellValue("Pedidos Local", fmt.Sprintf("A%d", i), idPed)
		f.SetCellValue("Pedidos Local", fmt.Sprintf("B%d", i), fecha)
		f.SetCellValue("Pedidos Local", fmt.Sprintf("C%d", i), cliente)
		f.SetCellValue("Pedidos Local", fmt.Sprintf("D%d", i), tipo)
		f.SetCellValue("Pedidos Local", fmt.Sprintf("E%d", i), producto)
		f.SetCellValue("Pedidos Local", fmt.Sprintf("F%d", i), cantidad)
		f.SetCellValue("Pedidos Local", fmt.Sprintf("G%d", i), precio)
		f.SetCellValue("Pedidos Local", fmt.Sprintf("H%d", i), subtotal)
		f.SetCellValue("Pedidos Local", fmt.Sprintf("I%d", i), total)
		i++
	}

	f.NewSheet("Pedidos Online")
	f.SetCellValue("Pedidos Online", "A1", "Pedido")
	f.SetCellValue("Pedidos Online", "B1", "Fecha")
	f.SetCellValue("Pedidos Online", "C1", "Cliente")
	f.SetCellValue("Pedidos Online", "D1", "Tipo")
	f.SetCellValue("Pedidos Online", "E1", "Producto")
	f.SetCellValue("Pedidos Online", "F1", "Cantidad")
	f.SetCellValue("Pedidos Online", "G1", "Total Pedido")

	rows5, _ := bases.DB.Query(`
		SELECT po.id_online, po.fecha, po.cliente, po.tipo_pedido, pr.nombre, pod.cantidad, po.total
		FROM pedidos_online po
		JOIN pedidos_online_detalle pod ON po.id_online = pod.id_online
		JOIN productos pr ON pod.id_pro = pr.id_pro
		WHERE po.estado = 'listo'
		ORDER BY po.fecha DESC
	`)
	defer rows5.Close()
	i = 2
	for rows5.Next() {
		var idOnline, cantidad, total int
		var fecha, cliente, tipo, producto string
		rows5.Scan(&idOnline, &fecha, &cliente, &tipo, &producto, &cantidad, &total)
		f.SetCellValue("Pedidos Online", fmt.Sprintf("A%d", i), idOnline)
		f.SetCellValue("Pedidos Online", fmt.Sprintf("B%d", i), fecha)
		f.SetCellValue("Pedidos Online", fmt.Sprintf("C%d", i), cliente)
		f.SetCellValue("Pedidos Online", fmt.Sprintf("D%d", i), tipo)
		f.SetCellValue("Pedidos Online", fmt.Sprintf("E%d", i), producto)
		f.SetCellValue("Pedidos Online", fmt.Sprintf("F%d", i), cantidad)
		f.SetCellValue("Pedidos Online", fmt.Sprintf("G%d", i), total)
		i++
	}

	f.DeleteSheet("Sheet1")

	ctx.Response.Header.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Response.Header.Set("Content-Disposition", "attachment; filename=reporte_tocke.xlsx")
	f.Write(ctx)
}
