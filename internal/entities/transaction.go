package entities

import "time"

type Transaction struct {
	Id        int64
	UserId    int64
	Amount    float64
	Kind      TransactionKind
	CreatedAt time.Time
}

type TransactionKind string

const (
	Deposit  TransactionKind = "deposit"
	BuyCheck TransactionKind = "buy-check"
)
