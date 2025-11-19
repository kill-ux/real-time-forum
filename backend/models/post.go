package models

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gofrs/uuid/v5"
)

// Post represents a forum post
type Post struct {
	ID         int64  `json:"id"`
	UserID     int    `json:"user_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Categories string `json:"categories"`
	CreatedAt  int    `json:"created_at"`
	Image      string `json:"image"`
}

// PostWithUser combines post data with user and like information
type PostWithUser struct {
	Post
	GetLikes
	User User `json:"user"`
}

// length checks if string length is within specified range
func length(a, b int, e string) bool {
	return len(e) < a || len(e) > b
}

// BeforCreatePost validates post data from form request
func (post *Post) BeforCreatePost(r *http.Request) error {
	post.Content = strings.TrimSpace(r.FormValue("content"))
	post.Title = strings.TrimSpace(r.FormValue("title"))

	// Validate title and content length
	if length(3, 100, post.Title) && length(10, 2000, post.Content) {
		return errors.New("size not allowed")
	}
	return nil
}

const maxFileSize = 1000000 // 1MB file size limit

// HandleImage processes and saves uploaded images
func HandleImage(path string, file multipart.File, fileheader *multipart.FileHeader) string {
	// Check file size
	if fileheader.Size > maxFileSize {
		return ""
	}

	// Read file content
	buffer := make([]byte, fileheader.Size)
	_, err := file.Read(buffer)
	if err != nil {
		return ""
	}

	// Validate file extension
	extensions := []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg"}
	extIndex := slices.IndexFunc(extensions, func(ext string) bool {
		return strings.HasSuffix(fileheader.Filename, ext)
	})
	if extIndex == -1 {
		return ""
	}

	// Generate unique filename and save
	imageName, _ := uuid.NewV4()
	err = os.WriteFile("../frontend/assets/images/"+path+"/"+imageName.String()+extensions[extIndex], buffer, 0o644)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return imageName.String() + extensions[extIndex]
}
