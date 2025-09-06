package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"groveapi/internal/http/handlers"
)

func RegisterReservationRoutes(r fiber.Router, db *gorm.DB) {
	h := handlers.NewReservationHTTP(db)

	g := r.Group("/reservations")
	g.Get("/my", h.ListCurrentUserReservations)
	g.Get("/", h.ListResourceReservations)
	g.Post("/",  h.CreateReservation)
	g.Patch("/:id", h.UpdateReservation)
	g.Delete("/:id", h.DeleteReservation)
}
