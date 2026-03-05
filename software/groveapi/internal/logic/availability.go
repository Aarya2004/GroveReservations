package logic

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"groveapi/internal/store"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Slot struct {
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
	Available bool      `json:"available"`
}

type AvailabilityResult struct {
	ResourceID string `json:"resource_id"`
	From       string `json:"from"`
	To         string `json:"to"`
	Slots      []Slot `json:"slots"`
}

type DayHours struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}

func GetAvailability(ctx context.Context, db *gorm.DB, resourceID uuid.UUID, from, to time.Time) (AvailabilityResult, error) {
	var resource store.Resource
	if err := db.WithContext(ctx).Where("id = ?", resourceID).First(&resource).Error; err != nil {
		return AvailabilityResult{}, ErrNotFound
	}

	// Parse open hours
	openHours := map[string]DayHours{}
	if len(resource.OpenHours) > 0 && string(resource.OpenHours) != "{}" {
		_ = json.Unmarshal(resource.OpenHours, &openHours)
	}

	// Get existing reservations in range
	var reservations []store.Reservation
	db.WithContext(ctx).
		Where("resource_id = ? AND ends_at > ? AND starts_at < ? AND status IN ?",
			resourceID, from, to, []string{"HELD", "CONFIRMED"}).
		Find(&reservations)

	slotDuration := time.Duration(resource.SlotMinutes) * time.Minute
	bufferDuration := time.Duration(resource.BufferMinutes) * time.Minute

	var slots []Slot

	// Generate slots day by day
	for day := from; day.Before(to); day = day.AddDate(0, 0, 1) {
		dayName := dayNameLower(day.Weekday())
		hours, ok := openHours[dayName]
		if !ok {
			// Default: 06:00-22:00 if no open hours defined
			hours = DayHours{Open: "06:00", Close: "22:00"}
		}

		openTime, err := time.Parse("15:04", hours.Open)
		if err != nil {
			continue
		}
		closeTime, err := time.Parse("15:04", hours.Close)
		if err != nil {
			continue
		}

		dayStart := time.Date(day.Year(), day.Month(), day.Day(),
			openTime.Hour(), openTime.Minute(), 0, 0, day.Location())
		dayEnd := time.Date(day.Year(), day.Month(), day.Day(),
			closeTime.Hour(), closeTime.Minute(), 0, 0, day.Location())

		for slotStart := dayStart; slotStart.Add(slotDuration).Before(dayEnd) || slotStart.Add(slotDuration).Equal(dayEnd); slotStart = slotStart.Add(slotDuration + bufferDuration) {
			slotEnd := slotStart.Add(slotDuration)
			available := true
			for _, r := range reservations {
				if r.StartsAt.Before(slotEnd) && r.EndsAt.After(slotStart) {
					available = false
					break
				}
			}
			slots = append(slots, Slot{
				StartsAt:  slotStart,
				EndsAt:    slotEnd,
				Available: available,
			})
		}
	}

	sort.Slice(slots, func(i, j int) bool {
		return slots[i].StartsAt.Before(slots[j].StartsAt)
	})

	return AvailabilityResult{
		ResourceID: resourceID.String(),
		From:       from.Format(time.RFC3339),
		To:         to.Format(time.RFC3339),
		Slots:      slots,
	}, nil
}

func dayNameLower(d time.Weekday) string {
	names := [7]string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}
	return names[d]
}
