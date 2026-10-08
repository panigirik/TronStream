package database

import (
	"TronStream/internal/entities"
	"TronStream/internal/infrastructure/background_jobs"
	"TronStream/internal/infrastructure/wallet"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository struct {
	pgx    *pgxpool.Pool
	wallet *wallet.Service
}

func NewWalletRepository(pgx *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{
		pgx:    pgx,
		wallet: wallet.NewWalletService(nil),
	}
}

func (w *WalletRepository) AddWallet(ctx context.Context, userId int64) error {
	tx, err := w.pgx.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	generatedWallet, err := w.wallet.GenerateKey()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `INSERT INTO wallets
		(user_id, address, private_key_encrypted, public_key, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		userId, generatedWallet.Address, generatedWallet.PrivateKey,
		generatedWallet.PublicKey, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (w *WalletRepository) Deposit(ctx context.Context, userId int64, amount float64) error {
	tx, err := w.pgx.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	_, err = tx.Exec(ctx, `UPDATE balances SET balance = balance + $1, updated_at = now()
		WHERE user_id = $2`, amount, userId)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// ListWatchedWallets отдаёт адреса, за которыми следит поллер депозитов.
func (w *WalletRepository) ListWatchedWallets(ctx context.Context) ([]entities.Wallet, error) {
	rows, err := w.pgx.Query(ctx, `SELECT user_id, address, created_at FROM wallets ORDER BY id`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var wallets []entities.Wallet

	for rows.Next() {
		var wallet entities.Wallet

		if err := rows.Scan(&wallet.UserId, &wallet.Address, &wallet.CreatedAt); err != nil {
			return nil, err
		}

		wallets = append(wallets, wallet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wallets, nil
}

func (w *WalletRepository) GetUsersWithPositiveBalance(ctx context.Context, user_id string) ([]background_jobs.UserDepositInfo, error) {
	rows, err := w.pgx.Query(ctx, `SELECT user_id, balance FROM Balance WHERE balance > 0`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []background_jobs.UserDepositInfo

	for rows.Next() {
		var user background_jobs.UserDepositInfo

		if err := rows.Scan(&user.UserId, &user.Balance); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (w *WalletRepository) GetWalletByUserId(ctx context.Context, user_id int64) (entities.Wallet, error) {
	query := `SELECT id, user_id, address, private_key_encrypted, public_key, created_at FROM Wallets WHERE user_id = $1`
	var wallet entities.Wallet
	err := w.pgx.QueryRow(ctx, query, user_id).Scan(
		&wallet.Id,
		&wallet.UserId,
		&wallet.Address,
		&wallet.PrivateKeyEncrypted,
		&wallet.PublicKey,
		&wallet.CreatedAt)
	if err != nil {
		return wallet, err
	}

	return wallet, nil
}

func (w *WalletRepository) ResetUserBalance(ctx context.Context, user_id int64) error {
	tx, err := w.pgx.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	_, err = w.pgx.Exec(ctx, `UPDATE WALLET SET Amount = 0 WHERE user_id = $1`, user_id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (w *WalletRepository) GetUserIdByWallet(ctx context.Context, walletAddress string) (int64, error) {
	query := `SELECT user_id FROM Wallets WHERE address = $1`
	var userId int64
	err := w.pgx.QueryRow(ctx, query, walletAddress).Scan(&userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, err
		}
		return 0, err
	}

	return userId, nil
}

// Проверка на этапе компиляции, что репозиторий закрывает интерфейс джоба.
var _ background_jobs.WatchedWalletSource = (*WalletRepository)(nil)
