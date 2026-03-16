package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/djwhocodes/auth-service/internal/model"
	"github.com/djwhocodes/auth-service/internal/repository"
	"github.com/djwhocodes/auth-service/internal/security"
)

type AuthService struct {
	users  *repository.UserRepository
	tokens *repository.TokenRepository
	jwt    *security.JWTManager
}

func NewAuthService(
	u *repository.UserRepository,
	t *repository.TokenRepository,
	j *security.JWTManager,
) *AuthService {

	return &AuthService{
		users:  u,
		tokens: t,
		jwt:    j,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	email string,
	password string,
) error {

	hash, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	user := &model.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hash,
		IsActive:     true,
	}

	return s.users.Create(ctx, user)
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (string, string, error) {

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	ok := security.VerifyPassword(password, user.PasswordHash)
	if !ok {
		return "", "", err
	}

	accessToken, err := s.jwt.GenerateAccessToken(
		user.ID,
		[]string{"user"},
		15*time.Minute,
	)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := security.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	hash := security.HashToken(refreshToken)

	rt := &model.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	err = s.tokens.Save(ctx, rt)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(
	ctx context.Context,
	refreshToken string,
) (string, string, error) {

	hash := security.HashToken(refreshToken)

	token, err := s.tokens.FindByHash(ctx, hash)
	if err != nil {
		return "", "", err
	}

	if token.Revoked || time.Now().After(token.ExpiresAt) {
		return "", "", err
	}

	err = s.tokens.Revoke(ctx, hash)
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.jwt.GenerateAccessToken(
		token.UserID,
		[]string{"user"},
		15*time.Minute,
	)

	newRefresh, err := security.GenerateRefreshToken()

	newHash := security.HashToken(newRefresh)

	rt := &model.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    token.UserID,
		TokenHash: newHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	s.tokens.Save(ctx, rt)

	return accessToken, newRefresh, nil
}

func (s *AuthService) Logout(
	ctx context.Context,
	refreshToken string,
) error {

	hash := security.HashToken(refreshToken)

	return s.tokens.Revoke(ctx, hash)
}
