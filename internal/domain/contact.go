package domain

import "time"

type EstadoFormulario string

const (
	EstadoNuevo      EstadoFormulario = "Nuevo"
	EstadoEnProceso  EstadoFormulario = "EnProceso"
	EstadoRespondido EstadoFormulario = "Respondido"
	EstadoCerrado    EstadoFormulario = "Cerrado"
)

type Contact struct {
	ID             int64            `db:"form_id" json:"id"`
	UsuarioID      *int64           `db:"user_id" json:"usuarioId,omitempty"`
	Nombre         string           `db:"name" json:"nombre"`
	Email          string           `db:"email" json:"email"`
	Telefono       *string          `db:"phone" json:"telefono,omitempty"`
	Mensaje        *string          `db:"message" json:"mensaje,omitempty"`
	Estado         EstadoFormulario `db:"status" json:"estado"`
	FechaEnvio     time.Time        `db:"sent_date" json:"fechaEnvio"`
	FechaRespuesta *time.Time       `db:"response_date" json:"fechaRespuesta,omitempty"`
}

type CreateContactRequest struct {
	Nombre   string  `json:"nombre" validate:"required,min=3,max=100"`
	Email    string  `json:"email" validate:"required,email,max=150"`
	Telefono *string `json:"telefono" validate:"omitempty,min=5,max=20"`
	Mensaje  *string `json:"mensaje" validate:"omitempty,min=1"`
}

type UpdateEstadoRequest struct {
	Estado         EstadoFormulario `json:"estado" validate:"required,oneof=Nuevo EnProceso Respondido Cerrado"`
	FechaRespuesta *time.Time       `json:"fechaRespuesta,omitempty"`
}
