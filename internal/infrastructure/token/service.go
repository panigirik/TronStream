package token

import (
	"TronStream/internal/database"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret              string
	authTokenRepository database.AuthTokenRepository
	JwtService          *JwtService
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *TokenService) Generate(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secret))
}

func (ts *TokenService) GenerateTokenPair(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	var oldToken, err = ts.authTokenRepository.GetByRefresh(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if oldToken.IsExpired() {
		return nil, errors.New("token is expired")
	}

	accessToken, err := ts.JwtService.GenerateAccessToken(oldToken.UserId)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := generateRefreshToken()
	{
		if err != nil {
			return nil, err
		}
	}

	if err := ts.authTokenRepository.Rotate(
		ctx,
		oldToken.UserId,
		newRefreshToken,
	); err != nil {
		return nil, err
	}

	return &TokenResponse{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil
}

func generateRefreshToken() (string, error) {
	data := make([]byte, 32)

	if _, err := rand.Read(data); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}
