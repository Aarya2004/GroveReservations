package store

import (
	"time"

	"github.com/google/uuid"
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

type ResourceType string

const (
	ResourceTennisCourt ResourceType = "TENNIS_COURT"
)

type Role string

const (
	RoleAdmin  Role = "ADMIN"
	RoleMember Role = "MEMBER"
	RoleGuest  Role = "GUEST"
)

// ---- Models (exported fields + tags) ----

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"not null"`
	Email       string    `gorm:"uniqueIndex;not null"`
	VillaNumber int
	PhoneNumber string
	Role        Role      `gorm:"type:text;not null;default:MEMBER"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type Resource struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name          string      `gorm:"not null"`
	Location      string
	ResourceType  ResourceType `gorm:"type:text;not null"`
	OpeningHours  string       // expand later to JSON if needed
	ClosingHours  string
	MaxSlotLength int          `gorm:"not null;default:60"`
	CreatedAt     time.Time    `gorm:"autoCreateTime"`
	UpdatedAt     time.Time    `gorm:"autoUpdateTime"`
}

type Reservation struct {
	ID         uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID         `gorm:"type:uuid;not null;index"`
	ResourceID uuid.UUID         `gorm:"type:uuid;not null;index"`
	StartsAt   time.Time         `gorm:"not null"`
	EndsAt     time.Time         `gorm:"not null"`
	Status     string            `gorm:"type:text;not null"` // keep as string for simplicity
	CreatedAt  time.Time         `gorm:"autoCreateTime"`
	UpdatedAt  time.Time         `gorm:"autoUpdateTime"`

	// Exclude generated column from ORM
	TimeRange any `gorm:"-:all"`
}

type AuditLog struct {
	ID           uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID   `gorm:"type:uuid;not null;index"`
	Action       string      `gorm:"not null"`
	ResourceID   *uuid.UUID  `gorm:"type:uuid"`
	ResourceType *string
	Details      string
	Timestamp    time.Time   `gorm:"autoCreateTime"`
}
