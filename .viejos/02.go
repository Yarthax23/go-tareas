//02 post de cliente recibido por servidor

package main

import (
	"fmt"
	"io"
	"net/http"
)

func handlePost_02(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error leyendo el cuerpo", http.StatusInternalServerError)
		return
	}

	fmt.Println("Recibido POST con body:", string(body))
	fmt.Fprintf(w, "¡POST recibido!\n")
}

func main_02() {
	http.HandleFunc("/post", handlePost)
	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
