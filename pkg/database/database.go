package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
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
		if err := os.MkdirAll(filepath.Dir(db.DSN), 0o700); err != nil {
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

// QueryRow runs the query and returns a single row.
func (db *DB) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return db.db.QueryRowContext(ctx, query, args...)
}
