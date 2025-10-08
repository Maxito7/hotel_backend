package domain

import "time"

// TipoHabitacion represents the room type
type TipoHabitacion struct {
	ID               int     `json:"id"`
	Titulo           string  `json:"titulo"`
	Descripcion      string  `json:"descripcion"`
	CapacidadAdultos int     `json:"capacidadAdultos"`
	CapacidadNinhos  int     `json:"capacidadNinhos"`
	CantidadCamas    int     `json:"cantidadCamas"`
	Precio           float64 `json:"precio"`
}

// Habitacion represents a room in the hotel with its type information
type Habitacion struct {
	ID                 int            `json:"id"`
	Nombre             string         `json:"nombre"`
	Numero             string         `json:"numero"`
	Capacidad          int            `json:"capacidad"`
	Estado             string         `json:"estado"`
	DescripcionGeneral string         `json:"descripcionGeneral"`
	TipoHabitacion     TipoHabitacion `json:"tipoHabitacion"`
	MediaID            int            `json:"-"` // El tag "-" hace que este campo se omita en la serialización JSON
}

// HabitacionRepository defines the interface for room data operations
type HabitacionRepository interface {
	// GetAllRooms returns all rooms in the system
	GetAllRooms() ([]Habitacion, error)
	// GetAvailableRooms returns rooms that are available for the given date range
	GetAvailableRooms(fechaEntrada, fechaSalida time.Time) ([]Habitacion, error)
}
