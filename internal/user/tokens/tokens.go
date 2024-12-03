package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	UserId     string
	LastUserIp string
	TokenId    string
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserId     string
	LastUserIp string
	TokenId    string
	jwt.RegisteredClaims
}

const (
	secret = "JYkX-GWnXAsuVJhqmeWxbyxOEXcx8Wm7"
)

func GenerateTokensPair(userId string, ip string) (aToken string, rToken string, err error) {
	bytes := make([]byte, 32)
	_, _ = rand.Read(bytes)
	tokenId := base64.URLEncoding.EncodeToString(bytes)[:32]

	atClaims := AccessTokenClaims{
		UserId:     userId,
		LastUserIp: ip,
		TokenId:    tokenId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	rtClaims := RefreshTokenClaims{
		UserId:     userId,
		LastUserIp: ip,
		TokenId:    tokenId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := GenerateJWT(atClaims)
	if err != nil {
		return "", "", nil
	}

	refreshToken, err := GenerateJWT(rtClaims)
	if err != nil {
		return "", "", nil
	}

	return accessToken, refreshToken, nil
}

func GenerateAccessToken(userId string, ip string, tokenId string) (aToken string, err error) {
	atClaims := AccessTokenClaims{
		UserId:     userId,
		LastUserIp: ip,
		TokenId:    tokenId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := GenerateJWT(atClaims)
	if err != nil {
		return "", nil
	}

	return accessToken, nil
}

func ParseJWT(jwtToken string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func GenerateJWT(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
