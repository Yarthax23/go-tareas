package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/yarthax23/go-tareas/db"
	m "github.com/yarthax23/go-tareas/models"
	u "github.com/yarthax23/go-tareas/utils"
)

func PostTask(w http.ResponseWriter, r *http.Request) {

	var task m.Tarea
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		u.WriteError(w, http.StatusBadRequest, "Error al leer JSON")
		return
	}
	insertarTarea(*task.Contenido, w)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, contenido, resuelto, fecha FROM tareas ORDER BY id ASC")
	if err != nil {
		u.WriteError(w, http.StatusInternalServerError, "Error leyendo de Database")
		return
	}
	defer rows.Close()

	var tasks []m.Tarea
	for rows.Next() {
		var task m.Tarea
		var fecha sql.NullTime
		if err := rows.Scan(&task.ID, &task.Contenido, &task.Resuelto, &fecha); err != nil {
			u.WriteError(w, http.StatusInternalServerError, "Error leyendo fila")
			return
		}
		if fecha.Valid {
			task.Fecha = &fecha.Time
		}
		tasks = append(tasks, task)
	}

	u.WriteJSON(w, http.StatusOK, tasks)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := u.ExtraerID(r)
	if err != nil {
		u.WriteError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var task m.Tarea
	var fecha sql.NullTime
	row := db.DB.QueryRow("SELECT id, contenido, resuelto, fecha FROM tareas WHERE id = $1", id)
	if err = row.Scan(&task.ID, &task.Contenido, &task.Resuelto, &fecha); err != nil {
		u.WriteError(w, http.StatusNotFound, "No se encontró la TAREA con ese ID")
		return
	}
	if fecha.Valid {
		task.Fecha = &fecha.Time
	}
	u.WriteJSON(w, http.StatusOK, task)
}

func PutTask(w http.ResponseWriter, r *http.Request) {
	id, err := u.ExtraerID(r)
	if err != nil {
		u.WriteError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// Leer el nuevo contenido (JSON esperado)
	var task m.Tarea
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		u.WriteError(w, http.StatusBadRequest, "Error al leer JSON")
		return
	}

	// Validación (con punteros)
	if task.Contenido == nil || task.Resuelto == nil || task.Fecha == nil {
		u.WriteError(w, http.StatusBadRequest, "PUT requiere contenido, resuelto y fecha \n fecha ejemplo: (2025-5-25T12:25:52Z)")
		return
	}

	/*// Validación (sin punteros)
	if strings.TrimSpace(task.Contenido) == "" || task.Fecha.IsZero() {
		u.WriteError(w, http.StatusBadRequest, "PUT requiere contenido, resuelto y fecha")
		return
	}*/

	query := `
		UPDATE tareas
		SET contenido = $1, resuelto = $2, fecha = $3
		WHERE id = $4
	`
	// Actualización completa
	res, err := db.DB.Exec(query, task.Contenido, task.Resuelto, task.Fecha, id)
	if err != nil || u.Afectadas(res) == 0 {
		u.WriteError(w, http.StatusInternalServerError, "Error actualizando tarea")
		return
	}

	u.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Tarea actualizada (PUT)",
		"id":      id,
	})
}

func PatchTask(w http.ResponseWriter, r *http.Request) {
	id, err := u.ExtraerID(r)
	if err != nil {
		u.WriteError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	datos, err := parseUpdateData(r)
	if err != nil {
		u.WriteError(w, http.StatusBadRequest, "Error al leer JSON")
		return
	}

	query, args := buildUpdateQuery(id, datos)
	res, err := db.DB.Exec(query, args...)
	if err != nil || u.Afectadas(res) == 0 {
		u.WriteError(w, http.StatusInternalServerError, "Error al actualizar en la base de datos")
		return
	}

	u.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Tarea editada (PATCH)",
		"id":      id,
	})
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := u.ExtraerID(r)
	if err != nil {
		u.WriteError(w, http.StatusBadRequest, "ID invalido")
		return
	}

	res, err := db.DB.Exec("DELETE FROM tareas WHERE id = $1", id)
	if err != nil || u.Afectadas(res) == 0 {
		u.WriteError(w, http.StatusInternalServerError, "Error al intentar borrar")
		return
	}

	u.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Tarea borrada correctamente",
		"id":      id,
	})
}

func insertarTarea(contenido string, w http.ResponseWriter) {
	var id int
	err := db.DB.QueryRow("INSERT INTO tareas (contenido) VALUES ($1) RETURNING id", contenido).Scan(&id)
	if err != nil {
		u.WriteError(w, http.StatusInternalServerError, "Error guardando en la database")
		return
	}
	u.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"mensaje": "Tarea recibida",
		"id":      id,
	})
}

func buildUpdateQuery(id int, datos *m.Tarea) (string, []interface{}) {
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
func parseUpdateData(r *http.Request) (*m.Tarea, error) {
	var task m.Tarea
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("error decodificando JSON: %w", err)
	}

	if task.Contenido == nil && task.Fecha == nil && task.Resuelto == nil {
		return nil, fmt.Errorf("no se recibió ningún campo para actualizar")
	}

	return &task, nil
}
