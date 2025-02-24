package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"forum/db"
	"forum/models"
	"forum/utils"
	"forum/utils/middlewares"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var comment models.CommentWithUser

	err := utils.ParseBody(r, &comment.Comment)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}

	if err := comment.Comment.BeforCreateComment(); err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := user.ID

	comment.Comment.UserID = userID
	comment.Comment.CreatedAt = int(time.Now().Unix())

	result, err := db.DB.Exec(`INSERT INTO comments VALUES (NULL,?,?,?,?)`, utils.GetExecFields(comment.Comment, "ID")...)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}

	comment.Comment.ID, err = result.LastInsertId()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}
	comment.User.ID = user.ID
	comment.User.Nickname = user.Nickname
	comment.User.FirstName = user.FirstName
	comment.User.LastName = user.LastName
	comment.User.Image = user.Image
	utils.RespondWithJSON(w, http.StatusCreated, comment)
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Post_id int `json:"post_id"`
		Before  int `json:"before"`
	}

	err := utils.ParseBody(r, &data)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}
	rows, err := db.DB.Query(`
        SELECT c.*,u.id ,u.first_name,u.last_name,u.nickname,u.image
        FROM comments c
        JOIN users u ON c.user_id = u.id
		WHERE post_id = ? AND c.created_at < ?
        ORDER BY c.created_at DESC
        LIMIT 10`, data.Post_id, data.Before)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []models.CommentWithUser
	for rows.Next() {
		var comment models.CommentWithUser
		if err := rows.Scan(append(utils.GetScanFields(&comment.Comment), &comment.User.ID, &comment.User.FirstName, &comment.User.LastName, &comment.User.Nickname, &comment.User.Image)...); err != nil {
			fmt.Println(err)
			utils.RespondWithError(w, http.StatusInternalServerError)
			return
		}
		fmt.Println(comment.Comment)
		comment.GetLikes.C_ID = &comment.Comment.ID
		comment.GetLikes.NameID = "comment_id"
		comment.GetLikes, err = GetLikes(w, r, comment.GetLikes)
		if err != nil {
			return
		}
		comments = append(comments, comment)
	}

	utils.RespondWithJSON(w, http.StatusOK, comments)
}
