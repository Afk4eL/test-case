package usecase

import (
	"context"
	"crypto/sha512"
	"errors"
	"test-case/internal/user"
	"test-case/internal/user/tokens"
	"test-case/internal/utils/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userPgRepo user.UserPgRepository
}

func NewUserUsecase(userPgRepo user.UserPgRepository) *userUsecase {
	return &userUsecase{userPgRepo: userPgRepo}
}

func (u *userUsecase) GetTokens(ctx context.Context, userId uuid.UUID) (accessToken string, refreshToken string, err error) {
	const op = "usecase.GetTokens"

	user, err := u.userPgRepo.FindUser(ctx, userId)
	if err != nil {
		return "", "", err
	}

	curIp := ctx.Value("user_ip").(string)

	accessToken, refreshToken, err = tokens.GenerateTokensPair(userId.String(), curIp)
	if err != nil {
		return "", "", nil
	}

	hashedToken := sha512.Sum512([]byte(refreshToken))
	hash, err := bcrypt.GenerateFromPassword(hashedToken[:], bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	user.RefreshToken = string(hash)
	user.LastUserIP = curIp
	if err := u.userPgRepo.SaveRefreshToken(ctx, *user); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *userUsecase) Refresh(ctx context.Context, aToken string, rToken string) (accessToken string, err error) {
	const op = "usecase.Refresh"

	curIp := ctx.Value("user_ip").(string)

	atClaims := &tokens.AccessTokenClaims{}
	err = tokens.ParseJWT(aToken, atClaims)
	if err != nil {
		return "", err
	}
	rtClaims := &tokens.RefreshTokenClaims{}
	err = tokens.ParseJWT(rToken, rtClaims)
	if err != nil {
		return "", err
	}

	userId, err := uuid.Parse(atClaims.UserId)
	if err != nil {
		return "", err
	}
	user, err := u.userPgRepo.FindUser(ctx, userId)
	if err != nil {
		return "", err
	}

	if atClaims.TokenId != rtClaims.TokenId ||
		atClaims.UserId != rtClaims.UserId {
		return "", errors.New("Unrelated tokens")
	}

	hashedToken := sha512.Sum512([]byte(rToken))
	if err := bcrypt.CompareHashAndPassword([]byte(user.RefreshToken), hashedToken[:]); err != nil {
		return "", err
	}

	if user.LastUserIP != curIp {
		logger.Logger.Info().Msg("IP CHANGED")
	}

	accessToken, err = tokens.GenerateAccessToken(userId.String(), curIp, rtClaims.TokenId)
	if err != nil {
		return "", nil
	}

	return accessToken, nil
}
