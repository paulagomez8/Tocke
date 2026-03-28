package pedidos

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/valyala/fasthttp"
)

const (
	twilioSID   = "TU_ACCOUNT_SID"
	twilioToken = "TU_AUTH_TOKEN"
	twilioFrom  = "whatsapp:+14155238886"
)

func EnviarWhatsapp(ctx *fasthttp.RequestCtx) {
	telefonoRaw := string(ctx.PostArgs().Peek("telefono"))
	mensaje := string(ctx.PostArgs().Peek("mensaje"))

	// Limpiar y agregar prefijo +569 fijo
	re := regexp.MustCompile(`\D`)
	digits := re.ReplaceAllString(telefonoRaw, "")

	// Si no hay exactamente 8 dígitos, igual intentamos enviar (el frontend ya debería validarlo)
	toNumber := "whatsapp:+569" + digits

	apiURL := "https://api.twilio.com/2010-04-01/Accounts/" + twilioSID + "/Messages.json"

	data := url.Values{}
	data.Set("From", twilioFrom)
	data.Set("To", toNumber)
	data.Set("Body", mensaje)

	req, _ := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	req.SetBasicAuth(twilioSID, twilioToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode >= 400 {
		ctx.SetStatusCode(500)
		ctx.WriteString(`{"ok":false}`)
		return
	}

	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(map[string]bool{"ok": true})
}
