package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"mitsimi.dev/codeShare/internal/domain"
	"mitsimi.dev/codeShare/internal/logger"
)

func setupLikeBookmarkTestDB(t *testing.T) (*sql.DB, *LikeRepository, *BookmarkRepository, *SnippetRepository, *UserRepository) {
	err := logger.Init(logger.Config{
		Environment: "development",
		Level:       "debug",
	})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	storage, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	likeRepo := NewLikeRepository(storage.DB())
	bookmarkRepo := NewBookmarkRepository(storage.DB())
	snippetRepo := NewSnippetRepository(storage.DB())
	userRepo := NewUserRepository(storage.DB())
	return storage.DB(), likeRepo, bookmarkRepo, snippetRepo, userRepo
}

func TestLikeRepository_ToggleLike(t *testing.T) {
	db, likeRepo, _, snippetRepo, userRepo := setupLikeBookmarkTestDB(t)
	defer db.Close()

	// Create a user and a snippet
	user := &domain.UserCreation{ID: "user-1", Username: "user1", Email: "user1@example.com"}
	createdUser, err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	snippet := &domain.Snippet{ID: "snippet-1", Title: "Snippet 1", Content: "Content 1", Language: "go", Author: createdUser}
	err = snippetRepo.Create(context.Background(), snippet)
	assert.NoError(t, err)

	t.Run("like", func(t *testing.T) {
		err := likeRepo.ToggleLike(context.Background(), createdUser.ID, snippet.ID, true)
		assert.NoError(t, err)

		// Verify the like was created
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM user_likes WHERE user_id = ? AND snippet_id = ?", createdUser.ID, snippet.ID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("unlike", func(t *testing.T) {
		err := likeRepo.ToggleLike(context.Background(), createdUser.ID, snippet.ID, false)
		assert.NoError(t, err)

		// Verify the like was removed
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM user_likes WHERE user_id = ? AND snippet_id = ?", createdUser.ID, snippet.ID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

func TestBookmarkRepository_ToggleSave(t *testing.T) {
	db, _, bookmarkRepo, snippetRepo, userRepo := setupLikeBookmarkTestDB(t)
	defer db.Close()

	// Create a user and a snippet
	user := &domain.UserCreation{ID: "user-1", Username: "user1", Email: "user1@example.com"}
	createdUser, err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	snippet := &domain.Snippet{ID: "snippet-1", Title: "Snippet 1", Content: "Content 1", Language: "go", Author: createdUser}
	err = snippetRepo.Create(context.Background(), snippet)
	assert.NoError(t, err)

	t.Run("save", func(t *testing.T) {
		err := bookmarkRepo.ToggleSave(context.Background(), createdUser.ID, snippet.ID, true)
		assert.NoError(t, err)

		// Verify the save was created
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM user_saves WHERE user_id = ? AND snippet_id = ?", createdUser.ID, snippet.ID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("unsave", func(t *testing.T) {
		err := bookmarkRepo.ToggleSave(context.Background(), createdUser.ID, snippet.ID, false)
		assert.NoError(t, err)

		// Verify the save was removed
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM user_saves WHERE user_id = ? AND snippet_id = ?", createdUser.ID, snippet.ID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}
