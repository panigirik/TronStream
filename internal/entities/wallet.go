package entities

import "time"

type Wallet struct {
	Id                  int64
	UserId              int64
	Address             string
	PrivateKeyEncrypted string
	PublicKey           string
	CreatedAt           time.Time
}
