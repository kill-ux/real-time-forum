package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"forum/db"
	"forum/models"
	"forum/utils"
	"forum/utils/middlewares"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	if err := post.BeforCreatePost(r); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest)
		return
	}

	newCats := isValidCategory(r.Form["categories"])
	post.Categories = strings.Join(newCats, ", ")

	// Process image if exists
	file, fileheader, err := r.FormFile("image")
	// var image string
	if err == nil {
		post.Image = models.HandleImage("posts", file, fileheader)
	}

	user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized)
		return
	}

	post.UserID = user.ID
	post.CreatedAt = int(time.Now().Unix())
	result, err := db.DB.Exec(`INSERT INTO posts 
        VALUES (Null,?,?,?,?,?,?)`, utils.GetExecFields(post, "ID")...)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}
	post.ID, err = result.LastInsertId()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}

	post_user := models.PostWithUser{
		Post: post,
		User: user,
	}
	utils.RespondWithJSON(w, http.StatusCreated, post_user)
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	var before struct {
		Before int `json:"before"`
	}

	err := utils.ParseBody(r, &before)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest)
	}

	rows, err := db.DB.Query(`
        SELECT p.*,u.id ,u.first_name,u.last_name,u.nickname,u.image
        FROM posts p
        JOIN users u ON p.user_id = u.id
		WHERE p.created_at < ?
        ORDER BY p.created_at DESC
        LIMIT 10`, before.Before)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.PostWithUser
	for rows.Next() {
		var post models.PostWithUser
		if err := rows.Scan(append(utils.GetScanFields(&post.Post), &post.User.ID, &post.User.FirstName, &post.User.LastName, &post.User.Nickname, &post.User.Image)...); err != nil {
			fmt.Println(err)
			utils.RespondWithError(w, http.StatusInternalServerError)
			return
		}
		post.P_ID = &post.ID
		post.NameID = "post_id"
		post.GetLikes, err = GetLikes(w, r, post.GetLikes)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}

	// var likes models.GetLikes
	// likes.P_ID = like.P_ID
	// likes.NameID = like.NameID

	// utils.RespondWithJSON(w, http.StatusCreated, likes)

	utils.RespondWithJSON(w, http.StatusOK, posts)
}

func getValidCategories() []string {
	return []string{"tech", "programming", "health", "finance", "food", "science", "memes", "others"}
}

func isValidCategory(categories []string) (newCats []string) {
	validCategories := getValidCategories()
	for _, validCat := range validCategories {
		for _, category := range categories {
			if validCat == category {
				newCats = append(newCats, validCat)
			}
		}
	}
	if len(newCats) == 0 {
		newCats = append(newCats, "others")
	}
	return
}
