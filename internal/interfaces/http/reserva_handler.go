package http

import (
	"strconv"
	"time"

	"github.com/Maxito7/hotel_backend/internal/application"
	"github.com/Maxito7/hotel_backend/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type ReservaHandler struct {
	service *application.ReservaService
}

// NewReservaHandler crea una nueva instancia del handler de reservas
func NewReservaHandler(service *application.ReservaService) *ReservaHandler {
	return &ReservaHandler{
		service: service,
	}
}

// CreateReservaRequest representa la petición para crear una reserva
type CreateReservaRequest struct {
	CantidadAdultos int                       `json:"cantidadAdultos"`
	CantidadNinhos  int                       `json:"cantidadNinhos"`
	ClienteID       string                    `json:"clienteId"`
	Descuento       float64                   `json:"descuento"`
	Habitaciones    []CreateHabitacionReserva `json:"habitaciones"`
}

// CreateHabitacionReserva representa una habitación a reservar
type CreateHabitacionReserva struct {
	HabitacionID int     `json:"habitacionId"`
	Precio       float64 `json:"precio"`
	FechaEntrada string  `json:"fechaEntrada"` // Formato: YYYY-MM-DD
	FechaSalida  string  `json:"fechaSalida"`  // Formato: YYYY-MM-DD
}

// UpdateEstadoRequest representa la petición para actualizar el estado de una reserva
type UpdateEstadoRequest struct {
	Estado string `json:"estado"`
}

// VerificarDisponibilidadRequest representa la petición para verificar disponibilidad
type VerificarDisponibilidadRequest struct {
	HabitacionID int    `json:"habitacionId"`
	FechaEntrada string `json:"fechaEntrada"` // Formato: YYYY-MM-DD
	FechaSalida  string `json:"fechaSalida"`  // Formato: YYYY-MM-DD
}

// CreateReserva crea una nueva reserva
func (h *ReservaHandler) CreateReserva(c *fiber.Ctx) error {
	var req CreateReservaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de solicitud inválido",
		})
	}

	// Validaciones básicas
	if req.ClienteID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El clienteId es requerido",
		})
	}

	if len(req.Habitaciones) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Debe incluir al menos una habitación",
		})
	}

	// Convertir habitaciones
	habitaciones := make([]domain.ReservaHabitacion, len(req.Habitaciones))
	for i, hab := range req.Habitaciones {
		fechaEntrada, err := time.Parse("2006-01-02", hab.FechaEntrada)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Formato de fechaEntrada inválido. Use YYYY-MM-DD",
			})
		}

		fechaSalida, err := time.Parse("2006-01-02", hab.FechaSalida)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Formato de fechaSalida inválido. Use YYYY-MM-DD",
			})
		}

		habitaciones[i] = domain.ReservaHabitacion{
			HabitacionID: hab.HabitacionID,
			Precio:       hab.Precio,
			FechaEntrada: fechaEntrada,
			FechaSalida:  fechaSalida,
			Estado:       1, // Activa
		}
	}

	// Crear la reserva
	reserva := &domain.Reserva{
		CantidadAdultos:   req.CantidadAdultos,
		CantidadNinhos:    req.CantidadNinhos,
		ClienteID:         req.ClienteID,
		Descuento:         req.Descuento,
		Estado:            domain.ReservaPendiente,
		FechaConfirmacion: time.Now(),
		Habitaciones:      habitaciones,
	}

	if err := h.service.CreateReserva(reserva); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Reserva creada exitosamente",
		"data":    reserva,
	})
}

// GetReservaByID obtiene una reserva por su ID
func (h *ReservaHandler) GetReservaByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID de reserva inválido",
		})
	}

	reserva, err := h.service.GetReservaByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": reserva,
	})
}

// GetReservasCliente obtiene todas las reservas de un cliente
func (h *ReservaHandler) GetReservasCliente(c *fiber.Ctx) error {
	clienteID := c.Params("clienteId")
	if clienteID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "clienteId es requerido",
		})
	}

	reservas, err := h.service.GetReservasCliente(clienteID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": reservas,
	})
}

// UpdateReservaEstado actualiza el estado de una reserva
func (h *ReservaHandler) UpdateReservaEstado(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID de reserva inválido",
		})
	}

	var req UpdateEstadoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de solicitud inválido",
		})
	}

	// Convertir el estado a EstadoReserva
	estado := domain.EstadoReserva(req.Estado)

	if err := h.service.UpdateReservaEstado(id, estado); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Estado de reserva actualizado exitosamente",
	})
}

// CancelarReserva cancela una reserva
func (h *ReservaHandler) CancelarReserva(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID de reserva inválido",
		})
	}

	if err := h.service.CancelarReserva(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reserva cancelada exitosamente",
	})
}

// ConfirmarReserva confirma una reserva pendiente (sin enviar email)
func (h *ReservaHandler) ConfirmarReserva(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID de reserva inválido",
		})
	}

	if err := h.service.ConfirmarReserva(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reserva confirmada exitosamente",
	})
}

// ConfirmarPago confirma el pago de una reserva y envía email automáticamente
func (h *ReservaHandler) ConfirmarPago(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID de reserva inválido",
		})
	}

	if err := h.service.ConfirmarReserva(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Pago confirmado y email enviado exitosamente",
	})
}

// VerificarDisponibilidad verifica si una habitación está disponible
func (h *ReservaHandler) VerificarDisponibilidad(c *fiber.Ctx) error {
	var req VerificarDisponibilidadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de solicitud inválido",
		})
	}

	fechaEntrada, err := time.Parse("2006-01-02", req.FechaEntrada)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de fechaEntrada inválido. Use YYYY-MM-DD",
		})
	}

	fechaSalida, err := time.Parse("2006-01-02", req.FechaSalida)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de fechaSalida inválido. Use YYYY-MM-DD",
		})
	}

	disponible, err := h.service.VerificarDisponibilidad(req.HabitacionID, fechaEntrada, fechaSalida)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"disponible": disponible,
	})
}

// GetReservasEnRango obtiene todas las reservas en un rango de fechas
func (h *ReservaHandler) GetReservasEnRango(c *fiber.Ctx) error {
	fechaInicioStr := c.Query("fechaInicio")
	fechaFinStr := c.Query("fechaFin")

	if fechaInicioStr == "" || fechaFinStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "fechaInicio y fechaFin son requeridos",
		})
	}

	fechaInicio, err := time.Parse("2006-01-02", fechaInicioStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de fechaInicio inválido. Use YYYY-MM-DD",
		})
	}

	fechaFin, err := time.Parse("2006-01-02", fechaFinStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Formato de fechaFin inválido. Use YYYY-MM-DD",
		})
	}

	reservas, err := h.service.GetReservasEnRango(fechaInicio, fechaFin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": reservas,
	})
}
