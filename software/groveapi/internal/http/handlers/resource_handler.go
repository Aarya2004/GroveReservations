package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"groveapi/internal/logic"
	"net/http"
)

// GET `/resources` – List all resources.
// POST `/resources` *(admin)* – Create resource `{ name, type, location, rules }`.
// GET `/resources/:id` – Get resource details and rules.
// PATCH `/resources/:id` *(admin)* – Update resource/rules.
// DELETE `/resources/:id` *(admin)* – Remove resource.

type ResourceHTTP struct {
	DB *gorm.DB
}

func NewResourceHTTP(db *gorm.DB) *ResourceHTTP { return &ResourceHTTP{DB: db} }

func (h *ResourceHTTP) ListResources(c *fiber.Ctx) error   { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }

func (h *ResourceHTTP) CreateResource(c *fiber.Ctx) error {
	var in logic.ResourceInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_json"})
	}
	out, err := logic.CreateResource(c.Context(), h.DB, in)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusCreated).JSON(out)
}

func (h *ResourceHTTP) GetResource(c *fiber.Ctx) error     { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }
func (h *ResourceHTTP) UpdateResource(c *fiber.Ctx) error     { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }
func (h *ResourceHTTP) DeleteResource(c *fiber.Ctx) error     { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }