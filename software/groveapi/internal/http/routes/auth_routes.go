package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"groveapi/internal/http/handlers"
	supabase "github.com/supabase-community/supabase-go"
)

func RegisterAuthRoutes(r fiber.Router, db *gorm.DB, sb *supabase.Client) {
	h := handlers.NewAuthHTTP(db, sb)

	g := r.Group("/auth")
	g.Post("/login",  h.Login)
	g.Post("/logout", h.Logout)
	g.Post("/register", h.Register)

	u := r.Group("/users")
	u.Get("/me", h.GetCurrentUser)

	admin := r.Group("/admin/users") // TODO: Need to modify this to include identity in request
	admin.Get("/", h.ListUsers)
	admin.Patch("/:id", h.UpdateUser)
	admin.Delete("/:id", h.DeactivateUser)
	admin.Post("/", h.AdminCreateUser)
}
