//10 Modularizar código
// Reutilizar funciones
// Pasar a usar mux complementando a net/http
//	>> nos deshicimos de validarMetodo
//		(r.Method("x", "y")) deja pasar sólo esos métodos
//	>> nos deshicimos de extraerID cambió la implementación

package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Post struct {
	ID        int       `json:"id"`
	Contenido string    `json:"contenido"`
	Resuelto  string    `json:"resuelto"`
	Fecha     time.Time `json:"fecha"`
}

func main() {
	initDB()

	r := mux.NewRouter()

	// Rutas básicas con métodos definidos
	r.HandleFunc("/ver", handleVer).Methods("GET")
	r.HandleFunc("/ver/{id}", handleVerPorID).Methods("GET")
	r.HandleFunc("/post", handlePost).Methods("POST")
	r.HandleFunc("/post/{id}", handleBorrarPorID).Methods("DELETE")
	r.HandleFunc("/post/{id}", handleActualizarPost).Methods("PUT")
	r.HandleFunc("/post/{id}", handleEditarPost).Methods("PATCH")
	r.HandleFunc("/post/{id}/toggle", handleToggle).Methods("POST")

	// Para debuggear
	/*
		fmt.Println("Rutas registradas:")
		r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			t, _ := route.GetPathTemplate()
			fmt.Println(t)
			return nil
		})
	*/
	//

	log.Println("🌐 Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// db.go
var db *sql.DB

func initDB() {
	var err error
	connStr := "user=golang password=clave123 dbname=todos host=localhost sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error abriendo conexión con PostgreSQL:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("No se pudo conectar a la DB:", err)
	}
	log.Println("✅ Conectado a la base de datos PostgreSQL")
}
