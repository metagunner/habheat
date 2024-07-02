package database

import (
	"database/sql"
	"embed"
	"fmt"
	"habheath"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

// DB represents the database connection.
type DB struct {
	db *sql.DB
	// Datasource name.
	DSN string
}

//go:embed migration/*.sql
var embedMigrations embed.FS

// NewDB returns a new instance of DB associated with the given datasource name.
func NewDB(dsn string) *DB {
	db := &DB{
		DSN: dsn,
	}
	return db
}

// Open creates a new DB for the given connection string.
func (db *DB) Open() (err error) {
	if db.DSN == "" {
		return fmt.Errorf("dsn required")
	}

	if db.DSN != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(db.DSN), 0700); err != nil {
			return err
		}
	}

	if db.db, err = sql.Open("sqlite3", db.DSN); err != nil {
		return err
	}

	if _, err := db.db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("foreign keys pragma: %w", err)
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}

	if err := goose.Up(db.db, "migration"); err != nil {
		panic(err)
	}

	return nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

// FormatError returns err as a error, if possible.
// Otherwise returns the original error.
func FormatError(err error) error {
	if err == nil {
		return nil
	}

	switch err.Error() {
	case "UNIQUE constraint failed: dial_memberships.dial_id, dial_memberships.user_id":
		return habheath.Errorf(habheath.ECONFLICT, "Dial membership already exists.")
	default:
		return err
	}
}
