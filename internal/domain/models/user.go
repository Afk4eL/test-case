package models

import (
	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();column:id"`
	RefreshToken string    `json:"refresh_token" gorm:"column:refresh_token"`
	LastUserIP   string    `json:"last_user_ip" gorm:"notnull"`
}

func (User) TableName() string {
	return "users"
}
