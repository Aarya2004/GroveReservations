package store

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ---- Enums ----
type ReservationStatus string

const (
	StatusHeld      ReservationStatus = "HELD"
	StatusConfirmed ReservationStatus = "CONFIRMED"
	StatusCancelled ReservationStatus = "CANCELLED"
	StatusNoShow    ReservationStatus = "NOSHOW"
	StatusCompleted ReservationStatus = "COMPLETED"
)

// ---- Models (exported fields + tags) ----

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Email       string    `gorm:"uniqueIndex;not null" json:"email"`
	VillaNumber int       `json:"villa_number"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `gorm:"type:text;not null;default:MEMBER" json:"role"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Active      bool      `gorm:"not null;default:true" json:"active"`
}

type Resource struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name           string         `gorm:"not null" json:"name"`
	Type           string         `gorm:"not null" json:"type"`
	Location       *string        `json:"location"`
	SlotMinutes    int            `gorm:"not null;default:60" json:"slot_minutes"`
	BufferMinutes  int            `gorm:"not null;default:0" json:"buffer_minutes"`
	MaxAdvanceDays int            `gorm:"not null;default:14" json:"max_advance_days"`
	OpenHours      datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'::jsonb" json:"open_hours"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

type Reservation struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ResourceID uuid.UUID `gorm:"type:uuid;not null;index" json:"resource_id"`
	StartsAt   time.Time `gorm:"not null" json:"starts_at"`
	EndsAt     time.Time `gorm:"not null" json:"ends_at"`
	Status     string    `gorm:"type:text;not null" json:"status"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Exclude generated column from ORM
	TimeRange any `gorm:"-:all" json:"-"`
}

type AuditLog struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Action       string     `gorm:"not null" json:"action"`
	ResourceID   *uuid.UUID `gorm:"type:uuid" json:"resource_id"`
	ResourceType *string    `json:"resource_type"`
	Details      string     `json:"details"`
	Timestamp    time.Time  `gorm:"autoCreateTime" json:"timestamp"`
}
