package logic_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"groveapi/internal/logic"
	"groveapi/internal/testutil"
)

func TestCreateResource_Valid(t *testing.T) {
	db := testutil.MustOpenTestDB(t)
	tx := testutil.WrapInTx(t, db)

	in := logic.ResourceInput{
		Name: "Court A",
		Type: "tennis_court",
	}

	out, err := logic.CreateResource(context.Background(), tx, in)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if out.Name != "Court A" {
		t.Errorf("expected name 'Court A', got '%s'", out.Name)
	}
	if out.SlotMinutes != 60 {
		t.Errorf("expected default slot_minutes 60, got %d", out.SlotMinutes)
	}
	if out.MaxAdvanceDays != 14 {
		t.Errorf("expected default max_advance_days 14, got %d", out.MaxAdvanceDays)
	}
	if out.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestCreateResource_EmptyName(t *testing.T) {
	db := testutil.MustOpenTestDB(t)
	tx := testutil.WrapInTx(t, db)

	in := logic.ResourceInput{
		Name: "",
		Type: "tennis_court",
	}

	_, err := logic.CreateResource(context.Background(), tx, in)
	if !errors.Is(err, logic.ErrBadInput) {
		t.Errorf("expected ErrBadInput, got: %v", err)
	}
}

func TestCreateResource_EmptyType(t *testing.T) {
	db := testutil.MustOpenTestDB(t)
	tx := testutil.WrapInTx(t, db)

	in := logic.ResourceInput{
		Name: "Court B",
		Type: "",
	}

	_, err := logic.CreateResource(context.Background(), tx, in)
	if !errors.Is(err, logic.ErrBadInput) {
		t.Errorf("expected ErrBadInput, got: %v", err)
	}
}

func TestCreateResource_InvalidOpenHours(t *testing.T) {
	db := testutil.MustOpenTestDB(t)
	tx := testutil.WrapInTx(t, db)

	in := logic.ResourceInput{
		Name:      "Court C",
		Type:      "tennis_court",
		OpenHours: json.RawMessage(`{invalid json`),
	}

	_, err := logic.CreateResource(context.Background(), tx, in)
	if !errors.Is(err, logic.ErrBadInput) {
		t.Errorf("expected ErrBadInput, got: %v", err)
	}
}

func TestListResources(t *testing.T) {
	db := testutil.MustOpenTestDB(t)
	tx := testutil.WrapInTx(t, db)

	// Create two resources
	for _, name := range []string{"Court X", "Court Y"} {
		_, err := logic.CreateResource(context.Background(), tx, logic.ResourceInput{
			Name: name,
			Type: "tennis_court",
		})
		if err != nil {
			t.Fatalf("failed to create resource %s: %v", name, err)
		}
	}

	resources, err := logic.ListResources(context.Background(), tx)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(resources) < 2 {
		t.Errorf("expected at least 2 resources, got %d", len(resources))
	}
}
