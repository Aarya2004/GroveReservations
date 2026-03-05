package testutil

import (
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustOpenTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	return db
}

// WrapInTx starts a transaction and rolls it back after the test.
// All operations performed on the returned *gorm.DB are isolated.
func WrapInTx(t *testing.T, db *gorm.DB) *gorm.DB {
	t.Helper()

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	t.Cleanup(func() {
		tx.Rollback()
	})

	return tx
}
