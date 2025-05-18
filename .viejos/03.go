//03 post de cliente recibido y guardado en memoria (no data base aún)
// si no cerrás el servidor podés ver el dato, si cerrás el servidor perdés el dato
// mutex (sync) permite ver de una sola solicitud

package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

var (
	todos []string   // Lista en memoria
	mutex sync.Mutex // Para que sea seguro si acceden varios al mismo tiempo
)

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error leyendo el cuerpo", http.StatusInternalServerError)
		return
	}

	texto := string(body)

	mutex.Lock()
	todos = append(todos, texto)
	mutex.Unlock()

	fmt.Println("POST recibido y guardado:", texto)
	fmt.Fprintln(w, "Guardado en memoria.")
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, t := range todos {
		fmt.Fprintf(w, "%d: %s\n", i+1, t)
	}
}

func main() {
	http.HandleFunc("/post", handlePost)
	http.HandleFunc("/ver", handleGet) // <- Este es el nuevo endpoint

	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
