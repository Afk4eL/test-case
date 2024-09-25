package repos

import (
	"fmt"
	"test-case/internal/models"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUser(guid string) (*models.User, error)
	WriteRefreshToken(tokenHash string, addedTime time.Duration, user *models.User) error
}

type userRepo struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{database: db}
}

func (r *userRepo) FindUser(guid string) (*models.User, error) {
	const op = "storage.repos.FindUser"

	var userFromDb models.User
	result := r.database.Where("guid = ?", guid).First(&userFromDb)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, fmt.Errorf("%s: %w", op, result.Error)
	}

	return &userFromDb, nil
}

func (r *userRepo) WriteRefreshToken(tokenHash string, addedTime time.Duration, user *models.User) error {
	const op = "storage.repos.WriteRefreshToken"

	user.RefreshToken = tokenHash
	user.ExpiresAt = time.Now().Add(addedTime).Unix()

	result := r.database.Save(user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return result.Error
		}
		return fmt.Errorf("%s: %w", op, result.Error)
	}

	return nil
}
