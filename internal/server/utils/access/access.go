package access

import (
	"time"

	util_jwt "test-case/internal/server/utils/jwt"

	"github.com/dgrijalva/jwt-go"
)

func GenerateAccessToken(guid string, ip string) (string, error) {
	payload := jwt.MapClaims{
		"user_guid": guid,
		"user_ip":   ip,
		"exp":       time.Now().Add(time.Second).Unix(),
	}
	accessToken, err := util_jwt.GenerateJWT(payload)
	if err != nil {
		return "", nil
	}

	return accessToken, nil
}
