package domain

import (
	"time"
)

type Snippet struct {
	ID        string
	Title     string
	Content   string
	Language  string
	Author    *User
	CreatedAt time.Time
	UpdatedAt time.Time
	Likes     int
	IsLiked   bool
	IsSaved   bool
}
