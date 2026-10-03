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
	query := `SELECT id, email, password_hash, deposit_address, role, created_at
	          FROM users 
	          WHERE email = $1 LIMIT 1`

	var user entities.User
	err := ur.db.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.DepositAddress,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *entities.User) (int64, error) {
	var userId int64

	err := ur.db.QueryRow(ctx, `INSERT into users (email, password_hash, deposit_address, role, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`, user.Email, user.PasswordHash, user.DepositAddress, user.Role, user.CreatedAt).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
