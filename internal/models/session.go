package models

import "time"

// Session represents a user session
type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    int64     `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
}
