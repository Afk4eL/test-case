package user

import (
	"context"

	"github.com/google/uuid"
)

type UserUseCase interface {
	GetTokens(ctx context.Context, userId uuid.UUID) (accessToken string, refreshToken string, err error)
	Refresh(ctx context.Context, aToken string, rToken string) (accessToken string, err error)
}
