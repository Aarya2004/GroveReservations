package handlers

import (
	"errors"
	"groveapi/internal/logic"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AvailabilityHTTP struct {
	DB *gorm.DB
}

func NewAvailabilityHTTP(db *gorm.DB) *AvailabilityHTTP { return &AvailabilityHTTP{DB: db} }

func (h *AvailabilityHTTP) GetAvailability(c *fiber.Ctx) error {
	resourceID := c.Query("resource_id")
	if resourceID == "" {
		return SendError(c, http.StatusBadRequest, "missing_param", "resource_id is required")
	}
	rid, err := uuid.Parse(resourceID)
	if err != nil {
		return SendError(c, http.StatusBadRequest, "invalid_uuid", "resource_id is not a valid UUID")
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr == "" || toStr == "" {
		return SendError(c, http.StatusBadRequest, "missing_param", "from and to are required")
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return SendError(c, http.StatusBadRequest, "invalid_date", "from must be RFC3339 format")
	}
	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return SendError(c, http.StatusBadRequest, "invalid_date", "to must be RFC3339 format")
	}

	result, err := logic.GetAvailability(c.Context(), h.DB, rid, from, to)
	if err != nil {
		if errors.Is(err, logic.ErrNotFound) {
			return SendError(c, http.StatusNotFound, "resource_not_found", "resource not found")
		}
		return SendError(c, http.StatusInternalServerError, "internal_error", "failed to get availability")
	}
	return c.Status(http.StatusOK).JSON(result)
}
