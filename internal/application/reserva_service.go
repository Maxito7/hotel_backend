package application

import (
	"fmt"
	"time"

	"github.com/Maxito7/hotel_backend/internal/domain"
	"github.com/Maxito7/hotel_backend/internal/email"
)

type ReservaService struct {
	reservaRepo           domain.ReservaRepository
	reservaHabitacionRepo domain.ReservaHabitacionRepository
	habitacionRepo        domain.HabitacionRepository
	emailClient           *email.Client
}

// NewReservaService crea una nueva instancia del servicio de reservas
func NewReservaService(
	reservaRepo domain.ReservaRepository,
	reservaHabitacionRepo domain.ReservaHabitacionRepository,
	habitacionRepo domain.HabitacionRepository,
	emailClient *email.Client,
) *ReservaService {
	return &ReservaService{
		reservaRepo:           reservaRepo,
		reservaHabitacionRepo: reservaHabitacionRepo,
		habitacionRepo:        habitacionRepo,
		emailClient:           emailClient,
	}
}

// CreateReserva crea una nueva reserva validando disponibilidad
func (s *ReservaService) CreateReserva(reserva *domain.Reserva) error {
	// Validar que la reserva tenga habitaciones
	if len(reserva.Habitaciones) == 0 {
		return fmt.Errorf("la reserva debe tener al menos una habitación")
	}

	// Validar fechas y disponibilidad de cada habitación
	for i, hab := range reserva.Habitaciones {
		// Validar que fecha de salida sea posterior a fecha de entrada
		if !hab.FechaSalida.After(hab.FechaEntrada) {
			return fmt.Errorf("la fecha de salida debe ser posterior a la fecha de entrada para la habitación %d", hab.HabitacionID)
		}

		// Verificar disponibilidad
		disponible, err := s.reservaHabitacionRepo.VerificarDisponibilidad(
			hab.HabitacionID,
			hab.FechaEntrada,
			hab.FechaSalida,
		)
		if err != nil {
			return fmt.Errorf("error al verificar disponibilidad: %w", err)
		}

		if !disponible {
			return fmt.Errorf("la habitación %d no está disponible para las fechas seleccionadas", hab.HabitacionID)
		}

		// Validar que se haya proporcionado un precio
		if hab.Precio <= 0 {
			return fmt.Errorf("el precio de la habitación %d debe ser mayor a 0", hab.HabitacionID)
		}

		reserva.Habitaciones[i].Precio = hab.Precio
	}

	// Calcular subtotal
	subtotal := 0.0
	for _, hab := range reserva.Habitaciones {
		// Calcular días de estancia
		dias := hab.FechaSalida.Sub(hab.FechaEntrada).Hours() / 24
		if dias < 1 {
			dias = 1
		}
		subtotal += hab.Precio * dias
	}
	reserva.Subtotal = subtotal

	// Si no se especificó descuento, establecerlo en 0
	if reserva.Descuento < 0 {
		reserva.Descuento = 0
	}

	// Validar que el descuento no sea mayor al subtotal
	if reserva.Descuento > reserva.Subtotal {
		return fmt.Errorf("el descuento no puede ser mayor al subtotal")
	}

	// Establecer fecha de confirmación si no se proporcionó
	if reserva.FechaConfirmacion.IsZero() {
		reserva.FechaConfirmacion = time.Now()
	}

	// Establecer estado inicial si no se especificó
	if reserva.Estado == "" {
		reserva.Estado = domain.ReservaPendiente
	}

	// Crear la reserva
	if err := s.reservaRepo.CreateReserva(reserva); err != nil {
		return fmt.Errorf("error al crear reserva: %w", err)
	}

	return nil
}

// GetReservaByID obtiene una reserva por su ID
func (s *ReservaService) GetReservaByID(id int) (*domain.Reserva, error) {
	return s.reservaRepo.GetReservaByID(id)
}

// GetReservasCliente obtiene todas las reservas de un cliente
func (s *ReservaService) GetReservasCliente(clienteID string) ([]domain.Reserva, error) {
	return s.reservaRepo.GetReservasCliente(clienteID)
}

// UpdateReservaEstado actualiza el estado de una reserva
func (s *ReservaService) UpdateReservaEstado(id int, estado domain.EstadoReserva) error {
	// Validar que el estado sea válido
	validEstados := map[domain.EstadoReserva]bool{
		domain.ReservaPendiente:  true,
		domain.ReservaConfirmada: true,
		domain.ReservaCancelada:  true,
		domain.ReservaCompletada: true,
	}

	if !validEstados[estado] {
		return fmt.Errorf("estado de reserva inválido: %s", estado)
	}

	// Obtener la reserva actual
	reserva, err := s.reservaRepo.GetReservaByID(id)
	if err != nil {
		return fmt.Errorf("error al obtener reserva: %w", err)
	}

	// Si se está cancelando, actualizar el estado de las habitaciones
	if estado == domain.ReservaCancelada {
		for _, hab := range reserva.Habitaciones {
			if err := s.reservaHabitacionRepo.UpdateReservaHabitacionEstado(
				id,
				hab.HabitacionID,
				0, // Estado cancelado
			); err != nil {
				return fmt.Errorf("error al cancelar habitación: %w", err)
			}
		}
	}

	return s.reservaRepo.UpdateReservaEstado(id, estado)
}

