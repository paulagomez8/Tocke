package impresora

import (
	"fmt"
	"os"
	"time"
)

type ItemTicket struct {
	Cantidad      int
	Nombre        string
	Modificadores []string
}

func ImprimirTicket(idPedido int, cliente string, items []ItemTicket) error {
	printer, err := os.OpenFile("/dev/usb/lp0", os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("no se pudo abrir la impresora: %v", err)
	}
	defer printer.Close()

	printer.Write([]byte{0x1B, 0x40})
	printer.Write([]byte{0x1B, 0x61, 0x01})
	printer.Write([]byte{0x1B, 0x21, 0x30})
	fmt.Fprintln(printer, "TOCKE SAN FELIPE")
	printer.Write([]byte{0x1B, 0x21, 0x00})
	fmt.Fprintf(printer, "- Pedido: %d -\n", idPedido)
	fmt.Fprintln(printer, time.Now().Format("02/01/06 15:04:05"))
	fmt.Fprintln(printer, "----------------------------")
	fmt.Fprintf(printer, "Cliente: %s\n", cliente)
	fmt.Fprintln(printer, "----------------------------")
	printer.Write([]byte{0x1B, 0x61, 0x00})

	for _, item := range items {
		fmt.Fprintf(printer, "%d   %s\n", item.Cantidad, item.Nombre)
		for _, mod := range item.Modificadores {
			fmt.Fprintf(printer, "    - %s\n", mod)
		}
		fmt.Fprintln(printer, "----------------------------")
	}

	fmt.Fprintln(printer, "\n\n\n")
	printer.Write([]byte{0x1D, 0x56, 0x41, 0x00})

	return nil
}
