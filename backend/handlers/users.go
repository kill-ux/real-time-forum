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

// GetMessageHistoryHandler retrieves paginated message history between two users
func GetMessageHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := user.ID
	
	// Parse receiver ID and pagination parameter
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
	
	// Validate pagination parameter
	if receiver_id.Before == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	
	// Fetch message history
	messages, err := models.GetMessageHistory(userID, receiver_id.ReceiverID, receiver_id.Before)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, messages)
}

// GetUsers retrieves all users with their last message time and unread count for the current user
func GetUsers(userID int) (users []models.Members, err error) {
	// Query users with message metadata sorted by last message time
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

	// Build users array
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
