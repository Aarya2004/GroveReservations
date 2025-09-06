package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"groveapi/internal/logic"
	"net/http"
)

// GET `/reservations/my` – List my reservations.
// GET `/reservations?resource_id=&from=&to=` *(admin)* – List reservations for a resource.
// POST `/reservations` – Create reservation `{ resource_id, starts_at, ends_at }`.
//   * Response: reservation object OR `409 Conflict`.
// PATCH `/reservations/:id` – Modify reservation (time change, cancel).
// DELETE `/reservations/:id` – Cancel reservation.

type ReservationHTTP struct {
	DB *gorm.DB
}

func NewReservationHTTP(db *gorm.DB) *ReservationHTTP { return &ReservationHTTP{DB: db} }

func (h *ReservationHTTP) CreateReservation(c *fiber.Ctx) error {
	var in logic.ReservationInput
	if err := c.BodyParser(&in); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"invalid_json"})
	}
	// TODO: inject current user from auth middleware; for now relies on JSON user_id
	out, err := logic.CreateReservation(c.Context(), h.DB, in)
	switch {
	case err == nil:
		return c.Status(http.StatusCreated).JSON(out)
	case errors.Is(err, logic.ErrConflict):
		return c.Status(http.StatusConflict).JSON(fiber.Map{"error":"conflict"})
	case errors.Is(err, logic.ErrRuleViolation):
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"rule_violation"})
	default:
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"internal"})
	}
}

func (h *ReservationHTTP) ListCurrentUserReservations(c *fiber.Ctx) error   { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }
func (h *ReservationHTTP) ListResourceReservations(c *fiber.Ctx) error   { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }
func (h *ReservationHTTP) UpdateReservation(c *fiber.Ctx) error     { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }
func (h *ReservationHTTP) DeleteReservation(c *fiber.Ctx) error     { /* TODO */ return c.SendStatus(http.StatusNotImplemented) }
