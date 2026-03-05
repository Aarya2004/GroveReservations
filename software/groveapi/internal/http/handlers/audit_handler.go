package handlers

import (
	"groveapi/internal/logic"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuditHTTP struct {
	DB *gorm.DB
}

func NewAuditHTTP(db *gorm.DB) *AuditHTTP { return &AuditHTTP{DB: db} }

func (h *AuditHTTP) ListAuditLogs(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	params := logic.AuditQueryParams{
		UserID: c.Query("user_id"),
		Action: c.Query("action"),
		From:   c.Query("from"),
		To:     c.Query("to"),
		Page:   page,
		Limit:  limit,
	}

	logs, total, err := logic.ListAuditLogs(c.Context(), h.DB, params)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, "internal_error", "failed to list audit logs")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"total": total,
		"page":  page,
		"limit": limit,
		"logs":  logs,
	})
}
