package entities

import "time"

type ActionLog struct {
	UserId    string
	Payload   string
	Method    string
	CreatedAt time.Time
}
