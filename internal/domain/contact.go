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
	ID             int64            `db:"formularioId" json:"id"`
	UsuarioID      *int64           `db:"usuarioId" json:"usuarioId,omitempty"`
	Nombre         string           `db:"nombre" json:"nombre"`
	Email          string           `db:"email" json:"email"`
	Telefono       *string          `db:"telefono" json:"telefono,omitempty"`
	Mensaje        *string          `db:"mensaje" json:"mensaje,omitempty"`
	Estado         EstadoFormulario `db:"estado" json:"estado"`
	FechaEnvio     time.Time        `db:"fechaEnvio" json:"fechaEnvio"`
	FechaRespuesta *time.Time       `db:"fechaRespuesta" json:"fechaRespuesta,omitempty"`
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
