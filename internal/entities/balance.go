package entities

import "time"

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

func (bal *Balance) spend(amount float64) {
	bal.Balance = bal.Balance - amount
}

func (bal *Balance) deposit(amount float64) {
	bal.Balance = bal.Balance + amount
}
