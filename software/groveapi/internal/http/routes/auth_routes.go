package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"groveapi/internal/http/handlers"
)

func RegisterAuthRoutes(r fiber.Router, db *gorm.DB) {
	h := handlers.NewAuthHTTP(db)

	g := r.Group("/auth")
	g.Post("/login",  h.Login)
	g.Post("/logout", h.Logout)
	// Maybe add registration

	u := r.Group("/users")
	u.Get("/me", h.GetCurrentUser)
	u.Get("/", h.ListUsers)
	u.Patch("/:id", h.UpdateUser)
	u.Delete("/:id", h.DeactivateUser)
}
