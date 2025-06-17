package domain

import "time"

type User struct {
	ID        string
	Username  string
	Email     string
	Avatar    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserCreation struct {
	Username     string
	Email        string
	PasswordHash string
}
