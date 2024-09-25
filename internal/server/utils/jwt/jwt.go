package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

const (
	_secretKey = "HNG4wHwDkO5DvSqQK1vb8EetGPrfAcBuR3UwU6Nejms"
)

func GenerateJWT(payload jwt.MapClaims) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)

	token, err := claims.SignedString([]byte(_secretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseJWT(jwtToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(_secretKey), nil
	})

	return token, err
}
