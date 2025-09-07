package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"groveapi/internal/http/handlers"
	supabase "github.com/supabase-community/supabase-go"
)

func RegisterReservationRoutes(r fiber.Router, db *gorm.DB, sb *supabase.Client) {
	h := handlers.NewReservationHTTP(db, sb)

	g := r.Group("/reservations")
	g.Get("/me", h.ListCurrentUserReservations)
	g.Get("/:id", h.ListResourceReservations)
	g.Post("/",  h.CreateReservation)
	g.Patch("/:id", h.UpdateReservation)
	g.Delete("/:id", h.DeleteReservation)
}
