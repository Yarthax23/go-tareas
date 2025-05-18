//01 Ver servidor local, puerto y ruta, prints en servidor y en cliente (curl)

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main_01() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s -> 204 No Content\n", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNoContent) // 204
	})

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s -> 200 OK\n", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusOK) // 200
		fmt.Fprintln(w, "Hola!")
	})

	// handler para cualquier ruta no definida
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ping" && r.URL.Path != "/hello" {
			log.Printf("Request: %s %s -> 404 Not Found\n", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	})

	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
