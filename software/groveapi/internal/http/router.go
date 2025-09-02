package http

import (
	"errors"
	"net/http"

	"groveapi/internal/logic"
	"groveapi/internal/store"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, s *store.Store) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/reservations", func(c *fiber.Ctx) error {
		var in logic.ReservationInput
		if err := c.BodyParser(&in); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"invalid_json"})
		}
		res, err := logic.CreateReservation(c.Context(), s.DB, in)
		switch {
		case err == nil:
			return c.Status(http.StatusCreated).JSON(res)
		case errors.Is(err, logic.ErrConflict):
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error":"conflict"})
		case errors.Is(err, logic.ErrRuleViolation):
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"rule_violation"})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"internal"})
		}
	})
}
