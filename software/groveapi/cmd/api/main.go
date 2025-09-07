package main

import (
	"log"

	"github.com/joho/godotenv"

	"groveapi/internal/config"
	http "groveapi/internal/http"
	"groveapi/internal/store"
	"groveapi/internal/sb"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	st := store.MustOpen(cfg.DatabaseURL)
	sb := sb.MustNewSupabaseClient(cfg.SupabaseUrl, cfg.SupabaseServiceKey)

	app := http.NewApp(st.DB, sb)

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := app.Listen(cfg.HTTPAddr); err != nil {
		log.Fatal(err)
	}
}