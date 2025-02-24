package models

import (
	"errors"
)

type Likes struct {
	ID       int64  `json:"id"`
	UserID   int    `json:"user_id"`
	P_ID     *int64 `json:"p_id"` // Use pointer to distinguish between zero value and nil
	C_ID     *int64 `json:"c_id"` // Use pointer to distinguish between zero value and nil
	Like     int    `json:"like"` // 1 for like, -1 for dislike
	LikeType string `json:"like_type"`
	NameID   string `json:"name_id"`
}

type GetLikes struct {
	Likes    int    `json:"likes"`
	DisLikes int    `json:"dislikes"`
	Like     *int   `json:"like"` // Pointer to handle null values
	NameID   string `json:"name_id"`
	P_ID     *int64 `json:"p_id"`
	C_ID     *int64 `json:"c_id"`
}

func (like *Likes) BeforCreateLikes() error {
	if like.NameID != "comment_id" && like.NameID != "post_id" {
		return errors.New("invalid NameID")
	}
	if like.Like != 1 && like.Like != -1 {
		return errors.New("invalid Like value")
	}
	like.LikeType = like.NameID[:len(like.NameID)-len("_id")] + "s"
	return nil
}
