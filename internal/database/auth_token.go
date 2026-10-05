package database

import (
	"TronStream/internal/entities"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthTokenRepository struct {
	db *pgxpool.Pool
}

func NewAuthTokenRepository(db *pgxpool.Pool) *AuthTokenRepository {
	return &AuthTokenRepository{db: db}
}

func (a *AuthTokenRepository) GetByRefresh(ctx context.Context, refreshToken string) (*entities.AuthToken, error) {
	query := `SELECT Id, userId, RefreshToken FROM AuthToken WHERE RefreshToken = $1`
	var authToken entities.AuthToken
	err := a.db.QueryRow(ctx, query, refreshToken).Scan(
		&authToken.Id,
		&authToken.UserId,
		&authToken.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &authToken, nil
}

func (a *AuthTokenRepository) Rotate(ctx context.Context, userId int64, NewRefreshToken string) error {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO auth_tokens (user_id, refresh_token, created_at, expires_at)
			VALUES ($1, $2, NOW(), NOW() + INTERVAL '30 days')`, userId, NewRefreshToken)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
