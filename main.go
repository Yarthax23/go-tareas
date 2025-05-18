//11 Usar HTML React

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/yarthax23/go-tareas/db"
)

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	// TAREAS Rutas básicas con métodos definidos
	r.HandleFunc("/tareas", handleGetTasks).Methods("GET")
	r.HandleFunc("/tareas/{id}", handleGetTask).Methods("GET")
	r.HandleFunc("/tareas", handlePostTask).Methods("POST")
	r.HandleFunc("/tareas/{id}", handlePutTask).Methods("PUT")
	r.HandleFunc("/tareas/{id}", handlePatchTask).Methods("PATCH")
	r.HandleFunc("/tareas/{id}", handleDeleteTask).Methods("DELETE")

	/* Ejemplos de uso
	curl localhost:8080/tareas
	curl localhost:8080/tareas/2				--formato json
	curl localhost:8080/tareas/2 | jq			--formato coloreado
	curl localhost:8080/tareas -X POST \
		-H "Content-Type: application/json" \
		-d '{"contenido" : "Ejemplo"}'
	curl localhost:8080/tareas/2 -X PATCH \
		-H "Content-Type: application/json" \
		-d '{"contenido" : "Ejemplo"}'
	curl localhost:8080/tareas/2 -X DELETE
	*/

	/* USUARIOS

	GET /usuarios → listar todos

	GET /usuarios/{id} → ver uno

	POST /usuarios → crear nuevo

	PUT /usuarios/{id} → reemplazar

	DELETE /usuarios/{id} → borrar

	Tareas por usuario

	GET /usuarios/{id}/tareas → ver tareas de un usuario

	POST /usuarios/{id}/tareas → crear tarea para un usuario
	*/
	// Para debuggear
	/*
		fmt.Println("Rutas registradas:")
		r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			t, _ := route.GetPathTemplate()
			fmt.Println(t)
			return nil
		})
	*/
	//

	log.Println("🌐 Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
