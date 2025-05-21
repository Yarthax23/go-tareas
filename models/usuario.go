package models

type UsuarioPost struct {
	ID     *int    `json:"id" validate="unique"`
	Nombre *string `json:"nombre" validate:"required,min=2,max=50"`
	Email  *string `json:"email" validate:"required,email"`
}

type UsuarioPut = UsuarioPost

type UsuarioPatch struct {
	ID     *int    `json:"id" validate="unique"`
	Nombre *string `json:"nombre,omitempty" validate:"omitempty,min=2,max=50"`
	Email  *string `json:"email,omitempty" validate:"omitempty,email"`
}

type UsuarioResponse struct {
	ID     int    `json:"id" validate="unique"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}
