package domain

import (
	"time"
)

// ReservaHabitacion representa la relación entre una reserva y una habitación
type ReservaHabitacion struct {
	ReservaID    int         `json:"reservaId"`
	HabitacionID int         `json:"habitacionId"`
	Precio       float64     `json:"precio"`
	FechaEntrada time.Time   `json:"fechaEntrada"`
	FechaSalida  time.Time   `json:"fechaSalida"`
	Estado       int         `json:"estado"` // 1: Activa, 0: Cancelada
	Habitacion   *Habitacion `json:"habitacion,omitempty"`
}

// ReservaHabitacionRepository define las operaciones disponibles con las reservas de habitaciones
type ReservaHabitacionRepository interface {
	// CreateReservaHabitacion crea una nueva reserva de habitación
	CreateReservaHabitacion(reservaHabitacion *ReservaHabitacion) error
	// GetReservaHabitacionesByReservaID obtiene todas las habitaciones de una reserva
	GetReservaHabitacionesByReservaID(reservaID int) ([]ReservaHabitacion, error)
	// UpdateReservaHabitacionEstado actualiza el estado de una reserva de habitación
	UpdateReservaHabitacionEstado(reservaID, habitacionID int, estado int) error
	// VerificarDisponibilidad verifica si una habitación está disponible para las fechas dadas
	VerificarDisponibilidad(habitacionID int, fechaEntrada, fechaSalida time.Time) (bool, error)
	// GetReservasEnRango obtiene todas las reservas activas en un rango de fechas
	GetReservasEnRango(fechaInicio, fechaFin time.Time) ([]ReservaHabitacion, error)
}
