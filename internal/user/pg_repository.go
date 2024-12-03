package user

import (
	"context"
	"test-case/internal/domain/models"

	"github.com/google/uuid"
)

type UserPgRepository interface {
	FindUser(ctx context.Context, userId uuid.UUID) (user *models.User, err error)
	SaveRefreshToken(ctx context.Context, user models.User) (err error)
}
