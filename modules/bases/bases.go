package bases

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

var DB *sql.DB

func ConectarDB() {
    var err error
    DB, err = sql.Open("mysql", "root:12345678@tcp(localhost)/Tocke")
    if err != nil {
        log.Fatal("Error al conectar la base de datos:", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatal("No se pudo hacer ping a la base de datos:", err)
    }

    log.Println("Conexión a la base de datos exitosa")
}
