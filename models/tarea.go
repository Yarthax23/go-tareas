package models

import "time"

type TareaResponse struct {
	ID        int       `json:"id" validate="unique"`
	Contenido string    `json:"contenido"`
	Resuelto  bool      `json:"resuelto" validate="boolean"`
	Fecha     time.Time `json:"fecha"`
	UsuarioID int       `json:"usuario_id"`
}

type TareaPOST struct {
	ID        *int       `json:"id" validate="unique"`
	Contenido *string    `json:"contenido" validate:"required,min=2,max=80"`
	Resuelto  *bool      `json:"resuelto" validate:"boolean"` // Opcional, DB default = false
	Fecha     *time.Time `json:"fecha"`                       // Opcional, DB default = CURRENT_TIMESTAMP
	UsuarioID *int       `json:"usuario_id" validate:"required"`
}

type TareaPUT struct {
	ID        *int       `json:"id" validate="unique"`
	Contenido *string    `json:"contenido" validate:"required,min=2,max=80"`
	Resuelto  *bool      `json:"resuelto" validate="required,boolean"`
	Fecha     *time.Time `json:"fecha" validate:"required"`
	UsuarioID *int       `json:"usuario_id" validate:"required"`
}

type TareaPATCH struct {
	ID        *int       `json:"id" validate="unique"`
	Contenido *string    `json:"contenido,omitempty" validate:"omitempty,min=2,max=80"`
	Resuelto  *bool      `json:"resuelto,omitempty" validate:"omitempty,boolean"`
	Fecha     *time.Time `json:"fecha,omitempty" validate:"omitempty"`
	UsuarioID *int       `json:"usuario_id,omitempty" validate:"omitempty"`
}
