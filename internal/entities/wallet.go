package entities

import "time"

type Wallet struct {
	Id                  int64     `db:"id"`
	UserId              int64     `db:"user_id"`
	Address             string    `db:"address"`
	PrivateKeyEncrypted string    `db:"private_key_encrypted"`
	PublicKey           string    `db:"public_key"`
	CreatedAt           time.Time `db:"created_at"`
}
