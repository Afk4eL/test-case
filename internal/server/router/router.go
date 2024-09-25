package router

import (
	get_token "test-case/internal/server/handlers/get-token"
	refresh_token "test-case/internal/server/handlers/refresh-token"
	"test-case/storage/repos"

	"github.com/go-chi/chi/v5"
)

func NewRouter(userRepo repos.UserRepository) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/get-token", get_token.GetToken(userRepo))
	router.Post("/refresh-token", refresh_token.RefreshToken(userRepo))

	return router
}
