package database

import (
	"TronStream/internal/entities"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CheckTransactionRepository struct {
	pgx *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *CheckTransactionRepository {
	return &CheckTransactionRepository{pgx: db}
}

func (t *CheckTransactionRepository) AddInternalTransaction(ctx context.Context, transaction entities.CheckTransaction) error {
	tx, err := t.pgx.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	_, err = tx.Exec(ctx, `INSERT INTO internal_transactions (user_id, amount, kind, created_at)
		VALUES ($1, $2, $3, $4)`,
		transaction.UserId,
		transaction.Amount,
		transaction.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
