package http

import (
	"groveapi/internal/http/routes"
	"groveapi/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	supabase "github.com/supabase-community/supabase-go"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB, sb *supabase.Client) *fiber.App {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		AppName:       "groveapi",
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Health check — pings DB
	app.Get("/health", func(c *fiber.Ctx) error {
		s := &store.Store{DB: db}
		if err := s.Ping(); err != nil {
			return c.Status(503).JSON(fiber.Map{"status": "unhealthy", "error": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Versioned API groups
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Mount feature routes
	routes.RegisterAuthRoutes(v1, db, sb)
	routes.RegisterResourceRoutes(v1, db, sb)
	routes.RegisterReservationRoutes(v1, db, sb)
	routes.RegisterAvailabilityRoutes(v1, db, sb)
	routes.RegisterAuditRoutes(v1, db, sb)

	return app
}
