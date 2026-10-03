package entities

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type BlockchainTransaction struct {
	Hash        bcrypt.HashVersionTooNewError
	Transaction Transaction
}

func NewBlockchainTransaction(hash bcrypt.HashVersionTooNewError, userId int64, amount float64, kind TransactionKind) *BlockchainTransaction {
	return &BlockchainTransaction{Hash: hash, Transaction: Transaction{UserId: userId, Amount: amount, Kind: kind, CreatedAt: time.Now()}}
}
