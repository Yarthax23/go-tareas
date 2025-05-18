package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/yarthax23/go-tareas/db"
	"github.com/yarthax23/go-tareas/models"
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

	var ts []models.Tarea
	for rows.Next() {
		var t models.Tarea
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

	var t models.Tarea
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
	var t models.Tarea
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

func insertarTarea(contenido string, w http.ResponseWriter) {
	var id int
	err := db.QueryRow("INSERT INTO tareas (contenido) VALUES ($1) RETURNING id", contenido).Scan(&id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error guardando en la base de datos")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"mensaje": "Tarea recibida",
		"id":      id,
	})
}

func buildUpdateQuery(id int, datos *models.Tarea) (string, []interface{}) {
	var campos []string
	var args []interface{}
	i := 1

	if datos.Contenido != nil {
		campos = append(campos, fmt.Sprintf("contenido = $%d", i))
		args = append(args, *datos.Contenido)
		i++
	}
	if datos.Fecha != nil {
		campos = append(campos, fmt.Sprintf("fecha = $%d", i))
		args = append(args, *datos.Fecha)
		i++
	}
	if datos.Resuelto != nil {
		campos = append(campos, fmt.Sprintf("resuelto = $%d", i))
		args = append(args, *datos.Resuelto)
		i++
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE tareas SET %s WHERE id = $%d", strings.Join(campos, ", "), i)
	return query, args
}

// Recibe el cuerpo de la petición (r.Body) y devuelve map con los campos permitidos
func parseUpdateData(r *http.Request) (*models.Tarea, error) {
	var t models.Tarea
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("error decodificando JSON: %w", err)
	}

	if t.Contenido == nil && t.Fecha == nil && t.Resuelto == nil {
		return nil, fmt.Errorf("no se recibió ningún campo para actualizar")
	}

	return &t, nil
}
