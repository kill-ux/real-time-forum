// models/message.go
package models

import (
	"database/sql"
	"fmt"

	"forum/db"
	"forum/utils"
)

type Message struct {
	ID         int    `json:"id"`
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
	CreatedAt  int    `json:"created_at"`
	IsRead     bool   `json:"is_read"`
}

type WSMessage struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Message `json:"message"`
	Members []Members `json:"members"`
	Typing  bool      `json:"is_typing"`
}

// GetMessageHistory retrieves the message history between two users.
// It fetches up to 10 messages created before the specified timestamp, ordered by creation time.
// Returns a slice of messages or an error if the query fails.
func GetMessageHistory(sender, receiver, time int) (messages []Message, err error) {
	rows, err := db.DB.Query(`
	    SELECT * FROM messages
	    WHERE created_at < ? AND ((sender_id = ? AND receiver_id = ?)
	    OR (sender_id = ? AND receiver_id = ?))
	    ORDER BY created_at DESC
	    LIMIT 10`, time, sender, receiver, receiver, sender)
	// Process rows
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return
	}
	// fmt.Println("rows", rows)

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

// StoreMessage inserts a new message into the database.
// It excludes the auto-generated ID field from the insert operation.
// Returns an error if the database operation fails.
func (message *Message) StoreMessage() error {
	// message.CreatedAt = int(time.Now().Unix())
	_, err := db.DB.Exec("INSERT INTO messages VALUES(NULL ,?,?,?,?,?)", utils.GetExecFields(message, "ID")...)
	return err
}

// UpdateRead marks all unread messages from a specific sender to a receiver as read.
// This is called when a user views their conversation with another user.
// Returns an error if the database operation fails.
func (message *Message) UpdateRead() error {
	_, err := db.DB.Exec("UPDATE messages SET is_read = true WHERE sender_id = ? AND receiver_id = ? AND is_read = false", message.ReceiverID, message.SenderID)
	return err
}
