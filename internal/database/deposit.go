package database

import (
	"TronStream/internal/entities"
	"TronStream/internal/infrastructure/background_jobs"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Проверка на этапе компиляции, что репозиторий закрывает интерфейс джоба.
var _ background_jobs.DepositLedger = (*DepositRepository)(nil)

type DepositRepository struct {
	pgx *pgxpool.Pool
}

func NewDepositRepository(pool *pgxpool.Pool) *DepositRepository {
	return &DepositRepository{pgx: pool}
}

func (r *DepositRepository) CreditDeposit(ctx context.Context, deposit entities.Deposit) (bool, error) {
	tx, err := r.pgx.Begin(ctx)
	if err != nil {
		return false, err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	tag, err := tx.Exec(ctx, `INSERT INTO deposits
		(tx_hash, user_id, to_address, from_address, amount_sun, block_number, block_timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (tx_hash) DO NOTHING`,
		deposit.TxHash,
		deposit.UserId,
		deposit.ToAddress,
		deposit.FromAddress,
		deposit.AmountSun,
		deposit.BlockNumber,
		deposit.BlockTimestamp)
	if err != nil {
		return false, err
	}

	// Единственный способ узнать, произошла ли вставка. Без этой проверки
	// баланс рос бы на каждом тике: из-за включительного курсора последняя
	// транзакция перезапрашивается постоянно.
	if tag.RowsAffected() == 0 {
		return false, tx.Commit(ctx)
	}

	_, err = tx.Exec(ctx, `INSERT INTO balances (user_id, balance)
		VALUES ($1, $2::numeric / 1000000)
		ON CONFLICT (user_id) DO UPDATE
			SET balance = balances.balance + EXCLUDED.balance,
				updated_at = now()`,
		deposit.UserId, deposit.AmountSun)
	if err != nil {
		return false, err
	}

	return true, tx.Commit(ctx)
}

func (r *DepositRepository) LastDepositTimestamp(ctx context.Context, walletAddress string) (int64, error) {
	query := `SELECT COALESCE(MAX(block_timestamp), 0) FROM deposits WHERE to_address = $1`

	var lastTimestamp int64
	if err := r.pgx.QueryRow(ctx, query, walletAddress).Scan(&lastTimestamp); err != nil {
		return 0, err
	}

	return lastTimestamp, nil
}
