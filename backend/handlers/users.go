package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"forum/db"
	"forum/models"
	"forum/utils"
	"forum/utils/middlewares"
)

// GetMessageHistoryHandler retrieves the message history between the authenticated user and another user.
// It fetches messages created before a specified timestamp for pagination.
// Returns 200 with an array of messages on success, or appropriate error codes on failure.
func GetMessageHistoryHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := user.ID
	var receiver_id struct {
		ReceiverID int `json:"receiver_id"`
		Before     int `json:"before"`
	}
	err := utils.ParseBody(r, &receiver_id)
	if err != nil {
		log.Println("Failed to parse request body:", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if receiver_id.Before == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	messages, err := models.GetMessageHistory(userID, receiver_id.ReceiverID, receiver_id.Before)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, messages)
}

// GetUsers retrieves a list of all users except the specified user, with messaging metadata.
// It includes last message time and unread message count for each user.
// Returns a slice of Members or an error if the query fails.
func GetUsers(userID int) (users []models.Members, err error) {
	query := `
        SELECT 
			u.id, 
			u.nickname, 
			u.first_name, 
			u.last_name,
			u.last_seen,
			u.image,
			COALESCE(MAX(m.created_at), 0) AS last_message_time,
			COALESCE(SUM(CASE WHEN m.is_read = false AND m.receiver_id = ? THEN 1 ELSE 0 END), 0) AS unread_count
		FROM 
			users u
		LEFT JOIN 
			messages m 
			ON (u.id = m.sender_id OR u.id = m.receiver_id)
			AND (m.sender_id = ? OR m.receiver_id = ?)
		WHERE 
			u.id != ?  
		GROUP BY 
			u.id
		ORDER BY 
			last_message_time DESC, 
			u.first_name ASC;
    `
	rows, err := db.DB.Query(query, userID, userID, userID, userID)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	// var users []models.Members
	for rows.Next() {
		var user models.Members
		if err = rows.Scan(utils.GetScanFields(&user)...); err != nil {
			fmt.Println(err)
			return
		}
		users = append(users, user)
	}
	return
}
