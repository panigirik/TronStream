package entities

import "time"

type Deposit struct {
	Id             int64     `db:"id"`
	TxHash         string    `db:"tx_hash"`
	UserId         int64     `db:"user_id"`
	ToAddress      string    `db:"to_address"`
	FromAddress    string    `db:"from_address"`
	AmountSun      int64     `db:"amount_sun"`
	BlockNumber    int64     `db:"block_number"`
	BlockTimestamp int64     `db:"block_timestamp"`
	CreatedAt      time.Time `db:"created_at"`
}
