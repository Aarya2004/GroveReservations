package store

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Gorm store holds DB
type Store struct{ DB *gorm.DB }

func MustOpen(dsn string) *Store {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("gorm open: %v", err)
	}
	return &Store{DB: db}
}
