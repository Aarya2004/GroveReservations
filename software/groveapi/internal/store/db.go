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

func (s *Store) Ping() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
