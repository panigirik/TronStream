package entities

import "time"

const (
	UserRole Role = iota + 1
	AdminRole
)

type Role int8

type User struct {
	Id             int64     `db:"id"`
	Email          string    `db:"email"`
	PasswordHash   string    `db:"password_hash"`
	DepositAddress string    `db:"deposit_address"`
	Role           Role      `db:"role"`
	CreatedAt      time.Time `db:"created_at"`
}
