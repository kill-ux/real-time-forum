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

type Post struct {
	ID         int64  `json:"id"`
	UserID     int    `json:"user_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Categories string `json:"categories"`
	CreatedAt  int    `json:"created_at"`
	Image      string `json:"image"`
}

type PostWithUser struct {
	Post
	GetLikes
	User User `json:"user"`
}

// length checks if a string's length is outside the specified range [a, b].
// Returns true if the string is too short or too long, false if within range.
func length(a, b int, e string) bool {
	return len(e) < a || len(e) > b
}

func (post *Post) BeforCreatePost(r *http.Request) error {
	post.Content = strings.TrimSpace(r.FormValue("content"))
	post.Title = strings.TrimSpace(r.FormValue("title"))

	if length(3, 100, post.Title) && length(10, 2000, post.Content) {
		return errors.New("size not allowed")
	}
	return nil
}

const maxFileSize = 1000000 // 1MB file size limit

// HandleImage processes and saves uploaded image files.
// It validates file size (max 1MB), checks file extension, generates a unique filename,
// and saves the file to the specified path in the frontend assets directory.
// Returns the generated filename with extension, or empty string if validation fails.
func HandleImage(path string, file multipart.File, fileheader *multipart.FileHeader) string {
	if fileheader.Size > maxFileSize {
		return ""
	}

	buffer := make([]byte, fileheader.Size)
	_, err := file.Read(buffer)
	if err != nil {
		return ""
	}

	extensions := []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg"}
	extIndex := slices.IndexFunc(extensions, func(ext string) bool {
		return strings.HasSuffix(fileheader.Filename, ext)
	})
	if extIndex == -1 {
		return ""
	}

	imageName, _ := uuid.NewV4()
	err = os.WriteFile("../frontend/assets/images/"+path+"/"+imageName.String()+extensions[extIndex], buffer, 0o644) // Safer permissions
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return imageName.String() + extensions[extIndex]
}
