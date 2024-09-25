package models

type User struct {
	UserID       string `gorm:"primarykey;column:guid"`
	RefreshToken string `gorm:"notnull;column:refresh_token"`
	UserIP       string `gorm:"notnull;column:user_ip"`
	ExpiresAt    int64  `gorm:"notnull;column:refresh_token_life_time"`
}

func (User) TableName() string {
	return "users"
}
