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

// CreateLikesHandler handles like/dislike actions on posts or comments.
// It creates, updates, or removes likes based on the user's previous action.
// Returns 201 with updated like counts on success, or appropriate error codes on failure.
func CreateLikesHandler(w http.ResponseWriter, r *http.Request) {
	var like models.Likes

	err := utils.ParseBody(r, &like)
	if err != nil {
		log.Println("Failed to parse request body:", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := like.BeforCreateLikes(); err != nil {
		log.Println("Validation error:", err)
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	like.UserID = user.ID

	var contentID int64
	if like.C_ID != nil {
		contentID = *like.C_ID
	} else if like.P_ID != nil {
		contentID = *like.P_ID
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "Either P_ID or C_ID must be provided")
		return
	}

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
			// Remove like/dislike if it's the same as the existing one
			result, err = db.DB.Exec(`DELETE FROM likes WHERE user_id = ? AND `+like.NameID+` = ?`, user.ID, contentID)
		} else {
			// Update like/dislike
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

// GetLikesHandler retrieves like/dislike counts and the current user's like status.
// It works for both posts and comments based on the provided ID.
// Returns 200 with like data on success, or 400/401 on errors.
func GetLikesHandler(w http.ResponseWriter, r *http.Request) {
	var likes models.GetLikes
	err := utils.ParseBody(r, &likes)
	if err != nil || (likes.NameID != "comment_id" && likes.NameID != "post_id") || (likes.P_ID != nil && likes.C_ID != nil) {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	likes, err = GetLikes(w, r, likes)
	if err != nil {
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, likes)
}

// GetLikes is a helper function that fetches like/dislike counts and user's like status.
// It queries the database for total likes, dislikes, and the authenticated user's vote.
// Returns the populated GetLikes struct or an error if the operation fails.
func GetLikes(w http.ResponseWriter, r *http.Request, likes models.GetLikes) (models.GetLikes, error) {
	user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return models.GetLikes{}, errors.New("unth")
	}
	userID := user.ID

	var contentID int64
	if likes.C_ID != nil {
		contentID = *likes.C_ID
	} else if likes.P_ID != nil {
		contentID = *likes.P_ID
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "Either P_ID or C_ID must be provided")
		return models.GetLikes{}, errors.New("unth")
	}

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

	if userLike.Valid {
		likes.Like = new(int)
		*likes.Like = int(userLike.Int64)
	}
	return likes, nil

	// utils.RespondWithJSON(w, http.StatusOK, likes)
}
