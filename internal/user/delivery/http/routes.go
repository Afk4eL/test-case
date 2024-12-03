package user_service

import (
	"test-case/internal/user"

	"github.com/gin-gonic/gin"
)

type userService struct {
	userUC   *user.UserUseCase
	handlers *userHandlers
}

func SetupRouter(handlers *userHandlers) *gin.Engine {
	router := gin.Default()

	router.GET("/get-tokens", handlers.GetTokens)

	router.POST("/refresh", handlers.Refresh)

	return router
}
