package database

import (
	"TronStream/internal/entities"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PayloadRepository struct {
	db *pgxpool.Pool
}

func NewPayload(db *pgxpool.Pool) *PayloadRepository {
	return &PayloadRepository{db: db}
}

func (p *PayloadRepository) SaveBatch(ctx context.Context, logs []entities.ActionLog) error {
	if len(logs) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	query := `INSERT INTO action_logs (user_id, payload, method, created_at) VALUES ($1, $2, $3, $4)`
	for _, log := range logs {
		batch.Queue(query, log.UserId, log.Payload, log.Method, log.CreatedAt)
	}

	br := p.db.SendBatch(ctx, batch)

	defer br.Close()
	var errs []error

	for i := 0; i < len(logs); i++ {
		_, err := br.Exec()
		if err != nil {
			errs = append(errs, fmt.Errorf("log %d failed: %w", i, err))
		}
	}

	if len(errs) > 0 {
		// Возвращаем комбинированную ошибку (доступно с Go 1.20)
		return errors.Join(errs...)
	}

	return nil
}
