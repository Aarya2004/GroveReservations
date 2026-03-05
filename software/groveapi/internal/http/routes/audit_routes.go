package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"groveapi/internal/http/handlers"
	"groveapi/internal/http/middleware"
	supabase "github.com/supabase-community/supabase-go"
)

func RegisterAuditRoutes(r fiber.Router, db *gorm.DB, sb *supabase.Client) {
	h := handlers.NewAuditHTTP(db)

	g := r.Group("/audit")
	g.Use(middleware.RequireAuth(sb, db))
	g.Use(middleware.RequireAdmin())
	g.Get("/", h.ListAuditLogs)
}
