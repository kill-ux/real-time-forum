package models

import (
	"errors"
)

// Likes represents a like or dislike on a post or comment
type Likes struct {
	ID       int64  `json:"id"`
	UserID   int    `json:"user_id"`
	P_ID     *int64 `json:"p_id"` // Post ID (pointer to distinguish null)
	C_ID     *int64 `json:"c_id"` // Comment ID (pointer to distinguish null)
	Like     int    `json:"like"` // 1 for like, -1 for dislike
	LikeType string `json:"like_type"`
	NameID   string `json:"name_id"` // "post_id" or "comment_id"
}

// GetLikes represents like/dislike counts and user's like status
type GetLikes struct {
	Likes    int    `json:"likes"`
	DisLikes int    `json:"dislikes"`
	Like     *int   `json:"like"` // User's like status (pointer to handle null)
	NameID   string `json:"name_id"`
	P_ID     *int64 `json:"p_id"`
	C_ID     *int64 `json:"c_id"`
}

// BeforCreateLikes validates like data before insertion
func (like *Likes) BeforCreateLikes() error {
	// Validate NameID is either "comment_id" or "post_id"
	if like.NameID != "comment_id" && like.NameID != "post_id" {
		return errors.New("invalid NameID")
	}
	// Validate Like value is 1 (like) or -1 (dislike)
	if like.Like != 1 && like.Like != -1 {
		return errors.New("invalid Like value")
	}
	// Set LikeType based on NameID
	like.LikeType = like.NameID[:len(like.NameID)-len("_id")] + "s"
	return nil
}
