package logic
// TODO: Refactor this into the reservation_handler file (or move logic from reservation_handler here)

import (
	"context"
	"errors"
	"time"

	"groveapi/internal/store"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrRuleViolation = errors.New("rule_violation")
	ErrConflict      = errors.New("conflict") // overlap
)

// Exported fields so BodyParser can fill them.
// JSON names match the API spec.
type ReservationInput struct {
	ResourceID string    `json:"resource_id"`
	UserID     string    `json:"user_id"` // temp until auth
	StartsAt   time.Time `json:"starts_at"`
	EndsAt     time.Time `json:"ends_at"`
}

type ReservationDTO struct {
	ID         string    `json:"id"`
	ResourceID string    `json:"resource_id"`
	UserID     string    `json:"user_id"`
	StartsAt   time.Time `json:"starts_at"`
	EndsAt     time.Time `json:"ends_at"`
	Status     string    `json:"status"`
}

func CreateReservation(ctx context.Context, db *gorm.DB, in ReservationInput) (ReservationDTO, error) {
	// Basic validation (add buffers/quotas later)
	if !in.EndsAt.After(in.StartsAt) {
		return ReservationDTO{}, ErrRuleViolation
	}

	resID, err := uuid.Parse(in.ResourceID)
	if err != nil {
		return ReservationDTO{}, ErrRuleViolation
	}
	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		return ReservationDTO{}, ErrRuleViolation
	}

	var row store.Reservation
	txErr := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		held := store.Reservation{
			ResourceID: resID,
			UserID:     userID,
			StartsAt:   in.StartsAt,
			EndsAt:     in.EndsAt,
			Status:     string(store.StatusHeld),
		}
		if err := tx.Create(&held).Error; err != nil {
			if isExclusionViolation(err) {
				return ErrConflict
			}
			return err
		}
		// Confirm
		if err := tx.Model(&held).Update("status", string(store.StatusConfirmed)).Error; err != nil {
			return err
		}
		row = held
		return nil
	})
	if txErr != nil {
		return ReservationDTO{}, txErr
	}

	return ReservationDTO{
		ID:         row.ID.String(),
		ResourceID: row.ResourceID.String(),
		UserID:     row.UserID.String(),
		StartsAt:   row.StartsAt,
		EndsAt:     row.EndsAt,
		Status:     row.Status,
	}, nil
}

func isExclusionViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23P01" // exclusion constraint violation
	}
	return false
}
