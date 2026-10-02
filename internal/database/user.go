package database

import (
	"context"

	"TronStream/internal/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `SELECT id, email, password_hash, role, created_at 
	          FROM users 
	          WHERE email = $1 LIMIT 1`

	var user entities.User
	err := ur.db.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
