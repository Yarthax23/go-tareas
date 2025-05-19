package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// --- FUNCIONES AUXILIARES --- //

var modoDesarrollo bool = false

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if !modoDesarrollo {
		json.NewEncoder(w).Encode(data)
	} else {
		jsonData, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "Error generando JSON")
			return
		}
		w.Write(jsonData)
	}
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

func Afectadas(res sql.Result) int64 {
	n, _ := res.RowsAffected()
	return n
}

func ExtraerID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["id"])
}
