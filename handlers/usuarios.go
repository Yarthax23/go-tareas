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
	var user m.Usuario
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		u.WriteError(w, http.StatusBadRequest, "Error al leer JSON")
		return
	}

	query := `
	INSERT INTO usuarios 
	(nombre, email) VALUES ($1, $2) RETURNING id
	`
	var id int
	if err := db.DB.QueryRow(query, *user.Nombre, *user.Email).Scan(&id); err != nil {
		u.WriteError(w, http.StatusInternalServerError, "Error guardando en la database")
		return
	}

	fmt.Fprintf(w, "Hola %s!\nTe doy la bienvenida!\n", *user.Nombre)
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

	var user m.Usuario
	if err = row.Scan(&user.ID, &user.Nombre, &user.Email); err != nil {
		u.WriteError(w, http.StatusNotFound, "No se encontró el usuario con ese ID")
		return
	}
	u.WriteJSON(w, http.StatusOK, user)
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

	var users []m.Usuario
	for rows.Next() {
		var user m.Usuario
		if err := rows.Scan(&user.ID, &user.Nombre, &user.Email); err != nil {
			u.WriteError(w, http.StatusInternalServerError, "Error leyendo fila")
			return
		}

		users = append(users, user)

	}
	u.WriteJSON(w, http.StatusOK, users)
}

func PutUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := u.ExtraerID(r)
	if err != nil {
		u.WriteError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var user m.Usuario
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		u.WriteError(w, http.StatusBadRequest, "Error al leer JSON")
		return
	}

	// Validación
	if user.Nombre == nil || user.Email == nil {
		u.WriteError(w, http.StatusBadRequest, "PUT requiere nombre y email")
		return
	}

	query := `
		UPDATE usuarios
		SET nombre = $1, email = $2
		WHERE id = $3
	`
	res, err := db.DB.Exec(query, user.Nombre, user.Email, id)
	if err != nil || u.Afectadas(res) == 0 {
		u.WriteError(w, http.StatusInternalServerError, "Error actualizando usuario")
		return
	}

	u.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Usuario actualizado (PUT)",
		"id":      id,
	})
}

func PatchUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := u.ExtraerID(r)
	if err != nil {
		u.WriteError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	//Parseo
	var user m.Usuario
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		u.WriteError(w, http.StatusBadRequest, "Error al leer JSON")
		return
	}
	if user.Nombre == nil && user.Email == nil {
		u.WriteError(w, http.StatusBadRequest, "Error al actualizar, escriba nombre o email")
		return
	}

	// Armar query y args
	var campos []string
	var args []interface{}
	var i int = 1

	if user.Nombre != nil {
		campos = append(campos, fmt.Sprintf("nombre = $%d", i))
		args = append(args, *user.Nombre)
		i++
	}
	if user.Email != nil {
		campos = append(campos, fmt.Sprintf("email = $%d", i))
		args = append(args, *user.Email)
		i++
	}

	query := fmt.Sprintf("UPDATE usuarios SET %s WHERE id = $%d", strings.Join(campos, ", "), id)
	args = append(args, id)

	res, err := db.DB.Exec(query, args...)
	if err != nil || u.Afectadas(res) == 0 {
		u.WriteError(w, http.StatusInternalServerError, "Error al actualizar la Database")
		return
	}

	u.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"mensaje": "Usuario actualizado (PUT)",
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
