package user_service

import (
	"context"
	"net/http"
	"test-case/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userHandlers struct {
	noteUC user.UserUseCase
}

func NewHandlers(noteUC user.UserUseCase) *userHandlers {
	return &userHandlers{noteUC: noteUC}
}

func (h *userHandlers) GetTokens(c *gin.Context) {
	const op = "user.delivery.http.GetTokens"

	userId, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Empty query"})
		return
	}

	ctx := context.WithValue(c, "user_ip", c.ClientIP())
	accessToken, refreshToken, err := h.noteUC.GetTokens(ctx, userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"AccessToken": accessToken, "RefreshToken": refreshToken})
}

func (h *userHandlers) Refresh(c *gin.Context) {
	const op = "user.delivery.http.Refresh"

	oldAccessToken := c.Query("at")
	oldRefreshToken := c.Query("rt")
	if oldAccessToken == "" || oldRefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Empty query"})
		return
	}

	ctx := context.WithValue(c, "user_ip", c.ClientIP())
	accessToken, err := h.noteUC.Refresh(ctx, oldAccessToken, oldRefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"AccessToken": accessToken})
}
