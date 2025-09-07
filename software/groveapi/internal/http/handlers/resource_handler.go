package handlers

import (
	"groveapi/internal/logic"
	"groveapi/internal/store"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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

func (h *ResourceHTTP) ListResources(c *fiber.Ctx) error   {
	var resources []store.Resource
	result := h.DB.Find(&resources)
	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal_error", "detail": result.Error.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"count": result.RowsAffected,
		"resources": resources,
	}) 
}

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

func (h *ResourceHTTP) GetResource(c *fiber.Ctx) error { 
	var resource store.Resource;
	resource.ID = uuid.MustParse(c.Params("id"))
	result := h.DB.First(&resource)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "resource_not_found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal_error", "detail": result.Error.Error()})
	}
	return c.Status(http.StatusOK).JSON(resource)
}

type UpdateResourceInput struct {
	Name           string
	Type           string
	Location       *string
	SlotMinutes    int
	BufferMinutes  int
	MaxAdvanceDays int
	OpenHours      datatypes.JSON
}

func (h *ResourceHTTP) UpdateResource(c *fiber.Ctx) error { 
	var resource store.Resource;
	var in UpdateResourceInput;

	if err := c.BodyParser(&in); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_json"})
	}
	resource.ID = uuid.MustParse(c.Params("id"))
	if err := h.DB.First(&resource).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "resource_not_found"})
	}
	h.DB.Model(&resource).Updates(store.Resource{Name: in.Name, Type: in.Type, Location: in.Location, SlotMinutes: in.SlotMinutes,
		BufferMinutes: in.BufferMinutes, MaxAdvanceDays: in.MaxAdvanceDays, OpenHours: in.OpenHours})
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id": resource.ID,
	})
}

func (h *ResourceHTTP) DeleteResource(c *fiber.Ctx) error { 
	var resource store.Resource;
	resource.ID = uuid.MustParse(c.Params("id"))
	if err := h.DB.First(&resource).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "resource_not_found"})
	}
	if err := h.DB.Delete(&resource).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal_error", "detail": err.Error()})
	}
	return c.SendStatus(http.StatusNoContent)
}
