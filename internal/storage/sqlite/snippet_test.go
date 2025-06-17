package storage

import (
	"testing"

	"mitsimi.dev/codeShare/internal/models"
)

func TestSQLiteCreateAndGetSnippet(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user first
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	user, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a snippet
	snippet := models.Snippet{
		Title:    "Test Snippet",
		Content:  "Test Content",
		Language: "go",
		Author:   user.ID,
	}

	snippetID, err := store.CreateSnippet(snippet)
	if err != nil {
		t.Fatalf("Failed to create snippet: %v", err)
	}

	// Get the snippet
	gotSnippet, err := store.GetSnippet(user.ID, snippetID)
	if err != nil {
		t.Fatalf("Failed to get snippet: %v", err)
	}

	// Verify snippet data
	if gotSnippet.Title != snippet.Title {
		t.Errorf("Expected title %q, got %q", snippet.Title, gotSnippet.Title)
	}
	if gotSnippet.Content != snippet.Content {
		t.Errorf("Expected content %q, got %q", snippet.Content, gotSnippet.Content)
	}
	if gotSnippet.Language != snippet.Language {
		t.Errorf("Expected language %q, got %q", snippet.Language, gotSnippet.Language)
	}
	if gotSnippet.Author != username {
		t.Errorf("Expected author %q, got %q", username, gotSnippet.Author)
	}
	if gotSnippet.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if gotSnippet.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
	if gotSnippet.Likes != 0 {
		t.Errorf("Expected 0 likes, got %d", gotSnippet.Likes)
	}
}

func TestSQLiteUpdateSnippet(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user first
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	user, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a snippet
	snippet := models.Snippet{
		Title:    "Original Title",
		Content:  "Original Content",
		Language: "go",
		Author:   user.ID,
	}

	snippetID, err := store.CreateSnippet(snippet)
	if err != nil {
		t.Fatalf("Failed to create snippet: %v", err)
	}

	// Update the snippet
	updatedSnippet := models.Snippet{
		ID:       snippetID,
		Title:    "Updated Title",
		Content:  "Updated Content",
		Language: "python",
		Author:   user.ID,
	}

	err = store.UpdateSnippet(updatedSnippet)
	if err != nil {
		t.Fatalf("Failed to update snippet: %v", err)
	}

	// Get the updated snippet
	gotSnippet, err := store.GetSnippet(user.ID, snippetID)
	if err != nil {
		t.Fatalf("Failed to get updated snippet: %v", err)
	}

	// Verify the update
	if gotSnippet.Title != updatedSnippet.Title {
		t.Errorf("Expected title %q, got %q", updatedSnippet.Title, gotSnippet.Title)
	}
	if gotSnippet.Content != updatedSnippet.Content {
		t.Errorf("Expected content %q, got %q", updatedSnippet.Content, gotSnippet.Content)
	}
	if gotSnippet.Language != updatedSnippet.Language {
		t.Errorf("Expected language %q, got %q", updatedSnippet.Language, gotSnippet.Language)
	}
}

func TestSQLiteDeleteSnippet(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user first
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	user, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a snippet
	snippet := models.Snippet{
		Title:   "Test Snippet",
		Content: "Test Content",
		Author:  user.ID,
	}

	snippetID, err := store.CreateSnippet(snippet)
	if err != nil {
		t.Fatalf("Failed to create snippet: %v", err)
	}

	// Delete the snippet
	err = store.DeleteSnippet(snippetID)
	if err != nil {
		t.Fatalf("Failed to delete snippet: %v", err)
	}

	// Try to get the deleted snippet
	_, err = store.GetSnippet(user.ID, snippetID)
	if err == nil {
		t.Error("Expected error when getting deleted snippet")
	}
}

