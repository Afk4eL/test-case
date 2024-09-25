package refresh_token

import (
	"log"
	"net/http"
	"test-case/internal/models"
	error_struct "test-case/internal/server/error"
	"test-case/internal/server/utils/access"
	util_jwt "test-case/internal/server/utils/jwt"
	"test-case/internal/server/utils/refresh"
	user_alert "test-case/internal/server/utils/user-alert"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	tokenLength = 32
)

type UserRepo interface {
	FindUser(guid string) (*models.User, error)
	WriteRefreshToken(tokenHash string, addedTime time.Duration, user *models.User) error
}

type Response struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func RefreshToken(userRepo UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get-token.RefreshToken"

		accessToken := r.URL.Query().Get("access")
		refreshToken := r.URL.Query().Get("refresh")

		if accessToken == "" || refreshToken == "" {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, error_struct.Error{Err: "Empty query"})

			return
		}

		token, err := util_jwt.ParseJWT(accessToken)
		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors != jwt.ValidationErrorExpired {
					render.Status(r, http.StatusUnauthorized)
					render.JSON(w, r, error_struct.Error{Err: "Incorrect params", Msg: err.Error()})

					return
				}
			} else {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, error_struct.Error{Err: "Incorrect params", Msg: err.Error()})
			}
		}

		claims := token.Claims.(*jwt.MapClaims)

		guid := (*claims)["user_guid"].(string)
		ip := (*claims)["user_ip"].(string)

		user, err := userRepo.FindUser(guid)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, error_struct.Error{Err: err.Error(), Msg: "User doesn't exist"})

				return
			}
			log.Printf("%s: %s", op, err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, error_struct.Error{Err: err.Error(), Msg: "Find user failed"})

			return
		}

		if user.UserIP != ip {
			user_alert.AlertOnEmail("New IP detected")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.RefreshToken), []byte(refreshToken)); err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, error_struct.Error{Err: "Incorrect refresh token"})

			return
		}

		if user.ExpiresAt < time.Now().Unix() {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, error_struct.Error{Err: "Refresh token is expired"})

			return
		}

		newAccessToken, err := access.GenerateAccessToken(guid, r.RemoteAddr)
		if err != nil {
			render.JSON(w, r, error_struct.Error{Err: err.Error(), Msg: "Generate JWT error"})

			return
		}

		newRefreshToken, hash, err := refresh.GenerateRefreshToken()
		if err != nil {
			render.JSON(w, r, error_struct.Error{Err: err.Error(), Msg: "Hash generate error"})

			return
		}

		userRepo.WriteRefreshToken(string(hash), time.Minute, user)

		render.JSON(w, r, Response{Access: newAccessToken, Refresh: newRefreshToken})
	}
}
