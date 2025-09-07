package logic

// TODO: Refactor this into the resource_handler file (or move logic from resource_handler here)

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"groveapi/internal/store"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	ErrBadInput = errors.New("bad_input")
)

type ResourceInput struct {
	Name           string          `json:"name"`
	Type           string          `json:"type"`
	Location       *string         `json:"location"`
	SlotMinutes    int             `json:"slot_minutes"`     // default 60 if 0
	BufferMinutes  int             `json:"buffer_minutes"`   // default 0 if <0
	MaxAdvanceDays int             `json:"max_advance_days"` // default 14 if 0
	OpenHours      json.RawMessage `json:"open_hours"`
}

type ResourceDTO struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Type           string          `json:"type"`
	Location       *string         `json:"location,omitempty"`
	SlotMinutes    int             `json:"slot_minutes"`
	BufferMinutes  int             `json:"buffer_minutes"`
	MaxAdvanceDays int             `json:"max_advance_days"`
	OpenHours      json.RawMessage `json:"open_hours"`
}

func CreateResource(ctx context.Context, db *gorm.DB, in ResourceInput) (ResourceDTO, error) {
	// Basic validation
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Type) == "" {
		return ResourceDTO{}, ErrBadInput
	}
	if in.SlotMinutes <= 0 {
		in.SlotMinutes = 60
	}
	if in.BufferMinutes < 0 {
		in.BufferMinutes = 0
	}
	if in.MaxAdvanceDays <= 0 {
		in.MaxAdvanceDays = 14
	}
	// Accept open_hours as any valid JSON object/array; if empty, use {}
	if len(in.OpenHours) == 0 || string(in.OpenHours) == "null" {
		in.OpenHours = json.RawMessage(`{}`)
	} else {
		var tmp any
		if err := json.Unmarshal(in.OpenHours, &tmp); err != nil {
			return ResourceDTO{}, ErrBadInput
		}
	}

	row := store.Resource{
		Name:           in.Name,
		Type:           strings.ToLower(in.Type),
		Location:       in.Location,
		SlotMinutes:    in.SlotMinutes,
		BufferMinutes:  in.BufferMinutes,
		MaxAdvanceDays: in.MaxAdvanceDays,
		OpenHours:      datatypes.JSON(in.OpenHours),
	}

	if err := db.WithContext(ctx).Create(&row).Error; err != nil {
		return ResourceDTO{}, err
	}

	return ResourceDTO{
		ID:             row.ID.String(),
		Name:           row.Name,
		Type:           row.Type,
		Location:       row.Location,
		SlotMinutes:    row.SlotMinutes,
		BufferMinutes:  row.BufferMinutes,
		MaxAdvanceDays: row.MaxAdvanceDays,
		OpenHours:      json.RawMessage(row.OpenHours), // cast back for JSON response
	}, nil
}
