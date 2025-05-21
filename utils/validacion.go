package utils

import "github.com/go-playground/validator/v10"

//err := validate.Struct(mystruct)
//validationErrors := err.(validator.ValidationErrors)

//validate := validate.New(validator.WithRequiredStructEnabled())

// var validate *validator.Validate
var Validator = validator.New()

func ValidarStruct(s interface{}) map[string]string {
	err := Validator.Struct(s)
	if err == nil {
		return nil
	}

	errores := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		campo := e.Field()
		switch e.Tag() {
		case "required":
			errores[campo] = "es obligatorio"
		case "email":
			errores[campo] = "no tiene formato de email válido"
		case "min":
			errores[campo] = "es demasiado corto"
		case "max":
			errores[campo] = "es demasiado largo"
		default:
			errores[campo] = "no es válido"
		}
	}
	return errores
}
