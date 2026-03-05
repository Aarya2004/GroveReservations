package logic

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"groveapi/internal/store"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	ErrBadInput = errors.New("bad_input")
	ErrNotFound = errors.New("not_found")
)

type ResourceInput struct {
	Name           string          `json:"name"`
	Type           string          `json:"type"`
	Location       *string         `json:"location"`
	SlotMinutes    int             `json:"slot_minutes"`
	BufferMinutes  int             `json:"buffer_minutes"`
	MaxAdvanceDays int             `json:"max_advance_days"`
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

func resourceToDTO(r store.Resource) ResourceDTO {
	return ResourceDTO{
		ID:             r.ID.String(),
		Name:           r.Name,
		Type:           r.Type,
		Location:       r.Location,
		SlotMinutes:    r.SlotMinutes,
		BufferMinutes:  r.BufferMinutes,
		MaxAdvanceDays: r.MaxAdvanceDays,
		OpenHours:      json.RawMessage(r.OpenHours),
	}
}

func applyResourceDefaults(in *ResourceInput) {
	if in.SlotMinutes <= 0 {
		in.SlotMinutes = 60
	}
	if in.BufferMinutes < 0 {
		in.BufferMinutes = 0
	}
	if in.MaxAdvanceDays <= 0 {
		in.MaxAdvanceDays = 14
	}
	if len(in.OpenHours) == 0 || string(in.OpenHours) == "null" {
		in.OpenHours = json.RawMessage(`{}`)
	}
}

func validateResourceInput(in ResourceInput) error {
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Type) == "" {
		return ErrBadInput
	}
	if len(in.OpenHours) > 0 && string(in.OpenHours) != "null" && string(in.OpenHours) != "{}" {
		var tmp any
		if err := json.Unmarshal(in.OpenHours, &tmp); err != nil {
			return ErrBadInput
		}
	}
	return nil
}

func CreateResource(ctx context.Context, db *gorm.DB, in ResourceInput) (ResourceDTO, error) {
	applyResourceDefaults(&in)
	if err := validateResourceInput(in); err != nil {
		return ResourceDTO{}, err
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
	return resourceToDTO(row), nil
}

func ListResources(ctx context.Context, db *gorm.DB) ([]ResourceDTO, error) {
	var resources []store.Resource
	if err := db.WithContext(ctx).Find(&resources).Error; err != nil {
		return nil, err
	}
	dtos := make([]ResourceDTO, len(resources))
	for i, r := range resources {
		dtos[i] = resourceToDTO(r)
	}
	return dtos, nil
}

func GetResource(ctx context.Context, db *gorm.DB, id uuid.UUID) (ResourceDTO, error) {
	var resource store.Resource
	if err := db.WithContext(ctx).Where("id = ?", id).First(&resource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ResourceDTO{}, ErrNotFound
		}
		return ResourceDTO{}, err
	}
	return resourceToDTO(resource), nil
}

func UpdateResource(ctx context.Context, db *gorm.DB, id uuid.UUID, in ResourceInput) (ResourceDTO, error) {
	var resource store.Resource
	if err := db.WithContext(ctx).Where("id = ?", id).First(&resource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ResourceDTO{}, ErrNotFound
		}
		return ResourceDTO{}, err
	}

	updates := store.Resource{
		Name:           in.Name,
		Type:           in.Type,
		Location:       in.Location,
		SlotMinutes:    in.SlotMinutes,
		BufferMinutes:  in.BufferMinutes,
		MaxAdvanceDays: in.MaxAdvanceDays,
	}
	if len(in.OpenHours) > 0 && string(in.OpenHours) != "null" {
		var tmp any
		if err := json.Unmarshal(in.OpenHours, &tmp); err != nil {
			return ResourceDTO{}, ErrBadInput
		}
		updates.OpenHours = datatypes.JSON(in.OpenHours)
	}

	if err := db.WithContext(ctx).Model(&resource).Updates(updates).Error; err != nil {
		return ResourceDTO{}, err
	}

	if err := db.WithContext(ctx).Where("id = ?", id).First(&resource).Error; err != nil {
		return ResourceDTO{}, err
	}
	return resourceToDTO(resource), nil
}

func DeleteResource(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	var resource store.Resource
	if err := db.WithContext(ctx).Where("id = ?", id).First(&resource).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return db.WithContext(ctx).Delete(&resource).Error
}
