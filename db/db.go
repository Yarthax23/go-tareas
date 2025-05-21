package db

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := "user=golang password=clave123 dbname=todos host=localhost sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error abriendo conexión con PostgreSQL:", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("No se pudo conectar a la DB:", err)
	}
	log.Println("✅ Conectado a la base de datos 'todos', PostgreSQL")
}
