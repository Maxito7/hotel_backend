package http

import (
	"github.com/Maxito7/hotel_backend/internal/application"
	"github.com/gofiber/fiber/v2"
)

type ServicioHandler struct {
	service *application.ServicioService
}

func NewServicioHandler(service *application.ServicioService) *ServicioHandler {
	return &ServicioHandler{
		service: service,
	}
}

func (h *ServicioHandler) GetAllServices(c *fiber.Ctx) error {
	servicios, err := h.service.GetAllServices()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener los servicios",
		})
	}

	return c.JSON(servicios)
}
