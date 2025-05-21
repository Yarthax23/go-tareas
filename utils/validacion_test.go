//Unit Tests

package utils_test

import (
	"testing"

	"github.com/yarthax23/go-tareas/models"
	"github.com/yarthax23/go-tareas/utils"
)

func TestValidarStruct_UsuarioValido(t *testing.T) {
	u := models.UsuarioPost{
		Nombre: strPtr("Juan"),
		Email:  strPtr("juan@example.com"),
	}

	errores := utils.ValidarStruct(u)
	if errores != nil {
		t.Errorf("Esperaba nil, pero recibió errores: %v", errores)
	}
}

func TestValidarStruct_UsuarioInvalido(t *testing.T) {
	u := models.UsuarioPost{
		Nombre: strPtr(""),
		Email:  strPtr("no-es-email"),
	}

	errores := utils.ValidarStruct(u)
	if errores == nil || len(errores) == 0 {
		t.Error("Esperaba errores de validación, pero recibió nil")
	}

	if errores["Nombre"] == "" || errores["Email"] == "" {
		t.Error("Faltan errores esperados en Nombre o Email")
	}
}

func strPtr(s string) *string {
	return &s
}