func TestSQLiteToggleLikeSnippet(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create two users
	user1, err := store.CreateUser("testuser1", "testuser1@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}
	user2, err := store.CreateUser("testuser2", "testuser2@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	// Create a snippet
	snippet := models.Snippet{
		Title:    "Test Snippet",
		Content:  "Test Content",
		Language: "text",
		Author:   user1.ID,
	}

	snippetID, err := store.CreateSnippet(snippet)
	if err != nil {
		t.Fatalf("Failed to create snippet: %v", err)
	}

	// Test liking a non-existent snippet
	t.Run("LikeNonExistentSnippet", func(t *testing.T) {
		err = store.ToggleLikeSnippet(user1.ID, "non-existent-id", true)
		if err == nil {
			t.Error("Expected error when liking non-existent snippet")
		}
	})

	// Test liking the snippet
	t.Run("LikeSnippet", func(t *testing.T) {
		err = store.ToggleLikeSnippet(user1.ID, snippetID, true)
		if err != nil {
			t.Fatalf("Failed to like snippet: %v", err)
		}

		// Verify like count and IsLiked flag
		gotSnippet, err := store.GetSnippet(user1.ID, snippetID)
		if err != nil {
			t.Fatalf("Failed to get snippet after liking: %v", err)
		}
		if gotSnippet.Likes != 1 {
			t.Errorf("Expected 1 like, got %d", gotSnippet.Likes)
		}
		if !gotSnippet.IsLiked {
			t.Error("Expected IsLiked to be true after liking")
		}
	})

	// Test liking the same snippet again (should be ignored due to UNIQUE constraint)
	t.Run("LikeSnippetAgain", func(t *testing.T) {
		err = store.ToggleLikeSnippet(user1.ID, snippetID, true)
		if err != nil {
			t.Fatalf("Failed to like snippet again: %v", err)
		}
		gotSnippet, err := store.GetSnippet(user1.ID, snippetID)
		if err != nil {
			t.Fatalf("Failed to get snippet after second like: %v", err)
		}
		if gotSnippet.Likes != 1 {
			t.Errorf("Expected 1 like after second like attempt, got %d", gotSnippet.Likes)
		}
	})

	// Test another user liking the same snippet
	t.Run("LikeSnippetByAnotherUser", func(t *testing.T) {
		err = store.ToggleLikeSnippet(user2.ID, snippetID, true)
		if err != nil {
			t.Fatalf("Failed to like snippet as second user: %v", err)
		}
		gotSnippet, err := store.GetSnippet(user2.ID, snippetID)
		if err != nil {
			t.Fatalf("Failed to get snippet after second user liked: %v", err)
		}
		if gotSnippet.Likes != 2 {
			t.Errorf("Expected 2 likes after second user liked, got %d", gotSnippet.Likes)
		}
		if !gotSnippet.IsLiked {
			t.Error("Expected IsLiked to be true for second user")
		}
	})

	// Test unliking the snippet
	t.Run("UnlikeSnippet", func(t *testing.T) {
		err = store.ToggleLikeSnippet(user1.ID, snippetID, false)
		if err != nil {
			t.Fatalf("Failed to unlike snippet: %v", err)
		}
		gotSnippet, err := store.GetSnippet(user1.ID, snippetID)
		if err != nil {
			t.Fatalf("Failed to get snippet after unliking: %v", err)
		}
		if gotSnippet.Likes != 1 {
			t.Errorf("Expected 1 like after unliking, got %d", gotSnippet.Likes)
		}
		if gotSnippet.IsLiked {
			t.Error("Expected IsLiked to be false after unliking")
		}
	})

	// Test unliking a snippet that wasn't liked (should be a no-op)
	t.Run("UnlikeSnippetNotLiked", func(t *testing.T) {
		err = store.ToggleLikeSnippet(user1.ID, snippetID, false)
		if err != nil {
			t.Fatalf("Failed to unlike already unliked snippet: %v", err)
		}
		gotSnippet, err := store.GetSnippet(user1.ID, snippetID)
		if err != nil {
			t.Fatalf("Failed to get snippet after second unlike: %v", err)
		}
		if gotSnippet.Likes != 1 {
			t.Errorf("Expected 1 like after second unlike attempt, got %d", gotSnippet.Likes)
		}
	})
}
