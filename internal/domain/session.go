package domain

import "time"

type UnixTime = int64

type Session struct {
	ID           string
	UserID       string
	Token        string
	RefreshToken string
	ExpiresAt    UnixTime // Unix timestamp
	CreatedAt    time.Time
}
