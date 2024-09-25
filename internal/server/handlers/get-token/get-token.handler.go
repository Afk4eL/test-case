package get_token

import (
	"log"
	"net/http"
	"test-case/internal/models"
	error_struct "test-case/internal/server/error"
	"test-case/internal/server/utils/access"
	"test-case/internal/server/utils/refresh"
	"time"

	"github.com/go-chi/render"
	"gorm.io/gorm"
)

type UserRepo interface {
	FindUser(guid string) (*models.User, error)
	WriteRefreshToken(tokenHash string, addedTime time.Duration, user *models.User) error
}

type Response struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func GetToken(userRepo UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.get-token.GetToken"

		guid := r.URL.Query().Get("guid")

		if guid == "" {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, error_struct.Error{Err: "Empty query"})

			return
		}

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

		user.UserIP = r.RemoteAddr
		accessToken, err := access.GenerateAccessToken(guid, r.RemoteAddr)
		if err != nil {
			log.Printf("%s: %s", op, err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, error_struct.Error{Err: err.Error(), Msg: "Generate JWT error"})

			return
		}

		refreshToken, hash, err := refresh.GenerateRefreshToken()
		if err != nil {
			log.Printf("%s: %s", op, err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, error_struct.Error{Err: err.Error(), Msg: "Hash generate error"})

			return
		}

		userRepo.WriteRefreshToken(string(hash), time.Minute, user)

		render.JSON(w, r, Response{Access: accessToken, Refresh: refreshToken})
	}
}
