package models

type Members struct {
	ID              int    `json:"id"`
	Nickname        string `json:"nickname"`
	FirstName       string `json:"firstname"`
	LastName        string `json:"lastname"`
	LastSeen        int64  `json:"last_seen"`
	Image           string `json:"image"`
	LastMessageTime string `json:"last_message_time"`
	UnreadCount     int    `json:"unread_count"`
}
