//11 Agregar otro recurso

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/yarthax23/go-tareas/db"
	h "github.com/yarthax23/go-tareas/handlers"
)

func main() {
	db.InitDB()
	defer db.DB.Close()
	r := mux.NewRouter()

	// TAREAS Rutas básicas con métodos definidos
	r.HandleFunc("/tareas", h.GetTasks).Methods("GET")
	r.HandleFunc("/tareas", h.PostTask).Methods("POST")
	r.HandleFunc("/tareas/{id}", h.GetTask).Methods("GET")
	r.HandleFunc("/tareas/{id}", h.PutTask).Methods("PUT")
	r.HandleFunc("/tareas/{id}", h.PatchTask).Methods("PATCH")
	r.HandleFunc("/tareas/{id}", h.DeleteTask).Methods("DELETE")

	// USUARIOS Rutas básicas con métodos definidos
	r.HandleFunc("/usuarios", h.GetUsuarios).Methods("GET")
	r.HandleFunc("/usuarios", h.PostUsuario).Methods("POST")
	r.HandleFunc("/usuarios/{id}", h.GetUsuario).Methods("GET")
	r.HandleFunc("/usuarios/{id}", h.PutUsuario).Methods("PUT")
	r.HandleFunc("/usuarios/{id}", h.PatchUsuario).Methods("PATCH")
	r.HandleFunc("/usuarios/{id}", h.DeleteUsuario).Methods("DELETE")
	/* Ejemplos de uso
	curl localhost:8080/tareas
	curl localhost:8080/tareas/2				--formato json
	curl localhost:8080/tareas/2 | jq			--formato coloreado
	curl localhost:8080/tareas -X POST \
		-d '{"contenido" : "Ejemplo"}'
	curl localhost:8080/tareas/2 -X PATCH \
		-d '{"contenido" : "Ejemplo"
		"resuelto" :  true, "fecha" : "2025-05-25T12:30:30Z"}'
	curl localhost:8080/tareas/2 -X DELETE
	*/

	/* Tareas por usuario ???

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
