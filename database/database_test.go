package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensure the test can open & close.
func TestDB(t *testing.T) {
	db, err := SetupTestDB()
	assert.NoError(t, err)

	err = CloseTestDB(db)
	assert.NoError(t, err)
}

func SetupTestDB() (*DB, error) {
	dsn := ":memory:"
	db := NewDB(dsn)
	if err := db.Open(); err != nil {
		return nil, err
	}

	return db, nil
}

func CloseTestDB(db *DB) error {
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
