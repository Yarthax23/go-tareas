package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/yarthax23/go-tareas/db"
	m "github.com/yarthax23/go-tareas/models"
	u "github.com/yarthax23/go-tareas/utils"
)

func PostUsuario(w http.ResponseWriter, r *http.Request) {
	var input m.UsuarioPost
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.WriteError(w, http.StatusBadRequest, "Error al leer JSON")
		return
	}

	// Validar
	if errores := u.ValidarStruct(input); errores != nil {
		u.WriteJSON(w, http.StatusBadRequest, errores)
		return
	}

	// Crear el usuario
	query := `
	INSERT INTO usuarios 
	(nombre, email) VALUES ($1, $2) RETURNING id
	`
	var id int
	if err := db.DB.QueryRow(query, *input.Nombre, *input.Email).Scan(&id); err != nil {
		u.WriteError(w, http.StatusInternalServerError, "Error guardando en la database")
		return
	}

	u.WriteJSON(w, http.StatusOK, input)
	fmt.Fprintf(w, "Hola %s!\nTe doy la bienvenida!\n", *input.Nombre)
}

func GetUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := u.ExtraerID(r)
	if err != nil {
		u.WriteError(w, http.StatusBadRequest, "ID inválido")
	}

	query := `
		SELECT id, nombre, email 
		FROM usuarios
		WHERE id = $1
	`
	row := db.DB.QueryRow(query, id)

	var input m.UsuarioResponse
	if err = row.Scan(&input.ID, &input.Nombre, &input.Email); err != nil {
		u.WriteError(w, http.StatusNotFound, "No se encontró el usuario con ese ID")
		return
	}
	u.WriteJSON(w, http.StatusOK, input)
}

func GetUsuarios(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, nombre, email
		FROM usuarios
		ORDER BY id ASC
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		u.WriteError(w, http.StatusInternalServerError, "Error leyendo Database")
	}
	defer rows.Close()

	var inputs []m.UsuarioResponse
	for rows.Next() {
		var input m.UsuarioResponse
		if err := rows.Scan(&input.ID, &input.Nombre, &input.Email); err != nil {
			u.WriteError(w, http.StatusInternalServerError, "Error leyendo fila")
			return
		}

		inputs = append(inputs, input)

	}
	u.WriteJSON(w, http.StatusOK, inputs)
}

func PutUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := u.ExtraerID(r)
	if err != nil {
		u.WriteError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var input m.UsuarioPut
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.WriteError(w, http.StatusBadRequest, "Error al leer JSON")
		return
	}

	// Validación
	/*if input.Nombre == nil || input.Email == nil {
		u.WriteError(w, http.StatusBadRequest, "PUT requiere nombre y email")
		return
	}*/
	if errores := u.ValidarStruct(input); errores != nil {
		u.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"mensaje": "Error de validación",
			"errores": errores,
		})
		return
	}

	query := `
		UPDATE usuarios
		SET nombre = $1, email = $2
		WHERE id = $3
	`
	res, err := db.DB.Exec(query, input.Nombre, input.Email, id)
	if err != nil || u.Afectadas(res) == 0 {
		u.WriteError(w, http.StatusInternalServerError, "Error actualizando usuario")
		return
	}

	u.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Usuario reemplazado",
		"id":      id,
	})
}

func PatchUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := u.ExtraerID(r)
	if err != nil {
		u.WriteError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// Parseo
	var input m.UsuarioPatch
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.WriteError(w, http.StatusBadRequest, "Error al leer JSON")
		return
	}

	// Validacion
	/*if input.Nombre == nil && input.Email == nil {
		u.WriteError(w, http.StatusBadRequest, "Error al actualizar, escriba nombre o email")
		return
	}*/
	if errores := u.ValidarStruct(input); errores != nil {
		u.WriteJSON(w, http.StatusBadRequest, errores)
		return
	}

	// Armar query y args
	var campos []string
	var args []interface{}
	var i int = 1

	if input.Nombre != nil {
		campos = append(campos, fmt.Sprintf("nombre = $%d", i))
		args = append(args, *input.Nombre)
		i++
	}
	if input.Email != nil {
		campos = append(campos, fmt.Sprintf("email = $%d", i))
		args = append(args, *input.Email)
		i++
	}

	query := fmt.Sprintf("UPDATE usuarios SET %s WHERE id = $%d", strings.Join(campos, ", "), i)
	args = append(args, id)

	res, err := db.DB.Exec(query, args...)
	if err != nil || u.Afectadas(res) == 0 {
		u.WriteError(w, http.StatusInternalServerError, "Error al actualizar la Database")
		return
	}

	u.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Usuario actualizado parcialmente",
		"id":      id,
	})
}

func DeleteUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := u.ExtraerID(r)
	if err != nil {
		u.WriteError(w, http.StatusBadRequest, "ID invalido")
		return
	}

	res, err := db.DB.Exec("DELETE FROM usuarios WHERE id = $1", id)
	if err != nil || u.Afectadas(res) == 0 {
		u.WriteError(w, http.StatusInternalServerError, "Error al intentar borrar")
		return
	}

	u.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Usuario borrado correctamente",
		"id":      id,
	})
}
