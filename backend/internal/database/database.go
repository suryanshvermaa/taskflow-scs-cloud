package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the standard sql.DB for SQLite.
type DB struct {
	*sql.DB
}

// Open creates or opens the SQLite database at the given path.
// It automatically creates the parent directory and the tasks table if needed.
func Open(path string) (*DB, error) {
	// Ensure the parent directory exists.
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("create db directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite3", path+"?_journal=WAL&_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	// SQLite works best with a single writer connection.
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return &DB{db}, nil
}

// migrate creates the tasks table if it does not already exist.
func migrate(db *sql.DB) error {
	const schema = `
	CREATE TABLE IF NOT EXISTS tasks (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		title       TEXT    NOT NULL,
		description TEXT    NOT NULL DEFAULT '',
		completed   INTEGER NOT NULL DEFAULT 0,
		created_at  TEXT    NOT NULL,
		updated_at  TEXT    NOT NULL
	);`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("create tasks table: %w", err)
	}

	return nil
}

