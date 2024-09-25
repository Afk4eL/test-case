package refresh

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

const (
	tokenLength = 32
)

func GenerateRefreshToken() (string, string, error) {
	bytes := make([]byte, tokenLength)
	_, _ = rand.Read(bytes)
	refreshToken := base64.URLEncoding.EncodeToString(bytes)[:tokenLength]

	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return refreshToken, string(hash), nil
}
