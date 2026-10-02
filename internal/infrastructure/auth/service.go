package auth

import (
	"TronStream/internal/database"
	"TronStream/internal/entities"
	"TronStream/internal/infrastructure/token"
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository *database.UserRepository
	TokenService   *token.TokenService
}

func (s *AuthService) SignUp(ctx context.Context, email string, password string) (*entities.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Id:           0,
		Email:        email,
		PasswordHash: string(hash),
		Role:         entities.UserRole,
		CreatedAt:    time.Now(),
	}

	return user, nil
}

func (s *AuthService) SignIn(ctx context.Context, email string, password string) (string, error) {
	user, err := s.UserRepository.GetByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", err
	}

	token, err := s.TokenService.Generate(user.Id)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, token string) (*entities.User, error) {
	return nil, errors.New("refresh token lookup is not implemented")
}
