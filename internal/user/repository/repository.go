package repository

import (
	"context"
	"fmt"
	"test-case/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	database *gorm.DB
}

func NewUserRepository(database *gorm.DB) *UserRepository {
	return &UserRepository{database: database}
}

func (r *UserRepository) FindUser(ctx context.Context, userId uuid.UUID) (user *models.User, err error) {
	const op = "repository.FindUser"

	var userFromDb models.User
	result := r.database.Where("id = ?", userId.String()).First(&userFromDb)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, fmt.Errorf("%s: %w", op, result.Error)
	}

	return &userFromDb, nil
}

func (r *UserRepository) SaveRefreshToken(ctx context.Context, user models.User) (err error) {
	const op = "repository.SaveRefreshToken"

	result := r.database.Save(user)
	if result.Error != nil {
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	return nil
}
