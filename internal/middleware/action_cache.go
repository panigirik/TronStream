package middleware

import (
	"sync"
	"time"
)

type OutboxAction struct {
	UserId    string    `json:"userId"`
	Payload   string    `json:"payload"`
	Method    string    `json:"method"`
	CreatedAt time.Time `json:"createdAt"`
}

type ActionCache struct {
	mu      sync.Mutex
	actions []OutboxAction
}

func NewActionCache() *ActionCache {
	return &ActionCache{mu: sync.Mutex{}, actions: make([]OutboxAction, 0)}
}

func (ac *ActionCache) Add(action OutboxAction) {
	ac.mu.Lock()
	ac.actions = append(ac.actions, action)
	ac.mu.Unlock()
}

func (ac *ActionCache) Flush(BatchSize int64) []OutboxAction {
	ac.mu.Lock()

	if len(ac.actions) == 0 {
		return nil
	}

	dump := ac.actions
	ac.actions = make([]OutboxAction, 0, BatchSize)
	ac.mu.Unlock()
	return dump
}
