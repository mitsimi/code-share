package storage

import (
	"codeShare/internal/models"
	"os"
	"testing"
	"time"
)

// setupTestDB creates a temporary database for testing
func setupTestDB(t *testing.T) (*SQLiteStorage, func()) {
	t.Helper()

	// Create a temporary database file
	dbPath := "test.db"
	store, err := NewSQLiteStorage(dbPath)
	if err != nil {
		t.Fatalf("Failed to create SQLite storage: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		store.Close()
		os.Remove(dbPath)
	}

	return store, cleanup
}

func TestCreateAndGetUser(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Test creating a user
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	userID, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test getting user by ID
	user, err := store.GetUser(userID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Verify user data
	if user.Username != username {
		t.Errorf("Expected username %q, got %q", username, user.Username)
	}
	if user.ID != userID {
		t.Errorf("Expected user ID %q, got %q", userID, user.ID)
	}
	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestGetUserByUsername(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	userID, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test getting user by username
	user, err := store.GetUserByUsername(username)
	if err != nil {
		t.Fatalf("Failed to get user by username: %v", err)
	}

	// Verify user data
	if user.ID != userID {
		t.Errorf("Expected user ID %q, got %q", userID, user.ID)
	}
	if user.Username != username {
		t.Errorf("Expected username %q, got %q", username, user.Username)
	}

	// Test getting non-existent user
	_, err = store.GetUserByUsername("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent user")
	}
}

func TestCreateAndGetSnippet(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user first
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	_, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a snippet
	snippet := models.Snippet{
		Title:   "Test Snippet",
		Content: "Test Content",
		Author:  username,
	}

	snippetID, err := store.CreateSnippet(snippet)
	if err != nil {
		t.Fatalf("Failed to create snippet: %v", err)
	}

	// Get the snippet
	gotSnippet, err := store.GetSnippet(snippetID)
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

func TestUpdateSnippet(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user first
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	_, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a snippet
	snippet := models.Snippet{
		Title:   "Original Title",
		Content: "Original Content",
		Author:  username,
	}

	snippetID, err := store.CreateSnippet(snippet)
	if err != nil {
		t.Fatalf("Failed to create snippet: %v", err)
	}

	// Update the snippet
	updatedSnippet := models.Snippet{
		ID:      snippetID,
		Title:   "Updated Title",
		Content: "Updated Content",
		Author:  username,
	}

	err = store.UpdateSnippet(updatedSnippet)
	if err != nil {
		t.Fatalf("Failed to update snippet: %v", err)
	}

	// Get the updated snippet
	gotSnippet, err := store.GetSnippet(snippetID)
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
}

// retryWithBackoff retries the given operation with exponential backoff
func retryWithBackoff(operation func() error) error {
	var err error
	for i := 0; i < 5; i++ { // Try up to 5 times
		err = operation()
		if err == nil {
			return nil
		}

		// Check if it's a busy error or transaction error
		if err.Error() == "database is locked (5) (SQLITE_BUSY)" ||
			err.Error() == "SQL logic error: cannot start a transaction within a transaction (1)" {
			// Wait with exponential backoff
			time.Sleep(time.Duration(1<<uint(i)) * 10 * time.Millisecond)
			continue
		}
		// If it's not a busy error, return immediately
		return err
	}
	return err
}

func TestDeleteSnippet(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user first
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	_, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a snippet
	snippet := models.Snippet{
		Title:   "Test Snippet",
		Content: "Test Content",
		Author:  username,
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
	_, err = store.GetSnippet(snippetID)
	if err == nil {
		t.Error("Expected error when getting deleted snippet")
	}
}
