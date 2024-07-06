package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Ensure the test can open & close.
func TestDB(t *testing.T) {
	db, err := SetupTestDB()
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}
