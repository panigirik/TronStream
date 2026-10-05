package entities

import "time"

type Transaction struct {
	Id        int64           `db:"id"`
	UserId    int64           `db:"user_id"`
	Amount    float64         `db:"amount"`
	Kind      TransactionKind `db:"kind"`
	CreatedAt time.Time       `db:"created_at"`
}

type TransactionKind string

const (
	Deposit  TransactionKind = "deposit"
	BuyCheck TransactionKind = "buy-check"
)
