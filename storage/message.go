package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Message struct {
	ID                 int64
	OriginalMessageID  int
	ForwardedMessageID int
	ChatID             int64
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

func (s *Storage) GetMessageByForwardedID(forwardedID int) (*Message, error) {
	query := `
		SELECT original_message_id, forwarded_message_id, chat_id
		FROM messages
		WHERE forwarded_message_id = ?
	`
	var msg Message
	err := s.db.QueryRow(query, forwardedID).Scan(&msg.OriginalMessageID, &msg.ForwardedMessageID, &msg.ChatID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("queryRow: %w", err)
	}
	return &msg, nil
}
