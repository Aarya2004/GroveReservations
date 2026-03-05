package logic_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"groveapi/internal/logic"
	"groveapi/internal/store"
	"groveapi/internal/testutil"

	"github.com/google/uuid"
)

func createTestResource(t *testing.T, db interface{ WithContext(context.Context) *store.Resource }) {
	// helper handled inline
}

func TestCreateReservation_Valid(t *testing.T) {
	db := testutil.MustOpenTestDB(t)
	tx := testutil.WrapInTx(t, db)

	// Create a resource first
	resource := store.Resource{
		Name:           "Test Court",
		Type:           "tennis_court",
		SlotMinutes:    60,
		BufferMinutes:  0,
		MaxAdvanceDays: 14,
		OpenHours:      []byte(`{}`),
	}
	if err := tx.Create(&resource).Error; err != nil {
		t.Fatalf("failed to create test resource: %v", err)
	}

	// Create a user
	user := store.User{
		Name:   "Test User",
		Email:  "test-" + uuid.New().String()[:8] + "@example.com",
		Role:   "MEMBER",
		Active: true,
	}
	if err := tx.Create(&user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	in := logic.ReservationInput{
		ResourceID: resource.ID.String(),
		UserID:     user.ID.String(),
		StartsAt:   time.Now().Add(1 * time.Hour).Truncate(time.Second),
		EndsAt:     time.Now().Add(2 * time.Hour).Truncate(time.Second),
	}

	out, err := logic.CreateReservation(context.Background(), tx, in)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if out.Status != "CONFIRMED" {
		t.Errorf("expected status CONFIRMED, got %s", out.Status)
	}
	if out.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestCreateReservation_EndsBeforeStarts(t *testing.T) {
	db := testutil.MustOpenTestDB(t)
	tx := testutil.WrapInTx(t, db)

	in := logic.ReservationInput{
		ResourceID: uuid.New().String(),
		UserID:     uuid.New().String(),
		StartsAt:   time.Now().Add(2 * time.Hour),
		EndsAt:     time.Now().Add(1 * time.Hour),
	}

	_, err := logic.CreateReservation(context.Background(), tx, in)
	if !errors.Is(err, logic.ErrRuleViolation) {
		t.Errorf("expected ErrRuleViolation, got: %v", err)
	}
}

func TestCreateReservation_InvalidUUIDs(t *testing.T) {
	db := testutil.MustOpenTestDB(t)
	tx := testutil.WrapInTx(t, db)

	tests := []struct {
		name       string
		resourceID string
		userID     string
	}{
		{"invalid resource_id", "not-a-uuid", uuid.New().String()},
		{"invalid user_id", uuid.New().String(), "not-a-uuid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := logic.ReservationInput{
				ResourceID: tt.resourceID,
				UserID:     tt.userID,
				StartsAt:   time.Now().Add(1 * time.Hour),
				EndsAt:     time.Now().Add(2 * time.Hour),
			}
			_, err := logic.CreateReservation(context.Background(), tx, in)
			if !errors.Is(err, logic.ErrRuleViolation) {
				t.Errorf("expected ErrRuleViolation, got: %v", err)
			}
		})
	}
}

func TestCreateReservation_Overlap(t *testing.T) {
	db := testutil.MustOpenTestDB(t)
	tx := testutil.WrapInTx(t, db)

	resource := store.Resource{
		Name:        "Overlap Court",
		Type:        "tennis_court",
		SlotMinutes: 60,
		OpenHours:   []byte(`{}`),
	}
	if err := tx.Create(&resource).Error; err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}

	user := store.User{
		Name:   "Overlap User",
		Email:  "overlap-" + uuid.New().String()[:8] + "@example.com",
		Role:   "MEMBER",
		Active: true,
	}
	if err := tx.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	startsAt := time.Now().Add(3 * time.Hour).Truncate(time.Second)
	endsAt := startsAt.Add(1 * time.Hour)

	in := logic.ReservationInput{
		ResourceID: resource.ID.String(),
		UserID:     user.ID.String(),
		StartsAt:   startsAt,
		EndsAt:     endsAt,
	}

	// First reservation should succeed
	_, err := logic.CreateReservation(context.Background(), tx, in)
	if err != nil {
		t.Fatalf("first reservation failed: %v", err)
	}

	// Second overlapping reservation should conflict
	_, err = logic.CreateReservation(context.Background(), tx, in)
	if !errors.Is(err, logic.ErrConflict) {
		t.Errorf("expected ErrConflict, got: %v", err)
	}
}
