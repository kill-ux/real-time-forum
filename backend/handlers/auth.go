package handlers

import (
	"log"
	"net/http"
	"time"

	"forum/db"
	"forum/models"
	"forum/utils"
	"forum/utils/middlewares"

	"github.com/gofrs/uuid/v5"
)

// RegisterHandler handles user registration requests
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Parse request body into user struct
	err := utils.ParseBody(r, &user)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}

	// Validate and hash password
	err = user.BeforeCreate()
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid data")
		return
	}

	// Insert user into database
	err = user.CreateUser()
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]any{"message": "successfully registered."})
}

// LoginHandler handles user login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	
	// Parse login credentials
	err := utils.ParseBody(r, &user)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}

	// Retrieve user from database
	u, err := models.GetUserBy(user.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid email or password")
		return
	}

	// Verify password
	if !user.VerifyPassword(u.Password) {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid email or password")
		return
	}

	// Generate session UUID
	user = *u
	user.Password = "HashedPassword"
	uuid, err := uuid.NewV4()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}
	
	// Set session expiration to 24 hours
	user.UUID = uuid.String()
	user.UUID_EXP = time.Now().Add(time.Hour * 24).Unix()
	
	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "uuid",
		Value:  uuid.String(),
		Path:   "/",
		MaxAge: int(user.UUID_EXP),
	})
	
	// Update user session in database
	err = user.UpdateUuid()
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]any{"message": "successfully logged in.", "user": user})
}

// LogoutHandler handles user logout requests
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	
	// Clear user session from database
	userID := user.ID
	err := models.Logout(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized)
		return
	}
	
	// Delete session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "uuid",
		MaxAge: -1,
	})
	utils.RespondWithJSON(w, http.StatusNoContent)
}

// CheckAuthHandler verifies if the user is authenticated
func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Get session cookie
	cookie, err := r.Cookie("uuid")
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Validate UUID format
	uuid, err := uuid.FromString(cookie.Value)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid session ID")
		return
	}
	
	// Query user by UUID and check expiration
	var user models.User
	err = db.DB.QueryRow("SELECT * FROM users WHERE uuid = ? AND uuid_exp > ?", uuid.String(), time.Now().Unix()).Scan(utils.GetScanFields(&user)...)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	
	// Hide password in response
	user.Password = "HashedPassword"
	utils.RespondWithJSON(w, http.StatusOK, user)
}
