//04 conectar servidor con la base de datos en el servidor PSQL y handle POSTs y GETs a la/desde la DB (Base de Datos)

package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error leyendo el cuerpo", http.StatusInternalServerError)
		return
	}

	contenido := string(body)

	// Insertar en la base de datos
	_, err = db.Exec("INSERT INTO posts (contenido) VALUES ($1)", contenido)
	if err != nil {
		log.Println("Error insertando en DB:", err)
		http.Error(w, "Error guardando en DB", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "¡POST recibido y guardado!\n")
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

func main() {
	var err error

	// Cadena de conexión a la base de datos
	connStr := "user=golang password=clave123 dbname=todos host=localhost sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error abriendo conexión con PostgreSQL:", err)
	}

	// Verificar conexión
	err = db.Ping()
	if err != nil {
		log.Fatal("No se pudo conectar a la DB:", err)
	}

	fmt.Println("✅ Conectado a la base de datos PostgreSQL")
	fmt.Println("🌐 Servidor corriendo en http://localhost:8080")

	http.HandleFunc("/post", handlePost)
	http.HandleFunc("/ver", handleVer)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
