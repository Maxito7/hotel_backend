package http

import (
	"strconv"

	"github.com/Maxito7/hotel_backend/internal/application"
	"github.com/Maxito7/hotel_backend/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type ContactHandler struct {
	service *application.ContactService
}

func NewContactHandler(s *application.ContactService) *ContactHandler {
	return &ContactHandler{service: s}
}

func (h *ContactHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/api/contacto")
	group.Post("/", h.Create)
	group.Get("/", h.List)
	group.Patch("/:id/estado", h.UpdateEstado)
}

func (h *ContactHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateContactRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}
	id, err := h.service.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id, "estado": "Nuevo"})
}

func (h *ContactHandler) List(c *fiber.Ctx) error {
	items, err := h.service.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func (h *ContactHandler) UpdateEstado(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	var req domain.UpdateEstadoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}
	if err := h.service.UpdateEstado(c.Context(), id, req.Estado); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
