package models

import (
	"database/sql"
	"fmt"

	"forum/db"
	"forum/utils"
)

// Message represents a private message between users
type Message struct {
	ID         int    `json:"id"`
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
	CreatedAt  int    `json:"created_at"`
	IsRead     bool   `json:"is_read"`
}

// WSMessage represents a WebSocket message with metadata
type WSMessage struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Message `json:"message"`
	Members []Members `json:"members"`
	Typing  bool      `json:"is_typing"`
}

// GetMessageHistory retrieves paginated message history between two users
func GetMessageHistory(sender, receiver, time int) (messages []Message, err error) {
	// Query messages between sender and receiver before specified time
	rows, err := db.DB.Query(`
	    SELECT * FROM messages
	    WHERE created_at < ? AND ((sender_id = ? AND receiver_id = ?)
	    OR (sender_id = ? AND receiver_id = ?))
	    ORDER BY created_at DESC
	    LIMIT 10`, time, sender, receiver, receiver, sender)
	
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return
	}

	// Build messages array
	for rows.Next() {
		var message Message
		if err = rows.Scan(utils.GetScanFields(&message)...); err != nil {
			fmt.Println(err)
			return
		}
		messages = append(messages, message)
	}
	return
}

// StoreMessage inserts a new message into the database
func (message *Message) StoreMessage() error {
	_, err := db.DB.Exec("INSERT INTO messages VALUES(NULL ,?,?,?,?,?)", utils.GetExecFields(message, "ID")...)
	return err
}

// UpdateRead marks all unread messages from a sender as read
func (message *Message) UpdateRead() error {
	_, err := db.DB.Exec("UPDATE messages SET is_read = true WHERE sender_id = ? AND receiver_id = ? AND is_read = false", message.ReceiverID, message.SenderID)
	return err
}
