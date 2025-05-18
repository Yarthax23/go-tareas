package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// --- FUNCIONES AUXILIARES --- //

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func insertarPost(contenido string, w http.ResponseWriter) {
	var id int
	err := db.QueryRow("INSERT INTO tareas (contenido) VALUES ($1) RETURNING id", contenido).Scan(&id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error guardando en la base de datos")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"mensaje": "Post recibido",
		"id":      id,
	})
}

func afectadas(res sql.Result) int64 {
	n, _ := res.RowsAffected()
	return n
}

func extraerID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["id"])
}

func buildUpdateQuery(id int, datos map[string]string) (string, []interface{}) {
	var campos []string
	var args []interface{}
	i := 1

	for _, campo := range []string{"contenido", "fecha", "resuelto"} {
		if val, ok := datos[campo]; ok {
			campos = append(campos, fmt.Sprintf("%s = $%d", campo, i))
			args = append(args, val)
			i++
		}
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE tareas SET %s WHERE id = $%d", strings.Join(campos, ", "), i)
	return query, args
}

// Recibe el cuerpo de la petición (r.Body) y devuelve map con los campos permitidos
func parseUpdateData(r *http.Request) (map[string]string, error) {
	var datos map[string]string
	err := json.NewDecoder(r.Body).Decode(&datos)
	if err != nil {
		return nil, fmt.Errorf("error decodificando JSON: %w", err)
	}

	// Filtrar solo los campos permitidos para actualizar
	camposPermitidos := []string{"contenido", "fecha", "resuelto"}
	datosFiltrados := make(map[string]string)
	for _, campo := range camposPermitidos {
		if val, ok := datos[campo]; ok {
			datosFiltrados[campo] = val
		}
	}

	if len(datosFiltrados) == 0 {
		return nil, fmt.Errorf("no se recibió ningún campo para actualizar")
	}

	// Validar el campo 'resuelto' si está presente
	if resuelto, ok := datosFiltrados["resuelto"]; ok {
		if resuelto != "SI" && resuelto != "NO" {
			return nil, fmt.Errorf("valor inválido para 'resuelto': debe ser 'SI' o 'NO'")
		}
	}

	return datosFiltrados, nil
}
