package entities

import "time"

type Balance struct {
	Id        int64     `db:"id"`
	UserId    int64     `db:"user_id"`
	Balance   float64   `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
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
