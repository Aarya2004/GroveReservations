package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"groveapi/internal/http/handlers"
	"groveapi/internal/http/middleware"
	supabase "github.com/supabase-community/supabase-go"
)

func RegisterReservationRoutes(r fiber.Router, db *gorm.DB, sb *supabase.Client) {
	h := handlers.NewReservationHTTP(db)

	g := r.Group("/reservations")
	g.Use(middleware.RequireAuth(sb, db))
	g.Get("/me", h.ListCurrentUserReservations)
	g.Get("/", h.ListResourceReservations)
	g.Post("/", h.CreateReservation)
	g.Patch("/:id", h.UpdateReservation)
	g.Delete("/:id", h.DeleteReservation)
}
