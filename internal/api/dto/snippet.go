package dto

import (
	"time"

	"mitsimi.dev/codeShare/internal/domain"
)

// Request DTOs
type CreateSnippetRequest struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	Language string `json:"language" validate:"required"`
}

type UpdateSnippetRequest struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	Language string `json:"language" validate:"required"`
}

// Response DTOs
type SnippetResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Language  string    `json:"language"`
	AuthorID  string    `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Likes     int       `json:"likes"`
	IsLiked   bool      `json:"isLiked"`
	IsSaved   bool      `json:"isSaved"`
}

// Conversion functions
func ToSnippetResponse(snippet *domain.Snippet) SnippetResponse {
	return SnippetResponse{
		ID:        snippet.ID,
		Title:     snippet.Title,
		Content:   snippet.Content,
		Language:  snippet.Language,
		AuthorID:  snippet.AuthorID,
		CreatedAt: snippet.CreatedAt,
		UpdatedAt: snippet.UpdatedAt,
		Likes:     snippet.Likes,
		IsLiked:   snippet.IsLiked,
		IsSaved:   snippet.IsSaved,
	}
}

func ToDomainSnippet(req CreateSnippetRequest, authorID string) *domain.Snippet {
	now := time.Now()
	return &domain.Snippet{
		Title:     req.Title,
		Content:   req.Content,
		Language:  req.Language,
		AuthorID:  authorID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func UpdateDomainSnippet(snippet *domain.Snippet, req UpdateSnippetRequest) {
	snippet.Title = req.Title
	snippet.Content = req.Content
	snippet.Language = req.Language
	snippet.UpdatedAt = time.Now()
}

type ToggleActionRequest struct {
	Action string `json:"action" validate:"required,oneof=like unlike save unsave"`
}

type ToggleActionResponse struct {
	Snippet SnippetResponse `json:"snippet"`
	Action  string          `json:"action"`
}
