package impresora

import (
	"fmt"
	"net"
	"time"
)

type ItemTicket struct {
	Cantidad      int
	Nombre        string
	Modificadores []string
}

const (
	// Cuando tengan la impresora por red, cambiar a "tcp" y poner la IP:puerto
	// modoConexion = "tcp"
	// direccion    = "192.168.1.100:9100"
	modoConexion = "tcp"
	direccion    = "192.168.1.100:9100" // <-- cambiar por la IP real de la impresora
)

func ImprimirTicket(idPedido int, cliente string, items []ItemTicket) error {
	conn, err := net.DialTimeout(modoConexion, direccion, 5*time.Second)
	if err != nil {
		return fmt.Errorf("no se pudo conectar a la impresora: %v", err)
	}
	defer conn.Close()

	conn.Write([]byte{0x1B, 0x40})       // reset
	conn.Write([]byte{0x1B, 0x61, 0x01}) // centrar
	conn.Write([]byte{0x1B, 0x21, 0x30}) // texto grande
	fmt.Fprintln(conn, "TOCKE SAN FELIPE")
	conn.Write([]byte{0x1B, 0x21, 0x00}) // texto normal
	fmt.Fprintf(conn, "- Pedido: %d -\n", idPedido)
	fmt.Fprintln(conn, time.Now().Format("02/01/06 15:04:05"))
	fmt.Fprintln(conn, "----------------------------")
	fmt.Fprintf(conn, "Cliente: %s\n", cliente)
	fmt.Fprintln(conn, "----------------------------")
	conn.Write([]byte{0x1B, 0x61, 0x00}) // alinear izquierda

	for _, item := range items {
		fmt.Fprintf(conn, "%dx  %s\n", item.Cantidad, item.Nombre)
		for _, mod := range item.Modificadores {
			fmt.Fprintf(conn, "    - %s\n", mod)
		}
		fmt.Fprintln(conn, "----------------------------")
	}

	fmt.Fprintln(conn, "\n\n\n")
	conn.Write([]byte{0x1D, 0x56, 0x41, 0x00}) // cortar papel

	return nil
}
