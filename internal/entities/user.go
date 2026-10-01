package entities

import (
	_ "fmt"
)

const (
	UserRole Role = iota + 1
	AdminRole
)

type Role int8

type User struct {
	Id   int64
	Name string
	Role string
}
