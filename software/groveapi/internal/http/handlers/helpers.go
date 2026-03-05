package handlers

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var errResponseSent = errors.New("response already sent")

func parseUUIDParam(c *fiber.Ctx, name string) (uuid.UUID, error) {
	raw := c.Params(name)
	id, err := uuid.Parse(raw)
	if err != nil {
		_ = c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid_uuid",
			"message": "parameter '" + name + "' is not a valid UUID",
		})
		return uuid.Nil, errResponseSent
	}
	return id, nil
}
