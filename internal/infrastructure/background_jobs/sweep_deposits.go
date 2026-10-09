package background_jobs

import (
	"TronStream/internal/entities"
	"context"
	"log"
)

type UserDepositInfo struct {
	UserId  int64
	Balance float64
}

const (
	MainAddress = "xx"
)

type WalletDetails struct {
	Address             string
	PrivateKeyEncrypted string
}

type WalletRepositroyInterface interface {
	GetUsersWithPositiveBalance(ctx context.Context) ([]UserDepositInfo, error)
	GetWalletByUserId(ctx context.Context, user_id int64) (entities.Wallet, error)
	TransferCrypto(ctx context.Context, user_id int64, amount float64, privateKey string) error
	ResetUserBalance(ctx context.Context, user_id int64) error
}

type SweepDeposits struct {
	repo WalletRepositroyInterface
}

func (s *SweepDeposits) Run(ctx context.Context) error {
	users, err := s.repo.GetUsersWithPositiveBalance(ctx)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	for _, user := range users {
		wallet, err := s.repo.GetWalletByUserId(ctx, user.UserId)
		if err != nil {
			log.Printf("[ERROR] failed to get wallet for user %d: %v", user.UserId, err)
			continue
		}

		err = s.repo.TransferCrypto(ctx, user.UserId, user.Balance, wallet.PrivateKeyEncrypted)
		if err != nil {
			log.Printf("[ERROR] failed to transfer crypto for user %d: %v", user.UserId, err)
			continue
		}

		err = s.repo.ResetUserBalance(ctx, user.UserId)
		if err != nil {
			log.Printf("[CRITICAL ERROR] crypto swept but failed to reset DB balance for user %d: %v", user.UserId, err)
		}
	}

	return nil
}
