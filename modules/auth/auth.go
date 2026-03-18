package auth

import (
	"crypto/rand"
	"encoding/hex"
	"tockesanfelipe/modules/bases"

	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
)

var sesionToken string

func GenerarToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	sesionToken = hex.EncodeToString(bytes)
	return sesionToken
}

func VerificarSesion(ctx *fasthttp.RequestCtx) bool {
	cookie := string(ctx.Request.Header.Cookie("sesion"))
	return cookie != "" && cookie == sesionToken
}

func Login(ctx *fasthttp.RequestCtx) {
	if string(ctx.Method()) == "GET" {
		tmpl := `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <title>Login - Tocke San Felipe</title>
    <style>
        body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background: #f4f4f4; }
        .box { background: white; padding: 40px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); width: 300px; }
        h2 { text-align: center; color: #333; margin-bottom: 25px; }
        input { width: 100%; padding: 10px; margin-bottom: 15px; border: 1px solid #ddd; border-radius: 5px; font-size: 15px; box-sizing: border-box; }
        .btn { width: 100%; padding: 10px; background-color: #27ae60; color: white; border: none; border-radius: 5px; font-size: 16px; cursor: pointer; }
        .btn:hover { background-color: #219150; }
        .error { color: red; text-align: center; margin-bottom: 10px; font-size: 14px; }
    </style>
</head>
<body>
    <div class="box">
        <h2>Tocke San Felipe</h2>
        <form method="POST" action="/login" autocomplete="off">
            <input type="text" name="usuario" placeholder="Usuario" required>
            <input type="password" name="password" placeholder="Contraseña" required>
            <button class="btn" type="submit">Entrar</button>
        </form>
    </div>
</body>
</html>`
		ctx.SetContentType("text/html")
		ctx.WriteString(tmpl)
		return
	}

	usuario := string(ctx.FormValue("usuario"))
	password := string(ctx.FormValue("password"))

	var hashPassword string
	err := bases.DB.QueryRow("SELECT password FROM usuarios WHERE nombre = ?", usuario).Scan(&hashPassword)
	if err != nil {
		mostrarErrorLogin(ctx)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		mostrarErrorLogin(ctx)
		return
	}

	token := GenerarToken()
	cookie := fasthttp.AcquireCookie()
	cookie.SetKey("sesion")
	cookie.SetValue(token)
	cookie.SetPath("/")
	ctx.Response.Header.SetCookie(cookie)
	fasthttp.ReleaseCookie(cookie)

	ctx.Redirect("/", 302)
}

func mostrarErrorLogin(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	ctx.WriteString(`<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <title>Login - Tocke San Felipe</title>
    <style>
        body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background: #f4f4f4; }
        .box { background: white; padding: 40px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); width: 300px; }
        h2 { text-align: center; color: #333; margin-bottom: 25px; }
        input { width: 100%; padding: 10px; margin-bottom: 15px; border: 1px solid #ddd; border-radius: 5px; font-size: 15px; box-sizing: border-box; }
        .btn { width: 100%; padding: 10px; background-color: #27ae60; color: white; border: none; border-radius: 5px; font-size: 16px; cursor: pointer; }
        .error { color: red; text-align: center; margin-bottom: 10px; font-size: 14px; }
    </style>
</head>
<body>
    <div class="box">
        <h2>Tocke San Felipe</h2>
        <p class="error">Usuario o contraseña incorrectos</p>
        <form method="POST" action="/login">
            <input type="text" name="usuario" placeholder="Usuario" required>
            <input type="password" name="password" placeholder="Contraseña" required>
            <button class="btn" type="submit">Entrar</button>
        </form>
    </div>
</body>
</html>`)
}
