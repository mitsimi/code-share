package models

import "time"

// Snippet represents a code snippet
type Snippet struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Likes     int       `json:"likes"`
	IsLiked   bool      `json:"isLiked"`
}

// SnippetRequest represents the data needed to create/update a snippet
type SnippetRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}
