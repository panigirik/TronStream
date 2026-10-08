package entities

import "time"

type CheckTransaction struct {
	Id        int64     `db:"id"`
	UserId    int64     `db:"user_id"`
	Amount    float64   `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}
