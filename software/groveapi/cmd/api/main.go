package main

import (
	"log"

	"github.com/joho/godotenv"

	"groveapi/internal/config"
	http "groveapi/internal/http"
	"groveapi/internal/store"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	st := store.MustOpen(cfg.DatabaseURL)

	app := http.NewApp(st.DB)

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := app.Listen(cfg.HTTPAddr); err != nil {
		log.Fatal(err)
	}
}