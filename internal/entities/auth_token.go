package entities

import "time"

type AuthToken struct {
	Id           int64
	UserId       int64     `json:"userId"`
	RefreshToken string    `json:"refreshToken"`
	CreatedAt    time.Time `json:"createdAt"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

func (a *AuthToken) IsExpired() bool {
	return time.Now().After(a.ExpiresAt)
}
