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

var modoDesarrollo bool = false

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if modoDesarrollo {
		jsonData, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Error generando JSON")
			return
		}
		w.Write(jsonData)
	} else {
		json.NewEncoder(w).Encode(data)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
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

func afectadas(res sql.Result) int64 {
	n, _ := res.RowsAffected()
	return n
}

func extraerID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["id"])
}

func buildUpdateQuery(id int, datos *Tarea) (string, []interface{}) {
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
func parseUpdateData(r *http.Request) (*Tarea, error) {
	var t Tarea
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("error decodificando JSON: %w", err)
	}

	if t.Contenido == nil && t.Fecha == nil && t.Resuelto == nil {
		return nil, fmt.Errorf("no se recibió ningún campo para actualizar")
	}

	return &t, nil
	/*
	   // Filtrar solo los campos permitidos para actualizar
	   camposPermitidos := []string{"contenido", "fecha", "resuelto"}
	   //datosFiltrados := make(map[string]string) // toma todo como string
	   datosFiltrados := make(map[string]interface{})

	   	for _, campo := range camposPermitidos {
	   		if val, ok := datos[campo]; ok {
	   			switch campo {
	   			case "resuelto":
	   				boolVal, ok := val.(bool)
	   				if !ok {
	   					return nil, fmt.Errorf("el campo 'resuelto' es booleano (true/false)")
	   				}
	   				datosFiltrados[campo] = boolVal
	   			default:
	   				strVal, ok := val.(string)
	   				if !ok {
	   					return nil, fmt.Errorf("el campo '%s' debe ser texto", campo)
	   				}
	   				datosFiltrados[campo] = strVal
	   			}
	   		}
	   	}

	   	if len(datosFiltrados) == 0 {
	   		return nil, fmt.Errorf("no se recibió ningún campo para actualizar")
	   	}

	   return datosFiltrados, nil
	*/
}
