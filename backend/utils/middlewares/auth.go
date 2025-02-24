package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"forum/db"
	"forum/models"
	"forum/utils"

	"github.com/gofrs/uuid/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("uuid")
		if err != nil {
			fmt.Println(err)
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		uuid, err := uuid.FromString(cookie.Value)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		var user models.User
		err = db.DB.QueryRow("SELECT id,first_name,last_name,nickname,image FROM users WHERE uuid = ? AND uuid_exp > ?", uuid.String(), time.Now().Unix()).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Nickname, &user.Image)
		if err != nil {
			fmt.Println(err)
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ForbidnMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("uuid")
		if err == nil {
			uuid, err := uuid.FromString(cookie.Value)
			if err == nil {
				var user models.User
				err = db.DB.QueryRow("SELECT id FROM users WHERE uuid = ? AND uuid_exp > ?", uuid.String(), time.Now().Unix()).Scan(&user.ID)
				if err == nil {
					// User is already authenticated, forbid access to login page
					utils.RespondWithError(w, http.StatusForbidden)
					return
				}
			}
		}
		// Otherwise, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
