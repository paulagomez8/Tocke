package reportes

import (
    "html/template"
    "log"
    "tockesanfelipe/modules/bases"
    "github.com/valyala/fasthttp"
)

type VentaDia struct {
    Fecha string
    Total int
    Pedidos int
}

type VentaMes struct {
    Mes   string
    Total int
    Pedidos int
}

type ProductoVendido struct {
    Nombre   string
    Cantidad int
    Total    int
}

type ReporteData struct {
    VentasDia      []VentaDia
    VentasMes      []VentaMes
    ProductosTop   []ProductoVendido
}

func VerReportes(ctx *fasthttp.RequestCtx) {
    var data ReporteData

    // Ventas por dia (ultimos 30 dias)
    rowsDia, err := bases.DB.Query(`
        SELECT DATE(fecha) as dia, SUM(total), COUNT(*)
        FROM pedidos
        WHERE fecha >= DATE_SUB(NOW(), INTERVAL 30 DAY)
        GROUP BY dia
        ORDER BY dia DESC
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
