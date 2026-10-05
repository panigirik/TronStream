package entities

import (
	"time"
)

type BlockchainTransaction struct {
	Hash        string      `db:"hash"`
	Transaction Transaction `db:""`
}

func NewBlockchainTransaction(hash string, userId int64, amount float64, kind TransactionKind) *BlockchainTransaction {
	return &BlockchainTransaction{Hash: hash, Transaction: Transaction{UserId: userId, Amount: amount, Kind: kind, CreatedAt: time.Now()}}
}
