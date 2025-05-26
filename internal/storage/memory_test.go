package storage

import (
	"testing"
	"time"

	"mitsimi.dev/codeShare/internal/models"
)

// setupTestMemory creates a new memory storage for testing
func setupTestMemory(t *testing.T) *MemoryStorage {
	t.Helper()
	return NewMemoryStorage()
}

func TestMemoryCreateAndGetUser(t *testing.T) {
	store := setupTestMemory(t)

	// Test creating a user
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	userID, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test getting user by ID
	user, err := store.GetUser(UserID(userID.ID))
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	// Verify user data
	if user.Username != username {
		t.Errorf("Expected username %q, got %q", username, user.Username)
	}
	if user.ID != userID.ID {
		t.Errorf("Expected user ID %q, got %q", userID.ID, user.ID)
	}
	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestMemoryGetUserByUsername(t *testing.T) {
	store := setupTestMemory(t)

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
	if user.ID != userID.ID {
		t.Errorf("Expected user ID %q, got %q", userID.ID, user.ID)
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

func TestMemoryCreateAndGetSnippet(t *testing.T) {
	store := setupTestMemory(t)

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

func TestMemoryUpdateSnippet(t *testing.T) {
	store := setupTestMemory(t)

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

func TestMemoryDeleteSnippet(t *testing.T) {
	store := setupTestMemory(t)

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

func TestMemoryLogin(t *testing.T) {
	store := setupTestMemory(t)

	// Create a user first
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	userID, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test successful login
	gotUserID, err := store.Login(email, password)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	if gotUserID != userID.ID {
		t.Errorf("Expected user ID %q, got %q", userID.ID, gotUserID)
	}

	// Test login with wrong password
	_, err = store.Login(email, "wrongpassword")
	if err == nil {
		t.Error("Expected error when logging in with wrong password")
	}

	// Test login with non-existent email
	_, err = store.Login("nonexistent@example.com", password)
	if err == nil {
		t.Error("Expected error when logging in with non-existent email")
	}
}

func TestMemorySessionManagement(t *testing.T) {
	store := setupTestMemory(t)

	// Create a user first
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	userID, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test creating a session
	token := "test-token"
	refreshToken := "test-refresh-token"
	expiresAt := time.Now().Add(24 * time.Hour).Unix()
	err = store.CreateSession(userID.ID, token, refreshToken, expiresAt)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Test getting the session
	session, err := store.GetSession(token)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}
	if session.UserID != userID.ID {
		t.Errorf("Expected user ID %q, got %q", userID.ID, session.UserID)
	}
	if session.Token != token {
		t.Errorf("Expected token %q, got %q", token, session.Token)
	}
	if session.ExpiresAt != expiresAt {
		t.Errorf("Expected expires at %v, got %v", expiresAt, session.ExpiresAt)
	}

	// Test getting non-existent session
	_, err = store.GetSession("nonexistent-token")
	if err == nil {
		t.Error("Expected error when getting non-existent session")
	}

	// Test deleting session
	err = store.DeleteSession(token)
	if err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}

	// Verify session is deleted
	_, err = store.GetSession(token)
	if err == nil {
		t.Error("Expected error when getting deleted session")
	}
}

func TestMemoryDeleteExpiredSessions(t *testing.T) {
	store := setupTestMemory(t)

	// Create a user first
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	userID, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create an expired session (expired 1 hour ago)
	expiredToken := "expired-token"
	expiredRefreshToken := "expired-refresh-token"
	expiredExpiresAt := time.Now().UTC().Add(-1 * time.Hour).Unix()
	err = store.CreateSession(userID.ID, expiredToken, expiredRefreshToken, expiredExpiresAt)
	if err != nil {
		t.Fatalf("Failed to create expired session: %v", err)
	}

	// Create a valid session (expires in 1 hour)
	validToken := "valid-token"
	validRefreshToken := "valid-refresh-token"
	validExpiresAt := time.Now().UTC().Add(1 * time.Hour).Unix()
	err = store.CreateSession(userID.ID, validToken, validRefreshToken, validExpiresAt)
	if err != nil {
		t.Fatalf("Failed to create valid session: %v", err)
	}

	// Delete expired sessions
	err = store.DeleteExpiredSessions()
	if err != nil {
		t.Fatalf("Failed to delete expired sessions: %v", err)
	}

	// Verify expired session is deleted
	_, err = store.GetSession(expiredToken)
	if err == nil {
		t.Error("Expected error when getting expired session")
	}

	// Verify valid session still exists
	session, err := store.GetSession(validToken)
	if err != nil {
		t.Fatalf("Failed to get valid session: %v", err)
	}
	if session.Token != validToken {
		t.Errorf("Expected token %q, got %q", validToken, session.Token)
	}
}
