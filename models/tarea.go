type Tarea struct {
	ID        *int       `json:"id"`
	Contenido *string    `json:"contenido"`
	Resuelto  *bool      `json:"resuelto"`
	Fecha     *time.Time `json:"fecha"`
	UsuarioID *int       `json:"usuario_id`
}

/*type TareaPut struct {
	Contenido string    `json:"contenido"`
	Resuelto  bool      `json:"resuelto"`
	Fecha     time.Time `json:"fecha"`
	UsuarioID int       `json:"usuario_id`
}*/
