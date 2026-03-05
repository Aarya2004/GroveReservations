package handlers

import (
	"errors"
	"groveapi/internal/logic"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ResourceHTTP struct {
	DB *gorm.DB
}

func NewResourceHTTP(db *gorm.DB) *ResourceHTTP { return &ResourceHTTP{DB: db} }

func (h *ResourceHTTP) ListResources(c *fiber.Ctx) error {
	resources, err := logic.ListResources(c.Context(), h.DB)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, "internal_error", "failed to list resources")
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"count":     len(resources),
		"resources": resources,
	})
}

func (h *ResourceHTTP) CreateResource(c *fiber.Ctx) error {
	var in logic.ResourceInput
	if err := c.BodyParser(&in); err != nil {
		return SendError(c, http.StatusBadRequest, "invalid_json", "could not parse request body")
	}
	out, err := logic.CreateResource(c.Context(), h.DB, in)
	if err != nil {
		if errors.Is(err, logic.ErrBadInput) {
			return SendError(c, http.StatusBadRequest, "bad_input", "invalid resource data")
		}
		return SendError(c, http.StatusInternalServerError, "internal_error", "failed to create resource")
	}
	return c.Status(http.StatusCreated).JSON(out)
}

func (h *ResourceHTTP) GetResource(c *fiber.Ctx) error {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		return nil
	}
	out, err := logic.GetResource(c.Context(), h.DB, id)
	if err != nil {
		if errors.Is(err, logic.ErrNotFound) {
			return SendError(c, http.StatusNotFound, "resource_not_found", "resource not found")
		}
		return SendError(c, http.StatusInternalServerError, "internal_error", "failed to get resource")
	}
	return c.Status(http.StatusOK).JSON(out)
}

func (h *ResourceHTTP) UpdateResource(c *fiber.Ctx) error {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		return nil
	}
	var in logic.ResourceInput
	if err := c.BodyParser(&in); err != nil {
		return SendError(c, http.StatusBadRequest, "invalid_json", "could not parse request body")
	}
	out, err := logic.UpdateResource(c.Context(), h.DB, id, in)
	if err != nil {
		if errors.Is(err, logic.ErrNotFound) {
			return SendError(c, http.StatusNotFound, "resource_not_found", "resource not found")
		}
		if errors.Is(err, logic.ErrBadInput) {
			return SendError(c, http.StatusBadRequest, "bad_input", "invalid resource data")
		}
		return SendError(c, http.StatusInternalServerError, "internal_error", "failed to update resource")
	}
	return c.Status(http.StatusOK).JSON(out)
}

func (h *ResourceHTTP) DeleteResource(c *fiber.Ctx) error {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		return nil
	}
	if err := logic.DeleteResource(c.Context(), h.DB, id); err != nil {
		if errors.Is(err, logic.ErrNotFound) {
			return SendError(c, http.StatusNotFound, "resource_not_found", "resource not found")
		}
		return SendError(c, http.StatusInternalServerError, "internal_error", "failed to delete resource")
	}
	return c.SendStatus(http.StatusNoContent)
}
