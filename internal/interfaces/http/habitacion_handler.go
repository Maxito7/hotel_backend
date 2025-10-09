package http

import (
	"fmt"
	"log"
	"time"

	"github.com/Maxito7/hotel_backend/internal/application"
	"github.com/gofiber/fiber/v2"
)

type HabitacionHandler struct {
	service *application.HabitacionService
}

func NewHabitacionHandler(service *application.HabitacionService) *HabitacionHandler {
	return &HabitacionHandler{
		service: service,
	}
}

type availableRoomsRequest struct {
	FechaEntrada string `json:"fechaEntrada"`
	FechaSalida  string `json:"fechaSalida"`
}

func (h *HabitacionHandler) GetAllRooms(c *fiber.Ctx) error {
	// Get all rooms
	habitaciones, err := h.service.GetAllRooms()
	if err != nil {
		log.Printf("Error getting rooms: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Error al obtener las habitaciones: %v", err),
		})
	}

	return c.JSON(habitaciones)
}

func (h *HabitacionHandler) GetFechasBloqueadas(c *fiber.Ctx) error {
	// Obtener parámetros de consulta
	desdeStr := c.Query("desde")
	hastaStr := c.Query("hasta", "") // Si no se proporciona, usaremos 3 meses por defecto

	// Parsear fecha desde
	desde, err := time.Parse("2006-01-02", desdeStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid desde format. Use YYYY-MM-DD",
		})
	}

	// Si no se proporciona fecha hasta, usar 3 meses después de desde
	var hasta time.Time
	if hastaStr == "" {
		hasta = desde.AddDate(0, 3, 0) // 3 meses después
	} else {
		hasta, err = time.Parse("2006-01-02", hastaStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid hasta format. Use YYYY-MM-DD",
			})
		}
	}

	// Validar que desde sea menor o igual que hasta
	if desde.After(hasta) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "desde must be before or equal to hasta",
		})
	}

	// Validar que desde no sea anterior a hoy
	if desde.Before(time.Now().Truncate(24 * time.Hour)) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "desde cannot be before today",
		})
	}

	// Obtener fechas bloqueadas
	fechasBloqueadas, err := h.service.GetFechasBloqueadas(desde, hasta)
	if err != nil {
		log.Printf("Error en GetFechasBloqueadas: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Error al obtener las fechas bloqueadas: %v", err),
		})
	}

	return c.JSON(fechasBloqueadas)
}

func (h *HabitacionHandler) GetAvailableRooms(c *fiber.Ctx) error {
	// Parse query parameters
	fechaEntradaStr := c.Query("fechaEntrada")
	fechaSalidaStr := c.Query("fechaSalida")

	if fechaEntradaStr == "" || fechaSalidaStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "fechaEntrada and fechaSalida are required",
		})
	}

	// Parse dates
	fechaEntrada, err := time.Parse("2006-01-02", fechaEntradaStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid fechaEntrada format. Use YYYY-MM-DD",
		})
	}

	fechaSalida, err := time.Parse("2006-01-02", fechaSalidaStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid fechaSalida format. Use YYYY-MM-DD",
		})
	}

	// Validate dates
	if fechaEntrada.After(fechaSalida) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "fechaEntrada must be before fechaSalida",
		})
	}

	// Get available rooms
	habitaciones, err := h.service.GetAvailableRooms(fechaEntrada, fechaSalida)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener las habitaciones disponibles",
		})
	}

	return c.JSON(habitaciones)
}
