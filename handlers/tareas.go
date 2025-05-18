package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func handlePostTask(w http.ResponseWriter, r *http.Request) {

	var t Tarea
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, "Error al leer JSON")
		return
	}
	insertarTarea(*t.Contenido, w)
}

func handleGetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, contenido, resuelto, fecha FROM tareas ORDER BY id ASC")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error leyendo de DB")
		return
	}
	defer rows.Close()

	var ts []Tarea
	for rows.Next() {
		var t Tarea
		var fecha sql.NullTime
		if err := rows.Scan(&t.ID, &t.Contenido, &t.Resuelto, &fecha); err != nil {
			writeError(w, http.StatusInternalServerError, "Error leyendo fila")
			return
		}
		if fecha.Valid {
			t.Fecha = &fecha.Time
		}
		ts = append(ts, t)
	}

	writeJSON(w, http.StatusOK, ts)
}

func handleGetTask(w http.ResponseWriter, r *http.Request) {
	id, err := extraerID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var t Tarea
	var fecha sql.NullTime
	row := db.QueryRow("SELECT id, contenido, resuelto, fecha FROM tareas WHERE id = $1", id)
	if err = row.Scan(&t.ID, &t.Contenido, &t.Resuelto, &fecha); err != nil {
		writeError(w, http.StatusNotFound, "No se encontró la TAREA con ese ID")
		return
	}
	if fecha.Valid {
		t.Fecha = &fecha.Time
	}
	writeJSON(w, http.StatusOK, t)
}

func handlePutTask(w http.ResponseWriter, r *http.Request) {
	id, err := extraerID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// Leer el nuevo contenido (JSON esperado)
	var t Tarea
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, "Error al leer JSON")
		return
	}

	// Validación (con punteros)
	if t.Contenido == nil || t.Resuelto == nil || t.Fecha == nil {
		writeError(w, http.StatusBadRequest, "PUT requiere contenido, resuelto y fecha \n fecha ejemplo: (2025-5-25T12:25:52Z)")
		return
	}

	/*// Validación (sin punteros)
	if strings.TrimSpace(t.Contenido) == "" || t.Fecha.IsZero() {
		writeError(w, http.StatusBadRequest, "PUT requiere contenido, resuelto y fecha")
		return
	}*/

	query := `
		UPDATE tareas
		SET contenido = $1, resuelto = $2, fecha = $3
		WHERE id = $4
	`
	// Actualización completa
	res, err := db.Exec(query, t.Contenido, t.Resuelto, t.Fecha, id)
	if err != nil || afectadas(res) == 0 {
		writeError(w, http.StatusInternalServerError, "Error actualizando tarea")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Tarea actualizada (PUT)",
		"id":      id,
	})
}

func handlePatchTask(w http.ResponseWriter, r *http.Request) {
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
		"mensaje": "Tarea editada (PATCH)",
		"id":      id,
	})
}

func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := extraerID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	res, err := db.Exec("DELETE FROM tareas WHERE id = $1", id)
	if err != nil || afectadas(res) == 0 {
		writeError(w, http.StatusInternalServerError, "Error al intentar borrar")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Tarea borrada correctamente",
		"id":      id,
	})
}
