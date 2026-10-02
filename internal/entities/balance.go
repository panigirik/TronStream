package entities

import (
	_ "fmt"
	"time"
)

type Balance struct {
	Id        int64
	UserId    int64
	Balance   float64
	CreatedAt time.Time
}

func NewBalance(userId int64) *Balance {
	return &Balance{
		UserId:    userId,
		Balance:   0.0,
		CreatedAt: time.Now(),
	}
}
