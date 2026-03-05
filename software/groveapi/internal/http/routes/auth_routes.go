package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"groveapi/internal/http/handlers"
	"groveapi/internal/http/middleware"
	supabase "github.com/supabase-community/supabase-go"
)

func RegisterAuthRoutes(r fiber.Router, db *gorm.DB, sb *supabase.Client) {
	h := handlers.NewAuthHTTP(db, sb)

	g := r.Group("/auth")
	g.Post("/login", h.Login)
	g.Post("/logout", h.Logout)
	g.Post("/register", h.Register)

	u := r.Group("/users")
	u.Use(middleware.RequireAuth(sb, db))
	u.Get("/me", h.GetCurrentUser)

	admin := r.Group("/admin/users")
	admin.Use(middleware.RequireAuth(sb, db))
	admin.Use(middleware.RequireAdmin())
	admin.Get("/", h.ListUsers)
	admin.Patch("/:id", h.UpdateUser)
	admin.Delete("/:id", h.DeactivateUser)
	admin.Post("/", h.AdminCreateUser)
}
