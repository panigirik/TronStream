package database

import (
	"TronStream/internal/infrastructure/wallet"
	"context"
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
