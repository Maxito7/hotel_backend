package http

import (
	"fmt"
	"log"
	"time"

	"github.com/Maxito7/hotel_backend/internal/application"
	"github.com/gofiber/fiber/v2"
)

// Zona horaria de Perú (UTC-5)
var peruLocation *time.Location

func init() {
	var err error
	peruLocation, err = time.LoadLocation("America/Lima")
	if err != nil {
		// Fallback a UTC-5 si no se puede cargar la zona horaria
		peruLocation = time.FixedZone("PET", -5*60*60)
	}
}

// parseDatePeru parsea una fecha en formato YYYY-MM-DD y la retorna en zona horaria de Perú
func parseDatePeru(dateStr string) (time.Time, error) {
	// Parsear en UTC primero y luego convertir a zona horaria de Perú
	utcTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	// Crear la fecha en zona horaria de Perú a las 00:00:00
	return time.Date(utcTime.Year(), utcTime.Month(), utcTime.Day(), 0, 0, 0, 0, peruLocation), nil
}

// getTodayPeru retorna la fecha de hoy a las 00:00:00 en zona horaria de Perú
func getTodayPeru() time.Time {
	now := time.Now().In(peruLocation)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, peruLocation)
}

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
	desde, err := parseDatePeru(desdeStr)
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
		hasta, err = parseDatePeru(hastaStr)
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
	hoy := getTodayPeru()
	if desde.Before(hoy) {
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
	fechaEntrada, err := parseDatePeru(fechaEntradaStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid fechaEntrada format. Use YYYY-MM-DD",
		})
	}

	fechaSalida, err := parseDatePeru(fechaSalidaStr)
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

	// Validar que fechaEntrada no sea anterior a hoy
	hoy := getTodayPeru()
	if fechaEntrada.Before(hoy) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "fechaEntrada cannot be before today",
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

// GetRoomTypes returns all available room types
func (h *HabitacionHandler) GetRoomTypes(c *fiber.Ctx) error {
	tipos, err := h.service.GetRoomTypes()
	if err != nil {
		log.Printf("Error getting room types: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Error al obtener los tipos de habitación: %v", err),
		})
	}
	return c.JSON(tipos)
}
