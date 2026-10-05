package auth

import (
	"TronStream/internal/database"
	"TronStream/internal/entities"
	"TronStream/internal/infrastructure/token"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	UserRepository   *database.UserRepository
	WallerRepository *database.WalletRepository
	TokenService     *token.Service
}

const depositAddress string = "TKSi6eG81XrSjbXEoHWzUq6Fg2ava9pDbs"

func (s *Service) SignUp(ctx context.Context, email string, password string) (*entities.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Email:          email,
		PasswordHash:   string(hash),
		DepositAddress: depositAddress,
		Role:           entities.UserRole,
		CreatedAt:      time.Now(),
	}
	userId, err := s.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Id = userId

	if s.WallerRepository == nil {
		return nil, errors.New("wallet repository is not configured")
	}
	if err := s.WallerRepository.AddWallet(ctx, userId); err != nil {
		return nil, fmt.Errorf("create wallet: %w", err)
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

	t, err := s.TokenService.Generate(user.Id)
	if err != nil {
		return "", err
	}
	return t, nil
}

func (s *Service) RefreshToken() (*entities.User, error) {
	return nil, errors.New("refresh token lookup is not implemented")
}
