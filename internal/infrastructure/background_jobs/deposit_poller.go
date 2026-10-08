package background_jobs

import (
	"TronStream/internal/entities"
	"TronStream/internal/infrastructure/tron_node"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type WatchedWalletSource interface {
	ListWatchedWallets(ctx context.Context) ([]entities.Wallet, error)
}

type DepositLedger interface {
	CreditDeposit(ctx context.Context, deposit entities.Deposit) (credited bool, err error)
	LastDepositTimestamp(ctx context.Context, walletAddress string) (int64, error)
}

type ChainClient interface {
	GetIncomingTransfers(ctx context.Context, query tron_node.TransfersQuery) (tron_node.TransfersPage, error)
}

type DepositPollerConfig struct {
	Interval      time.Duration
	PageLimit     int
	MaxPages      int
	OnlyConfirmed bool
	StartLookback time.Duration
}

type DepositPoller struct {
	wallets WatchedWalletSource
	ledger  DepositLedger
	chain   ChainClient
	cfg     DepositPollerConfig
	cursors map[string]int64
}

func NewDepositPoller(
	wallets WatchedWalletSource,
	ledger DepositLedger,
	chain ChainClient,
	cfg DepositPollerConfig,
) *DepositPoller {
	return &DepositPoller{
		wallets: wallets,
		ledger:  ledger,
		chain:   chain,
		cfg:     cfg,
		cursors: make(map[string]int64),
	}
}

func (p *DepositPoller) Start(ctx context.Context) {
	ticker := time.NewTicker(p.cfg.Interval)
	defer ticker.Stop()

	log.Println("[Deposit Poller] Фоновый воркер успешно запущен")

	p.pollOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("[Deposit Poller] Фоновый воркер останавливается...")
			return
		case <-ticker.C:
			p.pollOnce(ctx)
		}
	}
}

func (p *DepositPoller) pollOnce(ctx context.Context) {
	wallets, err := p.wallets.ListWatchedWallets(ctx)
	if err != nil {
		log.Printf("[Deposit Poller ERROR] не удалось получить список кошельков: %v", err)
		return
	}

	for _, wallet := range wallets {
		if ctx.Err() != nil {
			return
		}

		if err := p.pollWallet(ctx, wallet); err != nil {
			log.Printf("[Deposit Poller ERROR] кошелёк %s: %v", wallet.Address, err)

			var apiErr *tron_node.APIError
			if errors.As(err, &apiErr) && apiErr.IsRateLimited() {
				log.Println("[Deposit Poller] достигнут лимит запросов, пропускаем остаток тика")
				return
			}

			continue
		}
	}
}

func (p *DepositPoller) pollWallet(ctx context.Context, wallet entities.Wallet) error {
	cursor, ok := p.cursors[wallet.Address]
	if !ok {
		restored, err := p.ledger.LastDepositTimestamp(ctx, wallet.Address)
		if err != nil {
			return fmt.Errorf("восстановление курсора: %w", err)
		}

		cursor = restored
		if cursor == 0 {
			cursor = p.initialCursor(wallet)
		}
		p.cursors[wallet.Address] = cursor
	}

	maxProcessed := cursor
	defer func() { p.cursors[wallet.Address] = maxProcessed }()

	fingerprint := ""

	for page := 0; page < p.cfg.MaxPages; page++ {
		response, err := p.chain.GetIncomingTransfers(ctx, tron_node.TransfersQuery{
			AddressBase58: wallet.Address,
			MinTimestamp:  cursor,
			Limit:         p.cfg.PageLimit,
			Fingerprint:   fingerprint,
			OnlyConfirmed: p.cfg.OnlyConfirmed,
		})
		if err != nil {
			return err
		}

		for _, transfer := range response.Transfers {
			deposit := entities.Deposit{
				TxHash:         transfer.TxID,
				UserId:         wallet.UserId,
				ToAddress:      wallet.Address,
				FromAddress:    transfer.FromHex,
				AmountSun:      transfer.AmountSun,
				BlockNumber:    transfer.BlockNumber,
				BlockTimestamp: transfer.BlockTimestamp,
			}

			credited, err := p.ledger.CreditDeposit(ctx, deposit)
			if err != nil {
				// Дальше курсор не двигаем: всё последующее имеет
				// block_timestamp >= maxProcessed и будет перезапрошено.
				return fmt.Errorf("зачисление %s: %w", transfer.TxID, err)
			}
			if credited {
				log.Printf("[Deposit Poller] зачислено %d SUN пользователю %d (tx %s)",
					transfer.AmountSun, wallet.UserId, transfer.TxID)
			}

			if transfer.BlockTimestamp > maxProcessed {
				maxProcessed = transfer.BlockTimestamp
			}
		}

		if response.Fingerprint == "" || response.RawCount == 0 {
			break
		}
		fingerprint = response.Fingerprint
	}

	return nil
}

// initialCursor выбирает стартовую точку для кошелька, по которому ещё нет
// ни одного депозита: глубже создания кошелька смотреть бессмысленно,
// а глубже StartLookback — не нужно.
func (p *DepositPoller) initialCursor(wallet entities.Wallet) int64 {
	lookback := time.Now().Add(-p.cfg.StartLookback).UnixMilli()

	if created := wallet.CreatedAt.UnixMilli(); created > lookback {
		return created
	}

	return lookback
}
