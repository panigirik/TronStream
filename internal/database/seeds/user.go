package seeds

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserSeed struct {
	Db *pgxpool.Pool
}

func (u *UserSeed) SeedAdminUser(ctx context.Context) error {
	query := `INSERT INTO Users (id, email, password_hash, deposit_address, role, created_at)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (id) DO NOTHING`

	id := 1
	email := "zaharvaskovskij78@gmail.com"
	passwordHash := "7a3c3e34b12c8b81369b76c813a82f3a473187a552880c85c2b0e6378e91b5c4"
	depositAddress := "TUUxkyGfb2bVrjjmnkhiFNRSTnffEn14MN"
	role := "Admin"
	createdAt, err := time.Parse(time.RFC3339, "2020-07-20T00:00:00Z")
	if err != nil {
		return err
	}

	_, err = u.Db.Exec(ctx, query, id, email, passwordHash, depositAddress, role, createdAt)
	if err != nil {
		return err
	}

	return nil
}
