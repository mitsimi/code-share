package domain

import (
	"time"

	db "mitsimi.dev/codeShare/internal/db/sqlc"
)

type User struct {
	ID           string
	Username     string
	Email        string
	Avatar       *string
	PasswordHash string // This should be kept private and not exposed in the User struct
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserCreation struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
}

func ToDomainUser(user db.User) *User {
	var avatar *string
	if user.Avatar.Valid {
		avatar = &user.Avatar.String
	}

	return &User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Avatar:       avatar,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
