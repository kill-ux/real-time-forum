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

// RegisterHandler handles user registration requests.
// It validates user input, hashes the password, and creates a new user account.
// Returns 201 on success, 400 for invalid data, or 500 for server errors.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := utils.ParseBody(r, &user)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}

	err = user.BeforeCreate()
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid data")
		return
	}

	err = user.CreateUser()
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]any{"message": "successfully registered."})
	// Handle response
}

// LoginHandler handles user login requests.
// It verifies credentials, generates a session UUID, and sets a cookie.
// Returns 200 with user data on success, or 400 for invalid credentials.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := utils.ParseBody(r, &user)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}

	u, err := models.GetUserBy(user.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid email or password")
		return
	}

	if !user.VerifyPassword(u.Password) {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid email or password")
		return
	}

	user = *u
	user.Password = "HashedPassword"
	uuid, err := uuid.NewV4()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}
	user.UUID = uuid.String()
	user.UUID_EXP = time.Now().Add(time.Hour * 24).Unix()
	http.SetCookie(w, &http.Cookie{
		Name:   "uuid",
		Value:  uuid.String(),
		Path:   "/",
		MaxAge: int(user.UUID_EXP),
	})
	err = user.UpdateUuid()
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]any{"message": "successfully logged in.", "user": user})
	// Handle response
}

// LogoutHandler handles user logout requests.
// It invalidates the user's session and clears the authentication cookie.
// Returns 204 on success or 401 if unauthorized.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := user.ID
	err := models.Logout(userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "uuid",
		MaxAge: -1,
	})
	utils.RespondWithJSON(w, http.StatusNoContent)
}

// CheckAuthHandler verifies if a user is authenticated.
// It checks the session cookie and returns user information if valid.
// Returns 200 with user data on success, or 401 if unauthorized.
func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Get session ID from cookie
	cookie, err := r.Cookie("uuid")
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	uuid, err := uuid.FromString(cookie.Value)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid session ID")
		return
	}
	var user models.User
	err = db.DB.QueryRow("SELECT * FROM users WHERE uuid = ? AND uuid_exp > ?", uuid.String(), time.Now().Unix()).Scan(utils.GetScanFields(&user)...)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	user.Password = "HashedPassword"
	// Respond with user info
	utils.RespondWithJSON(w, http.StatusOK, user)
}
