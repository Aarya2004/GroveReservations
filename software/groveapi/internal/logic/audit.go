package logic

import (
	"context"

	"groveapi/internal/store"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func WriteAuditLog(ctx context.Context, db *gorm.DB, userID uuid.UUID, action string, resourceID *uuid.UUID, details string) error {
	log := store.AuditLog{
		UserID:     userID,
		Action:     action,
		ResourceID: resourceID,
		Details:    details,
	}
	return db.WithContext(ctx).Create(&log).Error
}

type AuditLogDTO struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	Action     string     `json:"action"`
	ResourceID *string    `json:"resource_id,omitempty"`
	Details    string     `json:"details,omitempty"`
	Timestamp  string     `json:"timestamp"`
}

type AuditQueryParams struct {
	UserID string
	Action string
	From   string
	To     string
	Page   int
	Limit  int
}

func ListAuditLogs(ctx context.Context, db *gorm.DB, params AuditQueryParams) ([]AuditLogDTO, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}

	query := db.WithContext(ctx).Model(&store.AuditLog{})

	if params.UserID != "" {
		query = query.Where("user_id = ?", params.UserID)
	}
	if params.Action != "" {
		query = query.Where("action = ?", params.Action)
	}
	if params.From != "" {
		query = query.Where("timestamp >= ?", params.From)
	}
	if params.To != "" {
		query = query.Where("timestamp <= ?", params.To)
	}

	var total int64
	query.Count(&total)

	var logs []store.AuditLog
	offset := (params.Page - 1) * params.Limit
	if err := query.Order("timestamp DESC").Offset(offset).Limit(params.Limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	dtos := make([]AuditLogDTO, len(logs))
	for i, l := range logs {
		dto := AuditLogDTO{
			ID:        l.ID.String(),
			UserID:    l.UserID.String(),
			Action:    l.Action,
			Details:   l.Details,
			Timestamp: l.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		}
		if l.ResourceID != nil {
			s := l.ResourceID.String()
			dto.ResourceID = &s
		}
		dtos[i] = dto
	}
	return dtos, total, nil
}
