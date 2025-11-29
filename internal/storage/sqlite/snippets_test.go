package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"mitsimi.dev/codeShare/internal/domain"
	"mitsimi.dev/codeShare/internal/logger"
)

func setupSnippetTestDB(t *testing.T) (*sql.DB, *SnippetRepository, *UserRepository) {
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

	snippetRepo := NewSnippetRepository(storage.DB())
	userRepo := NewUserRepository(storage.DB())
	return storage.DB(), snippetRepo, userRepo
}

func TestSnippetRepository_Create(t *testing.T) {
	db, snippetRepo, userRepo := setupSnippetTestDB(t)
	defer db.Close()

	// Create a user first
	user := &domain.UserCreation{
		ID:       "test-user-id",
		Username: "test-user",
		Email:    "test@example.com",
	}
	createdUser, err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		snippet := &domain.Snippet{
			ID:       "test-snippet-id",
			Title:    "Test Snippet",
			Content:  "Hello, world!",
			Language: "go",
			Author:   createdUser,
		}

		err := snippetRepo.Create(context.Background(), snippet)
		assert.NoError(t, err)

		// Verify the snippet was created
		var title string
		err = db.QueryRow("SELECT title FROM snippets WHERE id = ?", snippet.ID).Scan(&title)
		assert.NoError(t, err)
		assert.Equal(t, snippet.Title, title)
	})
}

func TestSnippetRepository_GetByID(t *testing.T) {
	db, snippetRepo, userRepo := setupSnippetTestDB(t)
	defer db.Close()

	// Create a user first
	user := &domain.UserCreation{
		ID:       "test-user-id",
		Username: "test-user",
		Email:    "test@example.com",
	}
	createdUser, err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	// Create a snippet
	snippet := &domain.Snippet{
		ID:       "test-snippet-id",
		Title:    "Test Snippet",
		Content:  "Hello, world!",
		Language: "go",
		Author:   createdUser,
	}
	err = snippetRepo.Create(context.Background(), snippet)
	assert.NoError(t, err)

	t.Run("found", func(t *testing.T) {
		foundSnippet, err := snippetRepo.GetByID(context.Background(), snippet.ID, createdUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundSnippet)
		assert.Equal(t, snippet.ID, foundSnippet.ID)
		assert.Equal(t, snippet.Title, foundSnippet.Title)
		assert.Equal(t, snippet.Content, foundSnippet.Content)
		assert.Equal(t, snippet.Language, foundSnippet.Language)
		assert.Equal(t, createdUser.ID, foundSnippet.Author.ID)
		assert.Equal(t, createdUser.Username, foundSnippet.Author.Username)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := snippetRepo.GetByID(context.Background(), "non-existent-id", createdUser.ID)
		assert.Error(t, err)
	})
}

func TestSnippetRepository_GetAllByAuthor(t *testing.T) {
	db, snippetRepo, userRepo := setupSnippetTestDB(t)
	defer db.Close()

	// Create a user first
	user := &domain.UserCreation{
		ID:       "test-user-id",
		Username: "test-user",
		Email:    "test@example.com",
	}
	createdUser, err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	// Create a couple of snippets
	snippet1 := &domain.Snippet{
		ID:       "test-snippet-id-1",
		Title:    "Test Snippet 1",
		Content:  "Hello, world!",
		Language: "go",
		Author:   createdUser,
	}
	err = snippetRepo.Create(context.Background(), snippet1)
	assert.NoError(t, err)

	snippet2 := &domain.Snippet{
		ID:       "test-snippet-id-2",
		Title:    "Test Snippet 2",
		Content:  "Hello, again!",
		Language: "rust",
		Author:   createdUser,
	}
	err = snippetRepo.Create(context.Background(), snippet2)
	assert.NoError(t, err)

	t.Run("found", func(t *testing.T) {
		snippets, err := snippetRepo.GetAllByAuthor(context.Background(), createdUser.ID, createdUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, snippets)
		assert.Len(t, snippets, 2)
	})
}

func TestSnippetRepository_GetAll(t *testing.T) {
	db, snippetRepo, userRepo := setupSnippetTestDB(t)
	defer db.Close()

	// Create a couple of users
	user1 := &domain.UserCreation{ID: "user-1", Username: "user1", Email: "user1@example.com"}
	createdUser1, err := userRepo.Create(context.Background(), user1)
	assert.NoError(t, err)

	user2 := &domain.UserCreation{ID: "user-2", Username: "user2", Email: "user2@example.com"}
	createdUser2, err := userRepo.Create(context.Background(), user2)
	assert.NoError(t, err)

	// Create a couple of snippets
	snippet1 := &domain.Snippet{ID: "snippet-1", Title: "Snippet 1", Content: "Content 1", Language: "go", Author: createdUser1}
	err = snippetRepo.Create(context.Background(), snippet1)
	assert.NoError(t, err)

	snippet2 := &domain.Snippet{ID: "snippet-2", Title: "Snippet 2", Content: "Content 2", Language: "rust", Author: createdUser2}
	err = snippetRepo.Create(context.Background(), snippet2)
	assert.NoError(t, err)

	t.Run("found", func(t *testing.T) {
		snippets, err := snippetRepo.GetAll(context.Background(), createdUser1.ID)
		assert.NoError(t, err)
		assert.NotNil(t, snippets)
		assert.Len(t, snippets, 2)
	})
}

func TestSnippetRepository_Update(t *testing.T) {
	db, snippetRepo, userRepo := setupSnippetTestDB(t)
	defer db.Close()

	// Create a user and a snippet
	user := &domain.UserCreation{ID: "user-1", Username: "user1", Email: "user1@example.com"}
	createdUser, err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	snippet := &domain.Snippet{ID: "snippet-1", Title: "Snippet 1", Content: "Content 1", Language: "go", Author: createdUser}
	err = snippetRepo.Create(context.Background(), snippet)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		updatedSnippet := &domain.Snippet{
			ID:       snippet.ID,
			Title:    "Updated Title",
			Content:  "Updated Content",
			Language: "python",
		}
		err := snippetRepo.Update(context.Background(), updatedSnippet)
		assert.NoError(t, err)

		// Verify the update
		foundSnippet, err := snippetRepo.GetByID(context.Background(), snippet.ID, createdUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, updatedSnippet.Title, foundSnippet.Title)
		assert.Equal(t, updatedSnippet.Content, foundSnippet.Content)
		assert.Equal(t, updatedSnippet.Language, foundSnippet.Language)
	})
}

func TestSnippetRepository_Delete(t *testing.T) {
	db, snippetRepo, userRepo := setupSnippetTestDB(t)
	defer db.Close()

	// Create a user and a snippet
	user := &domain.UserCreation{ID: "user-1", Username: "user1", Email: "user1@example.com"}
	createdUser, err := userRepo.Create(context.Background(), user)
	assert.NoError(t, err)

	snippet := &domain.Snippet{ID: "snippet-1", Title: "Snippet 1", Content: "Content 1", Language: "go", Author: createdUser}
	err = snippetRepo.Create(context.Background(), snippet)
	assert.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		err := snippetRepo.Delete(context.Background(), snippet.ID)
		assert.NoError(t, err)

		// Verify it was deleted
		_, err = snippetRepo.GetByID(context.Background(), snippet.ID, createdUser.ID)
		assert.Error(t, err)
	})
}
