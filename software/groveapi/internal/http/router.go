package http

import (
	"groveapi/internal/http/routes"

	"github.com/gofiber/fiber/v2"
	supabase "github.com/supabase-community/supabase-go"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB, sb *supabase.Client) *fiber.App {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		AppName:       "groveapi",
	})

	// Global middleware (logging, recover, request ID) – add later if you want
	app.Get("/health", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status":"ok"}) })

	// Versioned API groups
	api := app.Group("/api")
	v1  := api.Group("/v1")

	// Mount feature routes
	routes.RegisterAuthRoutes(v1, db, sb)
	routes.RegisterResourceRoutes(v1, db)
	routes.RegisterReservationRoutes(v1, db, sb)

	return app
}

// func RegisterRoutes(app *fiber.App, s *store.Store) {
// 	app.Get("/health", func(c *fiber.Ctx) error {
// 		return c.JSON(fiber.Map{"status": "ok"})
// 	})

// 	app.Post("/reservations", func(c *fiber.Ctx) error {
// 		log.Printf("Received request: %s %s", c.Method(), c.Body())
// 		var in logic.ReservationInput
// 		if err := c.BodyParser(&in); err != nil {
// 			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"invalid_json"})
// 		}
// 		res, err := logic.CreateReservation(c.Context(), s.DB, in)
// 		switch {
// 		case err == nil:
// 			return c.Status(http.StatusCreated).JSON(res)
// 		case errors.Is(err, logic.ErrConflict):
// 			return c.Status(http.StatusConflict).JSON(fiber.Map{"error":"conflict"})
// 		case errors.Is(err, logic.ErrRuleViolation):
// 			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error":"rule_violation"})
// 		default:
// 			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error":"internal"})
// 		}
// 	})
// }
