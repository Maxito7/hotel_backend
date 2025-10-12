package http

import (
	"encoding/json"

	"github.com/Maxito7/hotel_backend/internal/application"
	"github.com/gofiber/fiber/v2"
)

type SearchHandler struct {
	searchService *application.SearchService
}

func NewSearchHandler(searchService *application.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

type SearchRequest struct {
	Query      string `json:"query"`
	MaxResults int    `json:"max_results"`
	Depth      string `json:"depth"`
}

func (h *SearchHandler) Search(c *fiber.Ctx) error {
	// Parseo manual como backup
	var req SearchRequest

	// Intenta con BodyParser primero
	if err := c.BodyParser(&req); err != nil {
		// Si falla, parsea manualmente
		if err := json.Unmarshal(c.Body(), &req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid JSON format",
				"details": err.Error(),
			})
		}
	}

	if req.Query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query is required",
		})
	}

	input := application.SearchInput{
		Query:      req.Query,
		MaxResults: req.MaxResults,
		Depth:      req.Depth,
	}

	results, err := h.searchService.SearchWeb(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(results)
}
