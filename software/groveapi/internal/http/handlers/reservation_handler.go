package handlers

import (
	"errors"
	"groveapi/internal/logic"
	"groveapi/internal/store"
	"net/http"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	supabase "github.com/supabase-community/supabase-go"
	"gorm.io/gorm"
)

// GET `/reservations/me` – List my reservations.
// GET `/reservations?resource_id=&from=&to=` *(admin)* – List reservations for a resource.
// POST `/reservations` – Create reservation `{ resource_id, starts_at, ends_at }`.
//   * Response: reservation object OR `409 Conflict`.
// PATCH `/reservations/:id` – Modify reservation (time change, cancel).
// DELETE `/reservations/:id` – Cancel reservation.

type ReservationHTTP struct {
	DB *gorm.DB
	SB *supabase.Client
}

func NewReservationHTTP(db *gorm.DB, sb *supabase.Client) *ReservationHTTP { return &ReservationHTTP{DB: db, SB: sb} }

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

func (h *ReservationHTTP) ListCurrentUserReservations(c *fiber.Ctx) error { 
	u, err := h.SB.Auth.GetUser()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error_fetching_user", "detail": err.Error()})
	}
	var reservations []store.Reservation
	result := h.DB.Where("user_id = ?", u.ID).Find(&reservations)
	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal_error", "detail": result.Error.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"count": result.RowsAffected,
		"reservations": reservations,
	})
}

func (h *ReservationHTTP) ListResourceReservations(c *fiber.Ctx) error {
	var reservations []store.Reservation
	resourceId := c.Params("id")
	result := h.DB.Where("resource_id = ?", resourceId).Find(&reservations)
	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal_error", "detail": result.Error.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"count": result.RowsAffected,
		"reservations": reservations,
	})
}

type UpdateReservationInput struct {
	UserId        string
	ResourceId    string
	StartsAt      time.Time
	EndsAt        time.Time
	Status        string    
}

func (h *ReservationHTTP) UpdateReservation(c *fiber.Ctx) error { 
	var reservation store.Reservation;
	var in UpdateReservationInput;
	if err := c.BodyParser(&in); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_json"})
	}
	reservation.ID = uuid.MustParse(c.Params("id"))
	result := h.DB.First(&reservation)
	if result.Error != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "error_finding_reservation", "detail": result.Error.Error()})
	}
	userId, err := uuid.Parse(in.UserId)
	if err != nil {
		userId = reservation.UserID
	}
	resourceId, err := uuid.Parse(in.ResourceId)
	if err != nil {
		resourceId = reservation.ResourceID
	}

	h.DB.Model(&reservation).Updates(store.Reservation{UserID: userId, ResourceID: resourceId, StartsAt: in.StartsAt,
		EndsAt: in.EndsAt, Status: in.Status})
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id": reservation.ID,
	})
}

func (h *ReservationHTTP) DeleteReservation(c *fiber.Ctx) error {
	var reservation store.Reservation;
	reservation.ID = uuid.MustParse(c.Params("id"))
	if err := h.DB.First(&reservation).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "reservation_not_found"})
	}
	if err := h.DB.Delete(&reservation).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal_error", "detail": err.Error()})
	}
	return c.SendStatus(http.StatusNoContent) 
}
