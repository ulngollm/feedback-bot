package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/glebarez/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	dbPath, ok := os.LookupEnv("DB_PATH")
	if !ok {
		log.Fatal("DB_PATH environment variable is not set")
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			original_message_id INTEGER NOT NULL,
			forwarded_message_id INTEGER NOT NULL,
			chat_id INTEGER NOT NULL,
			text TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("error creating messages table: %w", err)
	}

	return &Storage{db: db}, nil
}
