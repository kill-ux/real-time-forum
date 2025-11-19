package models

import (
	"errors"
)

// Comment represents a comment on a post
type Comment struct {
	ID        int64  `json:"id"`
	UserID    int    `json:"user_id"`
	PostID    int    `json:"post_id"`
	Content   string `json:"content"`
	CreatedAt int    `json:"created_at"`
}

// CommentWithUser combines comment data with user and like information
type CommentWithUser struct {
	Comment
	GetLikes
	User Members `json:"user"`
}

// BeforCreateComment validates comment data before insertion
func (comment *Comment) BeforCreateComment() error {
	// Validate content length and post ID
	if length(1, 2000, comment.Content) || comment.PostID == 0 {
		return errors.New("size not allowed")
	}
	return nil
}
