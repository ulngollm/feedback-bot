package storage

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/glebarez/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New() *Storage {
	dbPath, ok := os.LookupEnv("DB_PATH")
	if !ok {
		log.Println("DB_PATH environment variable is not set")
		return nil
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Printf("error opening database: %v", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		log.Printf("error connecting to the database: %v", err)
		return nil
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
		log.Printf("error creating messages table: %v", err)
		return nil
	}

	return &Storage{db: db}
}
