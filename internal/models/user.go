package models

import (
	"time"

	db "mitsimi.dev/codeShare/internal/db/sqlc"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Avatar       string    `json:"avatar"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Not exposed in JSON
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// ... existing code ...

// FromDBUser converts a db.User to models.User
func FromDBUser(user db.User) User {
	return User{
		ID:           user.ID,
		Username:     user.Username,
		Avatar:       user.Avatar.String,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
