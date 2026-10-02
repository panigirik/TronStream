package entities

import (
	_ "fmt"
	"time"
)

const (
	UserRole Role = iota + 1
	AdminRole
)

type Role int8

type User struct {
	Id           int64
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}
