package database

import (
	"TronStream/internal/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	pgx *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{pgx: db}
}

func (t *TransactionRepository) AddInternalTransaction(ctx context.Context, transaction entities.InternalTransaction) error {
	tx, err := t.pgx.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO internal_transactions (user_id, amount, kind, created_at)
		VALUES ($1, $2, $3, $4)`,
		transaction.Transaction.UserId,
		transaction.Transaction.Amount,
		transaction.Transaction.Kind,
		transaction.Transaction.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (t *TransactionRepository) AddBlockchainTransaction(ctx context.Context, transaction entities.BlockchainTransaction) error {
	tx, err := t.pgx.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO blockchain_transactions (user_id, amount, kind, created_at)
		VALUES ($1, $2, $3, $4)`,
		transaction.Transaction.UserId,
		transaction.Transaction.Amount,
		transaction.Transaction.Kind,
		transaction.Transaction.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
