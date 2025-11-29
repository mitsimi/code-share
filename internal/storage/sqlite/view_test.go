package sqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mitsimi.dev/codeShare/internal/domain"
	"mitsimi.dev/codeShare/internal/logger"
)

func setupViewTestDB(t *testing.T) (*sql.DB, *ViewRepository, *SnippetRepository, *UserRepository) {
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

	viewRepo := NewViewRepository(storage.DB())
	snippetRepo := NewSnippetRepository(storage.DB())
	userRepo := NewUserRepository(storage.DB())
	return storage.DB(), viewRepo, snippetRepo, userRepo
}

func TestViewRepository_RecordView(t *testing.T) {
	db, viewRepo, snippetRepo, userRepo := setupViewTestDB(t)
	defer db.Close()

	// Create a user and a snippet
	user := &domain.UserCreation{ID: "user-1", Username: "user1", Email: "user1@example.com"}
	createdUser, err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	snippet := &domain.Snippet{ID: "snippet-1", Title: "Snippet 1", Content: "Content 1", Language: "go", Author: createdUser}
	err = snippetRepo.Create(context.Background(), snippet)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		viewerIdentifier := "viewer-1"
		ipAddress := "127.0.0.1"
		err := viewRepo.RecordView(context.Background(), snippet.ID, viewerIdentifier, ipAddress)
		assert.NoError(t, err)

		// Verify the view was recorded
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM snippet_views WHERE snippet_id = ? AND viewer_identifier = ?", snippet.ID, viewerIdentifier).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestViewRepository_CheckRecentView(t *testing.T) {
	db, viewRepo, snippetRepo, userRepo := setupViewTestDB(t)
	defer db.Close()

	// Create a user and a snippet
	user := &domain.UserCreation{ID: "user-1", Username: "user1", Email: "user1@example.com"}
	createdUser, err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	snippet := &domain.Snippet{ID: "snippet-1", Title: "Snippet 1", Content: "Content 1", Language: "go", Author: createdUser}
	err = snippetRepo.Create(context.Background(), snippet)
	assert.NoError(t, err)

	viewerIdentifier := "viewer-1"
	ipAddress := "127.0.0.1"

	t.Run("found", func(t *testing.T) {
		err := viewRepo.RecordView(context.Background(), snippet.ID, viewerIdentifier, ipAddress)
		assert.NoError(t, err)

		viewRecord, err := viewRepo.CheckRecentView(context.Background(), snippet.ID, viewerIdentifier)
		assert.NoError(t, err)
		assert.NotNil(t, viewRecord)
		assert.Equal(t, snippet.ID, viewRecord.SnippetID)
		assert.Equal(t, viewerIdentifier, viewRecord.ViewerIdentifier)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := viewRepo.CheckRecentView(context.Background(), snippet.ID, "non-existent-viewer")
		assert.Error(t, err)
	})
}

func TestViewRepository_IncrementViewCount(t *testing.T) {
	db, viewRepo, snippetRepo, userRepo := setupViewTestDB(t)
	defer db.Close()

	// Create a user and a snippet
	user := &domain.UserCreation{ID: "user-1", Username: "user1", Email: "user1@example.com"}
	createdUser, err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	snippet := &domain.Snippet{ID: "snippet-1", Title: "Snippet 1", Content: "Content 1", Language: "go", Author: createdUser}
	err = snippetRepo.Create(context.Background(), snippet)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		err := viewRepo.IncrementViewCount(context.Background(), snippet.ID)
		assert.NoError(t, err)

		// Verify the view count was incremented
		var views int
		err = db.QueryRow("SELECT views FROM snippets WHERE id = ?", snippet.ID).Scan(&views)
		assert.NoError(t, err)
		assert.Equal(t, 1, views)
	})
}

func TestViewRepository_CleanupOldViews(t *testing.T) {
	db, viewRepo, snippetRepo, userRepo := setupViewTestDB(t)
	defer db.Close()

	// Create a user and a snippet
	user := &domain.UserCreation{ID: "user-1", Username: "user1", Email: "user1@example.com"}
	createdUser, err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	snippet := &domain.Snippet{ID: "snippet-1", Title: "Snippet 1", Content: "Content 1", Language: "go", Author: createdUser}
	err = snippetRepo.Create(context.Background(), snippet)
	assert.NoError(t, err)

	// Record a recent view and an old view
	err = viewRepo.RecordView(context.Background(), snippet.ID, "viewer-1", "127.0.0.1")
	assert.NoError(t, err)
	_, err = db.Exec("UPDATE snippet_views SET last_viewed_at = ? WHERE viewer_identifier = ?", time.Now().Add(-31*24*time.Hour), "viewer-1")
	assert.NoError(t, err)

	err = viewRepo.RecordView(context.Background(), snippet.ID, "viewer-2", "127.0.0.2")
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		err := viewRepo.CleanupOldViews(context.Background())
		assert.NoError(t, err)

		// Verify that the old view was deleted
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM snippet_views").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}
