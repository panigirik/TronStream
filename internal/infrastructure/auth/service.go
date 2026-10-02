package auth

import (
	"TronStream/internal/database"
	"TronStream/internal/entities"
	"TronStream/internal/infrastructure/token"
	"context"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	UserRepository *database.UserRepository
	TokenService   *token.TokenService
}

func (s *Service) SignUp(ctx context.Context, email string, password string) (*entities.User, error) {
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

func (s *Service) SignIn(ctx context.Context, email string, password string) (string, error) {
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

func (s *Service) RefreshToken(ctx context.Context, token string) (*entities.User, error) {

}
