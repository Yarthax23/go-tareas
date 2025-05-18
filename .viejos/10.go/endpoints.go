package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
)

func handlePost(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		var p Post
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(w, http.StatusBadRequest, "Error al leer JSON")
			return
		}
		insertarPost(p.Contenido, w)

	case "text/plain":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Error leyendo texto plano")
			return
		}
		insertarPost(string(body), w)

	default:
		writeError(w, http.StatusUnsupportedMediaType, "Se espera Content-Type: application/json o text/plain")
		return
	}
}

func handleVer(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, contenido, resuelto, fecha FROM tareas ORDER BY id ASC")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error leyendo de DB")
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		var fecha sql.NullTime
		if err := rows.Scan(&p.ID, &p.Contenido, &p.Resuelto, &fecha); err != nil {
			writeError(w, http.StatusInternalServerError, "Error leyendo fila")
			return
		}
		if fecha.Valid {
			p.Fecha = fecha.Time
		}
		posts = append(posts, p)
	}

	writeJSON(w, http.StatusOK, posts)
}

func handleVerPorID(w http.ResponseWriter, r *http.Request) {
	id, err := extraerID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var p Post
	var fecha sql.NullTime
	row := db.QueryRow("SELECT id, contenido, resuelto, fecha FROM tareas WHERE id = $1", id)
	if err = row.Scan(&p.ID, &p.Contenido, &p.Resuelto, &fecha); err != nil {
		writeError(w, http.StatusNotFound, "No se encontró el POST con ese ID")
		return
	}
	if fecha.Valid {
		p.Fecha = fecha.Time
	}
	writeJSON(w, http.StatusOK, p)
}

func handleActualizarPost(w http.ResponseWriter, r *http.Request) {
	id, err := extraerID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// Leer el nuevo contenido (JSON esperado)
	var p Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "Error al leer JSON")
		return
	}

	// Ejecutar actualización
	res, err := db.Exec("UPDATE tareas SET contenido = $1 WHERE id = $2", p.Contenido, id)
	if err != nil || afectadas(res) == 0 {
		writeError(w, http.StatusInternalServerError, "Error actualizando post")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Post actualizado (PUT)",
		"id":      id,
	})
}

func handleEditarPost(w http.ResponseWriter, r *http.Request) {
	id, err := extraerID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	datos, err := parseUpdateData(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Error al leer JSON")
		return
	}

	query, args := buildUpdateQuery(id, datos)
	res, err := db.Exec(query, args...)
	if err != nil || afectadas(res) == 0 {
		writeError(w, http.StatusInternalServerError, "Error al actualizar en la base de datos")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Post editado (PATCH)",
		"id":      id,
	})
}

func handleToggle(w http.ResponseWriter, r *http.Request) {
	id, err := extraerID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var actual string
	err = db.QueryRow("SELECT resuelto FROM tareas WHERE id = $1", id).Scan(&actual)
	if err != nil {
		writeError(w, http.StatusNotFound, "No se encontró el post con ese ID")
		return
	}

	var nuevo string = "SI"
	if actual == "SI" {
		nuevo = "NO"
	}
	// Actualizar la fila con el nuevo valor
	_, err = db.Exec("UPDATE tareas SET resuelto = $1 WHERE id = $2", nuevo, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error actualizando la base de datos")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje":      "Estado cambiado correctamente",
		"id":           id,
		"nuevo_estado": nuevo,
	})
}

func handleBorrarPorID(w http.ResponseWriter, r *http.Request) {
	id, err := extraerID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	res, err := db.Exec("DELETE FROM tareas WHERE id = $1", id)
	if err != nil || afectadas(res) == 0 {
		http.Error(w, "Error al intentar borrar", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Post borrado correctamente",
		"id":      id,
	})
}
