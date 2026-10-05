package entities

import "time"

type AuthToken struct {
	Id           int64     `db:"id"`
	UserId       int64     `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	CreatedAt    time.Time `db:"created_at"`
	ExpiresAt    time.Time `db:"expires_at"`
}

func (a *AuthToken) IsExpired() bool {
	return time.Now().After(a.ExpiresAt)
}
