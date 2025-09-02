package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2"

	"groveapi/internal/config"
	httpx "groveapi/internal/http"
	"groveapi/internal/store"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	st := store.MustOpen(cfg.DatabaseURL)

	app := fiber.New(fiber.Config{ CaseSensitive: true, AppName: "groveapi" })
	httpx.RegisterRoutes(app, st)

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := app.Listen(cfg.HTTPAddr); err != nil {
		log.Fatal(err)
	}
}