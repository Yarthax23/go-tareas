//08 Agregamos columna listo por DEFAULT 'NO' con endpoint /toggle
// Agregamos funcion actualizar y editar
// Agregamos funcion auxiliar validarMetoto y extraerIDDesdeURL
// Hay errores a corregir en PUT y PATCH por la nueva columna, creo.

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

func validarMetodo(w http.ResponseWriter, r *http.Request, metodo string) bool {
	if r.Method != metodo {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func extraerIDDesdeURL(r *http.Request) (int, error) {
	partes := strings.Split(r.URL.Path, "/")
	idStr := partes[len(partes)-1]
	return strconv.Atoi(idStr)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if !validarMetodo(w, r, http.MethodPost) {
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

func handleActualizarPost(w http.ResponseWriter, r *http.Request) {
	if !validarMetodo(w, r, http.MethodPut) {
		return
	}

	intID, err := extraerIDDesdeURL(r)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Leer el nuevo contenido (JSON esperado)
	var p Post
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Error al leer JSON", http.StatusBadRequest)
		return
	}

	// Ejecutar actualización
	res, err := db.Exec("UPDATE posts SET contenido = $1 WHERE id = $2", p.Contenido, intID)
	if err != nil {
		log.Println("Error actualizando:", err)
		http.Error(w, "Error actualizando post", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "No se encontró el post con ese ID", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Post actualizado (PUT)",
		"id":      intID,
	})
}

func handleEditarPost(w http.ResponseWriter, r *http.Request) {
	if !validarMetodo(w, r, http.MethodPatch) {
		return
	}

	intID, err := extraerIDDesdeURL(r)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var datosParciales map[string]string
	err = json.NewDecoder(r.Body).Decode(&datosParciales)
	if err != nil {
		http.Error(w, "Error al leer JSON", http.StatusBadRequest)
		return
	}

	contenido, ok := datosParciales["contenido"]
	if !ok {
		http.Error(w, "No se recibió 'contenido'", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("UPDATE posts SET contenido = $1 WHERE id = $2", contenido, intID)
	if err != nil {
		log.Println("Error actualizando:", err)
		http.Error(w, "Error actualizando post", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "No se encontró el post con ese ID", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Post editado (PATCH)",
		"id":      intID,
	})
}

func handleToggle(w http.ResponseWriter, r *http.Request) {
	if !validarMetodo(w, r, http.MethodPost) {
		return
	}

	// Extraer ID de la URL
	intID, err := extraerIDDesdeURL(r)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Leer el valor actual de "listo"
	var actual string
	err = db.QueryRow("SELECT listo FROM posts WHERE id = $1", intID).Scan(&actual)
	if err != nil {
		http.Error(w, "No se encontró el post con ese ID", http.StatusNotFound)
		return
	}

	// Calcular el nuevo valor
	var nuevo string
	if actual == "SI" {
		nuevo = "NO"
	} else {
		nuevo = "SI"
	}

	// Actualizar la fila con el nuevo valor
	_, err = db.Exec("UPDATE posts SET listo = $1 WHERE id = $2", nuevo, intID)
	if err != nil {
		http.Error(w, "Error actualizando la base de datos", http.StatusInternalServerError)
		return
	}

	// Devolver confirmación
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje":      "Estado cambiado correctamente",
		"id":           intID,
		"nuevo_estado": nuevo,
	})
}

func handleVer(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, contenido, listo FROM posts ORDER BY id ASC")
	if err != nil {
		log.Println("Error consultando DB:", err)
		http.Error(w, "Error leyendo de DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var contenido, listo string
		err := rows.Scan(&id, &contenido, &listo)
		if err != nil {
			http.Error(w, "Error leyendo fila", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "[%d] %s %s\n", id, listo, contenido)
	}
}

func handleVerPorID(w http.ResponseWriter, r *http.Request) {
	if !validarMetodo(w, r, http.MethodGet) {
		return
	}

	intID, err := extraerIDDesdeURL(r)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Te gustaría ver el post con ID: %d\n", intID)

	var id int
	var contenido, listo string
	row := db.QueryRow("SELECT id, contenido, listo FROM posts WHERE id = $1", intID)
	err = row.Scan(&id, &contenido, &listo)
	if err != nil {
		http.Error(w, "No se encontró el POST con ese ID", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "[%d] %s %s\n", id, listo, contenido)
}

func handleBorrarPorID(w http.ResponseWriter, r *http.Request) {
	if !validarMetodo(w, r, http.MethodDelete) {
		return
	}

	intID, err := extraerIDDesdeURL(r)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Te gustaría borrar el post con ID: %d\n", intID)

	resultado, err := db.Exec("DELETE FROM posts WHERE id = $1", intID)
	if err != nil {
		http.Error(w, "Error al intentar borrar", http.StatusInternalServerError)
		return
	}

	filasAfectadas, err := resultado.RowsAffected()
	if err != nil {
		http.Error(w, "No se pudo verificar si se borró el post", http.StatusInternalServerError)
		return
	}

	if filasAfectadas == 0 {
		http.Error(w, "No se encontró ningún post con ese ID", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Post borrado correctamente",
		"id":      intID,
	})
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

	http.HandleFunc("/post", handlePost)                  // POST
	http.HandleFunc("/ver", handleVer)                    // GET
	http.HandleFunc("/ver/", handleVerPorID)              // GET
	http.HandleFunc("/borrar/", handleBorrarPorID)        // DELETE
	http.HandleFunc("/actualizar/", handleActualizarPost) // PUT
	http.HandleFunc("/editar/", handleEditarPost)         // PATCH
	http.HandleFunc("/toggle/", handleToggle)             // POST
	log.Fatal(http.ListenAndServe(":8080", nil))
}
