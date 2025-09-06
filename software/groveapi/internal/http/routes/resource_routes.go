package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"groveapi/internal/http/handlers"
	// "groveapi/internal/http/middleware" // e.g., RequireAdmin
)

func RegisterResourceRoutes(r fiber.Router, db *gorm.DB) {
	h := handlers.NewResourceHTTP(db)

	g := r.Group("/resources")
	g.Get("/", h.ListResources)
	g.Get("/:id", h.GetResource)

	// Admin-only
	// g.Use(middleware.RequireAdmin())
	g.Post("/",    h.CreateResource)
	g.Patch("/:id",h.UpdateResource)
	g.Delete("/:id",h.DeleteResource)
}
