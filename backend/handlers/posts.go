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

// CreatePostHandler handles post creation requests
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var post models.Post

	// Validate post data from form
	if err := post.BeforCreatePost(r); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest)
		return
	}

	// Validate and filter categories
	newCats := isValidCategory(r.Form["categories"])
	post.Categories = strings.Join(newCats, ", ")

	// Handle optional image upload
	file, fileheader, err := r.FormFile("image")
	if err == nil {
		post.Image = models.HandleImage("posts", file, fileheader)
	}

	// Get authenticated user from context
	user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized)
		return
	}

	// Set post metadata and insert into database
	post.UserID = user.ID
	post.CreatedAt = int(time.Now().Unix())
	result, err := db.DB.Exec(`INSERT INTO posts 
        VALUES (Null,?,?,?,?,?,?)`, utils.GetExecFields(post, "ID")...)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}
	
	// Get the newly created post ID
	post.ID, err = result.LastInsertId()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}

	// Return post with user information
	post_user := models.PostWithUser{
		Post: post,
		User: user,
	}
	utils.RespondWithJSON(w, http.StatusCreated, post_user)
}

// GetPostsHandler retrieves paginated posts with user and like information
func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	var before struct {
		Before int `json:"before"`
	}

	// Parse pagination parameter
	err := utils.ParseBody(r, &before)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest)
	}

	// Query posts with user information, ordered by creation time
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

	// Build posts array with likes information
	var posts []models.PostWithUser
	for rows.Next() {
		var post models.PostWithUser
		if err := rows.Scan(append(utils.GetScanFields(&post.Post), &post.User.ID, &post.User.FirstName, &post.User.LastName, &post.User.Nickname, &post.User.Image)...); err != nil {
			fmt.Println(err)
			utils.RespondWithError(w, http.StatusInternalServerError)
			return
		}
		
		// Fetch likes for each post
		post.P_ID = &post.ID
		post.NameID = "post_id"
		post.GetLikes, err = GetLikes(w, r, post.GetLikes)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}

	utils.RespondWithJSON(w, http.StatusOK, posts)
}

// getValidCategories returns the list of allowed post categories
func getValidCategories() []string {
	return []string{"tech", "programming", "health", "finance", "food", "science", "memes", "others"}
}

// isValidCategory filters user-provided categories against valid ones
func isValidCategory(categories []string) (newCats []string) {
	validCategories := getValidCategories()
	for _, validCat := range validCategories {
		for _, category := range categories {
			if validCat == category {
				newCats = append(newCats, validCat)
			}
		}
	}
	// Default to "others" if no valid categories provided
	if len(newCats) == 0 {
		newCats = append(newCats, "others")
	}
	return
}
