package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"groveapi/internal/http/handlers"
	"groveapi/internal/http/middleware"
	supabase "github.com/supabase-community/supabase-go"
)

func RegisterResourceRoutes(r fiber.Router, db *gorm.DB, sb *supabase.Client) {
	h := handlers.NewResourceHTTP(db)

	g := r.Group("/resources")
	g.Get("/", h.ListResources)
	g.Get("/:id", h.GetResource)

	// Admin-only
	admin := r.Group("/resources")
	admin.Use(middleware.RequireAuth(sb, db))
	admin.Use(middleware.RequireAdmin())
	admin.Post("/", h.CreateResource)
	admin.Patch("/:id", h.UpdateResource)
	admin.Delete("/:id", h.DeleteResource)
}
