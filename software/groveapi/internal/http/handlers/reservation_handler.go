package handlers

import (
	"errors"
	"groveapi/internal/logic"
	"groveapi/internal/store"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReservationHTTP struct {
	DB *gorm.DB
}

func NewReservationHTTP(db *gorm.DB) *ReservationHTTP { return &ReservationHTTP{DB: db} }

func (h *ReservationHTTP) CreateReservation(c *fiber.Ctx) error {
	var in logic.ReservationInput
	if err := c.BodyParser(&in); err != nil {
		return SendError(c, http.StatusBadRequest, "invalid_json", "could not parse request body")
	}

	userID, _ := c.Locals("user_id").(string)
	if userID != "" {
		in.UserID = userID
	}

	out, err := logic.CreateReservation(c.Context(), h.DB, in)
	switch {
	case err == nil:
		return c.Status(http.StatusCreated).JSON(out)
	case errors.Is(err, logic.ErrConflict):
		return SendError(c, http.StatusConflict, "conflict", "time slot is already booked")
	case errors.Is(err, logic.ErrRuleViolation):
		return SendError(c, http.StatusBadRequest, "rule_violation", "reservation violates booking rules")
	default:
		return SendError(c, http.StatusInternalServerError, "internal", "failed to create reservation")
	}
}

func (h *ReservationHTTP) ListCurrentUserReservations(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	var reservations []store.Reservation
	result := h.DB.Where("user_id = ?", userID).Find(&reservations)
	if result.Error != nil {
		return SendError(c, http.StatusInternalServerError, "internal_error", "failed to list reservations")
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"count":        result.RowsAffected,
		"reservations": reservations,
	})
}

func (h *ReservationHTTP) ListResourceReservations(c *fiber.Ctx) error {
	resourceID := c.Query("resource_id")
	if resourceID == "" {
		return SendError(c, http.StatusBadRequest, "missing_param", "resource_id query parameter is required")
	}
	if _, err := uuid.Parse(resourceID); err != nil {
		return SendError(c, http.StatusBadRequest, "invalid_uuid", "resource_id is not a valid UUID")
	}

	query := h.DB.Where("resource_id = ?", resourceID)

	if from := c.Query("from"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			query = query.Where("ends_at >= ?", t)
		}
	}
	if to := c.Query("to"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			query = query.Where("starts_at <= ?", t)
		}
	}

	var reservations []store.Reservation
	result := query.Find(&reservations)
	if result.Error != nil {
		return SendError(c, http.StatusInternalServerError, "internal_error", "failed to list reservations")
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"count":        result.RowsAffected,
		"reservations": reservations,
	})
}

type UpdateReservationInput struct {
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
	Status   string    `json:"status"`
}

func (h *ReservationHTTP) UpdateReservation(c *fiber.Ctx) error {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		return nil
	}

	var in UpdateReservationInput
	if err := c.BodyParser(&in); err != nil {
		return SendError(c, http.StatusBadRequest, "invalid_json", "could not parse request body")
	}

	var reservation store.Reservation
	if err := h.DB.Where("id = ?", id).First(&reservation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return SendError(c, http.StatusNotFound, "reservation_not_found", "reservation not found")
		}
		return SendError(c, http.StatusInternalServerError, "internal_error", "failed to find reservation")
	}

	updates := map[string]any{}
	if !in.StartsAt.IsZero() {
		updates["starts_at"] = in.StartsAt
	}
	if !in.EndsAt.IsZero() {
		updates["ends_at"] = in.EndsAt
	}
	if in.Status != "" {
		updates["status"] = in.Status
	}

	if len(updates) > 0 {
		result := h.DB.Model(&reservation).Updates(updates)
		if result.Error != nil {
			return SendError(c, http.StatusInternalServerError, "update_failed", "failed to update reservation")
		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"id": reservation.ID})
}

func (h *ReservationHTTP) DeleteReservation(c *fiber.Ctx) error {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		return nil
	}

	var reservation store.Reservation
	if err := h.DB.Where("id = ?", id).First(&reservation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return SendError(c, http.StatusNotFound, "reservation_not_found", "reservation not found")
		}
		return SendError(c, http.StatusInternalServerError, "internal_error", "failed to find reservation")
	}

	if err := h.DB.Delete(&reservation).Error; err != nil {
		return SendError(c, http.StatusInternalServerError, "internal_error", "failed to delete reservation")
	}
	return c.SendStatus(http.StatusNoContent)
}
