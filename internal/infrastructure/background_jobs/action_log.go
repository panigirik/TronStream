package background_jobs

import (
	"TronStream/internal/entities"
	"TronStream/internal/middleware"
	"context"
	"log"
	"time"
)

type OutboxSaver interface {
	SaveBatch(ctx context.Context, actions []entities.ActionLog) error
}

type OutboxJob struct {
	cache    *middleware.ActionCache
	repo     OutboxSaver
	interval time.Duration
}

func NewOutboxJob(cache *middleware.ActionCache, repo OutboxSaver, interval time.Duration) *OutboxJob {
	return &OutboxJob{
		cache:    cache,
		repo:     repo,
		interval: interval,
	}
}

func (j *OutboxJob) Start(ctx context.Context) {
	// Тикер будет пинговать нас каждые X секунд
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	log.Println("[Outbox Job] Фоновый воркер успешно запущен")

	for {
		select {
		case <-ctx.Done():
			// Если приложение завершает работу (Graceful Shutdown), выходим из цикла
			log.Println("[Outbox Job] Фоновый воркер останавливается...")
			return
		case <-ticker.C:
			// Проснулись по тикеру — вытряхиваем кэш
			// Задаем размер батча, например 500
			rawActions := j.cache.Flush(500)
			if len(rawActions) == 0 {
				continue // Кэш пустой, спать дальше
			}

			dbActions := make([]entities.ActionLog, len(rawActions))
			for i, a := range rawActions {
				dbActions[i] = entities.ActionLog{
					UserId:    a.UserId,
					Payload:   a.Payload,
					Method:    a.Method,
					CreatedAt: a.CreatedAt,
				}
			}

			// Пытаемся сохранить пачку в базу данных
			err := j.repo.SaveBatch(ctx, dbActions)
			if err != nil {
				log.Printf("[Outbox Job ERROR] Не удалось сохранить батч в БД: %v", err)
				// ПОДУМАЙ: Если база упала, данные из кэша мы уже стерли вызовом Flush().
				// В данном случае они потеряются. Если это критично, то нужно возвращать их в кэш.
			} else {
				log.Printf("[Outbox Job] Успешно сохранено %d записей в аутбокс", len(dbActions))
			}
		}
	}
}
