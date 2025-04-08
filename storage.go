package main

import (
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage() *Storage {
	dbPath, ok := os.LookupEnv("DB_PATH")
	if !ok {
		log.Fatal("DB_PATH environment variable is not set")
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Printf("gorm.Open: %s", err)
		return nil
	}

	err = db.AutoMigrate(&Message{})
	if err != nil {
		log.Printf("autoMigrate: %s", err)
		return nil
	}

	return &Storage{db: db}
}

// Message represents a forwarded message in the database
type Message struct {
	ID                 uint `gorm:"primarykey"`
	CreatedAt          time.Time
	OriginalMessageID  int   `gorm:"not null"`
	ForwardedMessageID int   `gorm:"not null"`
	ChatID             int64 `gorm:"not null"`
	Text               string
}

func (s *Storage) SaveMessage(msg Message) error {
	return s.db.Create(msg).Error
}

func (s *Storage) GetMessageByForwardedID(forwardedID int64) (*Message, error) {
	var msg Message
	err := s.db.Where("forwarded_message_id = ?", forwardedID).First(&msg).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (s *Storage) GetMessageByOriginalID(originalID int64) (*Message, error) {
	var msg Message
	err := s.db.Where("original_message_id = ?", originalID).First(&msg).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
