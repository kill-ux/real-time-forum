package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"forum/db"
	"forum/models"
	"forum/utils"
	"forum/utils/middlewares"
)

// CreateLikesHandler handles like/dislike creation, update, and removal
func CreateLikesHandler(w http.ResponseWriter, r *http.Request) {
	var like models.Likes

	// Parse like data from request body
	err := utils.ParseBody(r, &like)
	if err != nil {
		log.Println("Failed to parse request body:", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate like data
	if err := like.BeforCreateLikes(); err != nil {
		log.Println("Validation error:", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get authenticated user from context
	user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	like.UserID = user.ID

	// Determine content ID (post or comment)
	var contentID int64
	if like.C_ID != nil {
		contentID = *like.C_ID
	} else if like.P_ID != nil {
		contentID = *like.P_ID
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "Either P_ID or C_ID must be provided")
		return
	}

	// Check if user already liked/disliked this content
	var existingLike int
	query := "SELECT like FROM likes WHERE " + like.NameID + " = ? AND user_id = ?"
	row := db.DB.QueryRow(query, contentID, user.ID)
	err = row.Scan(&existingLike)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Failed to query existing like:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	var result sql.Result
	if err == sql.ErrNoRows {
		// Insert new like/dislike
		result, err = db.DB.Exec(`INSERT INTO likes (user_id, `+like.NameID+`, like) VALUES (?, ?, ?)`, user.ID, contentID, like.Like)
	} else {
		if existingLike == like.Like {
			// Remove like/dislike if clicking the same button again
			result, err = db.DB.Exec(`DELETE FROM likes WHERE user_id = ? AND `+like.NameID+` = ?`, user.ID, contentID)
		} else {
			// Update like to dislike or vice versa
			result, err = db.DB.Exec(`UPDATE likes SET like = ? WHERE user_id = ? AND `+like.NameID+` = ?`, like.Like, user.ID, contentID)
		}
	}

	if err != nil {
		log.Println("Failed to execute SQL query:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	like.ID, err = result.LastInsertId()
	if err != nil {
		log.Println("Failed to get last insert ID:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Return updated like counts
	var likes models.GetLikes
	likes.C_ID = like.C_ID
	likes.P_ID = like.P_ID
	likes.NameID = like.NameID
	likes, err = GetLikes(w, r, likes)
	if err != nil {
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, likes)
}

// GetLikesHandler retrieves like/dislike counts for a post or comment
func GetLikesHandler(w http.ResponseWriter, r *http.Request) {
	var likes models.GetLikes
	
	// Parse and validate request
	err := utils.ParseBody(r, &likes)
	if err != nil || (likes.NameID != "comment_id" && likes.NameID != "post_id") || (likes.P_ID != nil && likes.C_ID != nil) {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	
	// Fetch like counts
	likes, err = GetLikes(w, r, likes)
	if err != nil {
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, likes)
}

// GetLikes retrieves like/dislike counts and user's like status for a post or comment
func GetLikes(w http.ResponseWriter, r *http.Request, likes models.GetLikes) (models.GetLikes, error) {
	// Get authenticated user from context
	user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return models.GetLikes{}, errors.New("unth")
	}
	userID := user.ID

	// Determine content ID (post or comment)
	var contentID int64
	if likes.C_ID != nil {
		contentID = *likes.C_ID
	} else if likes.P_ID != nil {
		contentID = *likes.P_ID
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "Either P_ID or C_ID must be provided")
		return models.GetLikes{}, errors.New("unth")
	}

	// Query like counts and user's like status
	query := `
		SELECT 
			(SELECT COUNT(*) FROM likes WHERE ` + likes.NameID + ` = ? AND like = 1) as likes,
			(SELECT COUNT(*) FROM likes WHERE ` + likes.NameID + ` = ? AND like = -1) as dislikes,
			(SELECT like FROM likes WHERE ` + likes.NameID + ` = ? AND user_id = ?)
	`
	row := db.DB.QueryRow(query, contentID, contentID, contentID, userID)

	var userLike sql.NullInt64
	if err := row.Scan(&likes.Likes, &likes.DisLikes, &userLike); err != nil {
		log.Println("Failed to scan row:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return models.GetLikes{}, errors.New("unth")
	}

	// Set user's like status if they have liked/disliked
	if userLike.Valid {
		likes.Like = new(int)
		*likes.Like = int(userLike.Int64)
	}
	return likes, nil
}
