package models

import "time"

// Snippet represents a code snippet
type Snippet struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Language  string    `json:"language"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Likes     int       `json:"likes"`
	IsLiked   bool      `json:"isLiked"`
	IsSaved   bool      `json:"isSaved"`
}

// SnippetCreateRequest represents the data needed to create/update a snippet
type SnippetCreateRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Language string `json:"language"`
}
