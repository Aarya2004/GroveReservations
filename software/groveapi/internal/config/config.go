package config

import (
	"log"
	"os"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	SupabaseUrl string
	SupabaseServiceKey string
}

func Load() Config {
	return Config{
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL: must("DATABASE_URL"),
		SupabaseUrl: must("SUPABASE_URL"),
		SupabaseServiceKey: must("SUPABASE_SERVICE_KEY"),
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func must(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}
