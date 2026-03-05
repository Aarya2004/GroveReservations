package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"groveapi/internal/http/handlers"
	"groveapi/internal/http/middleware"
	supabase "github.com/supabase-community/supabase-go"
)

func RegisterAvailabilityRoutes(r fiber.Router, db *gorm.DB, sb *supabase.Client) {
	h := handlers.NewAvailabilityHTTP(db)

	g := r.Group("/availability")
	g.Use(middleware.RequireAuth(sb, db))
	g.Get("/", h.GetAvailability)
}
