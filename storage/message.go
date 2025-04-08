package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Message struct {
	ID                 int64
	OriginalMessageID  int   `gorm:"not null"`
	ForwardedMessageID int   `gorm:"not null"`
	ChatID             int64 `gorm:"not null"`
	Text               string
	CreatedAt          time.Time
}

func (s *Storage) SaveMessage(msg Message) error {
	query := `
		INSERT INTO messages (original_message_id, forwarded_message_id, chat_id, text, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, msg.OriginalMessageID, msg.ForwardedMessageID, msg.ChatID, msg.Text, msg.CreatedAt)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (s *Storage) GetMessageByForwardedID(forwardedID int64) (*Message, error) {
	query := `
		SELECT id, original_message_id, forwarded_message_id, chat_id, text, created_at, updated_at
		FROM messages
		WHERE forwarded_message_id = ?
	`
	msg := &Message{}
	err := s.db.QueryRow(query, forwardedID).Scan(msg)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("queryRow: %w", err)
	}
	return msg, nil
}