// CancelarReserva cancela una reserva completa
func (s *ReservaService) CancelarReserva(id int) error {
	return s.UpdateReservaEstado(id, domain.ReservaCancelada)
}

// ConfirmarReserva confirma una reserva pendiente y envía email de confirmación
func (s *ReservaService) ConfirmarReserva(id int) error {
	return s.confirmarReservaInternal(id, true) // true = enviar email
}

// ConfirmarReservaSinEmail confirma una reserva sin enviar email
func (s *ReservaService) ConfirmarReservaSinEmail(id int) error {
	return s.confirmarReservaInternal(id, false) // false = no enviar email
}

// confirmarReservaInternal es el método interno que maneja la confirmación
func (s *ReservaService) confirmarReservaInternal(id int, enviarEmail bool) error {
	// Actualizar estado
	if err := s.UpdateReservaEstado(id, domain.ReservaConfirmada); err != nil {
		return err
	}

	// Solo enviar email si se solicita
	if !enviarEmail {
		return nil
	}

	// Obtener datos completos de la reserva para el email
	reserva, err := s.GetReservaByID(id)
	if err != nil {
		// Log error pero no fallar si el email falla
		fmt.Printf("Error al obtener reserva para email: %v\n", err)
		return nil // La confirmación ya se hizo, solo falló el email
	}

	// Enviar email de confirmación
	if s.emailClient != nil {
		if err := s.enviarEmailConfirmacion(reserva); err != nil {
			// Log error pero no fallar
			fmt.Printf("Error al enviar email de confirmación: %v\n", err)
		}
	}

	return nil
}

// enviarEmailConfirmacion envía el email de confirmación de la reserva
func (s *ReservaService) enviarEmailConfirmacion(reserva *domain.Reserva) error {
	// Preparar información de habitaciones
	habitaciones := make([]email.HabitacionInfo, len(reserva.Habitaciones))
	for i, hab := range reserva.Habitaciones {
		// Verificar que la habitación no sea nil
		if hab.Habitacion == nil {
			return fmt.Errorf("habitación %d no tiene datos completos", hab.HabitacionID)
		}

		noches := int(hab.FechaSalida.Sub(hab.FechaEntrada).Hours() / 24)
		habitaciones[i] = email.HabitacionInfo{
			Nombre:       hab.Habitacion.Nombre,
			Numero:       hab.Habitacion.Numero,
			FechaEntrada: hab.FechaEntrada,
			FechaSalida:  hab.FechaSalida,
			Precio:       hab.Precio,
			Noches:       noches,
		}
	}

	// Preparar información de la reserva
	reservaInfo := email.ReservaInfo{
		ID:                reserva.ID,
		ClienteEmail:      reserva.ClienteID, // El clienteID es el email
		CantidadAdultos:   reserva.CantidadAdultos,
		CantidadNinhos:    reserva.CantidadNinhos,
		FechaConfirmacion: reserva.FechaConfirmacion,
		Subtotal:          reserva.Subtotal,
		Descuento:         reserva.Descuento,
		Total:             reserva.Subtotal - reserva.Descuento,
		Habitaciones:      habitaciones,
	}

	// Enviar email
	return s.emailClient.SendReservaConfirmacion(reservaInfo)
}

// CompletarReserva marca una reserva como completada
func (s *ReservaService) CompletarReserva(id int) error {
	return s.UpdateReservaEstado(id, domain.ReservaCompletada)
}

// VerificarDisponibilidad verifica si una habitación está disponible
func (s *ReservaService) VerificarDisponibilidad(habitacionID int, fechaEntrada, fechaSalida time.Time) (bool, error) {
	if !fechaSalida.After(fechaEntrada) {
		return false, fmt.Errorf("la fecha de salida debe ser posterior a la fecha de entrada")
	}

	return s.reservaHabitacionRepo.VerificarDisponibilidad(habitacionID, fechaEntrada, fechaSalida)
}

// GetReservasEnRango obtiene todas las reservas en un rango de fechas
func (s *ReservaService) GetReservasEnRango(fechaInicio, fechaFin time.Time) ([]domain.ReservaHabitacion, error) {
	if !fechaFin.After(fechaInicio) {
		return nil, fmt.Errorf("la fecha fin debe ser posterior a la fecha inicio")
	}

	return s.reservaHabitacionRepo.GetReservasEnRango(fechaInicio, fechaFin)
}
