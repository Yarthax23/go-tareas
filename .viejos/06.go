//06 Leer content-types text/plain (el que veníamos usando) y application/json (nuevo)
// Hicimos un switch case para reutilizar la función según el content-type
// Hicimos función auxiliar para insertar e imprimir el Post sea como sea

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Post struct {
	Contenido string `json:"contenido"`
}

var db *sql.DB

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		var p Post
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "Error al leer JSON", http.StatusBadRequest)
			return
		}
		insertarPost(p.Contenido, w)

	case "text/plain":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error leyendo texto plano", http.StatusBadRequest)
			return
		}
		insertarPost(string(body), w)

	default:
		http.Error(w, "Se espera Content-Type: application/json o text/plain", http.StatusUnsupportedMediaType)
		return
	}
}

func insertarPost(contenido string, w http.ResponseWriter) {
	var id int
	err := db.QueryRow("INSERT INTO posts (contenido) VALUES ($1) RETURNING id", contenido).Scan(&id)
	if err != nil {
		log.Println("Error insertando en DB:", err)
		http.Error(w, "Error guardando en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Post recibido",
		"id":      id,
	})
}

func handleVer(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, contenido FROM posts ORDER BY id ASC")
	if err != nil {
		log.Println("Error consultando DB:", err)
		http.Error(w, "Error leyendo de DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var contenido string
		err := rows.Scan(&id, &contenido)
		if err != nil {
			http.Error(w, "Error leyendo fila", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "[%d] %s\n", id, contenido)
	}
}

func handleVerPorID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	partes := strings.Split(path, "/")
	idStr := partes[len(partes)-1]

	intID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Te gustaría ver el post con ID: %d\n", intID)

	var id int
	var contenido string
	row := db.QueryRow("SELECT id, contenido FROM posts WHERE id = $1", intID)
	err = row.Scan(&id, &contenido)
	if err != nil {
		http.Error(w, "No se encontró el POST con ese ID", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "[%d] %s\n", id, contenido)
}

func main() {
	var err error

	// Cadena de conexión a la base de datos
	connStr := "user=golang password=clave123 dbname=todos host=localhost sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error abriendo conexión con PostgreSQL:", err)
	}
	defer db.Close()

	// Verificar conexión
	err = db.Ping()
	if err != nil {
		log.Fatal("No se pudo conectar a la DB:", err)
	}

	fmt.Println("✅ Conectado a la base de datos PostgreSQL")
	fmt.Println("🌐 Servidor corriendo en http://localhost:8080")

	http.HandleFunc("/post", handlePost)
	http.HandleFunc("/ver", handleVer)
	http.HandleFunc("/ver/{id}", handleVerPorID)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
