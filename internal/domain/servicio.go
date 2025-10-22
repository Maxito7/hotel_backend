package domain

// Servicio representa un servicio del hotel
type Servicio struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// ServicioRepository define la interfaz para operaciones de datos de servicios
type ServicioRepository interface {
	// GetAllServices retorna todos los servicios disponibles
	GetAllServices() ([]Servicio, error)
}
