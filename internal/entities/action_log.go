package entities

import "time"

type ActionLog struct {
	UserId    string    `db:"user_id"`
	Payload   string    `db:"payload"`
	Method    string    `db:"method"`
	CreatedAt time.Time `db:"created_at"`
}
