package storage

import (
	"testing"
	"time"
)

func TestSQLiteLogin(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user first
	username := "testuserlog"
	email := "testuserlog@example.com"
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

func TestSQLiteSessionManagement(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

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

func TestSQLiteDeleteExpiredSessions(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

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

func TestSQLiteTokenRefresh(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user first
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	createdUser, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create initial session
	initialToken := "initial-token"
	initialRefreshToken := "initial-refresh-token"
	initialExpiresAt := time.Now().Add(24 * time.Hour).Unix()
	err = store.CreateSession(createdUser.ID, initialToken, initialRefreshToken, initialExpiresAt)
	if err != nil {
		t.Fatalf("Failed to create initial session: %v", err)
	}

	// Get initial session
	initialSession, err := store.GetSession(initialToken)
	if err != nil {
		t.Fatalf("Failed to get initial session: %v", err)
	}

	// Verify initial session
	if initialSession.UserID != createdUser.ID {
		t.Errorf("Expected user ID %q, got %q", createdUser.ID, initialSession.UserID)
	}
	if initialSession.RefreshToken != initialRefreshToken {
		t.Errorf("Expected refresh token %q, got %q", initialRefreshToken, initialSession.RefreshToken)
	}

	// Update session with new refresh token
	newRefreshToken := "new-refresh-token"
	newExpiresAt := time.Now().Add(48 * time.Hour).Unix()
	err = store.UpdateSessionExpiry(initialToken, newExpiresAt, newRefreshToken)
	if err != nil {
		t.Fatalf("Failed to update session: %v", err)
	}

	// Get updated session
	updatedSession, err := store.GetSession(initialToken)
	if err != nil {
		t.Fatalf("Failed to get updated session: %v", err)
	}

	// Verify updated session
	if updatedSession.UserID != createdUser.ID {
		t.Errorf("Expected user ID %q, got %q", createdUser.ID, updatedSession.UserID)
	}
	if updatedSession.RefreshToken != newRefreshToken {
		t.Errorf("Expected new refresh token %q, got %q", newRefreshToken, updatedSession.RefreshToken)
	}
	if updatedSession.ExpiresAt != newExpiresAt {
		t.Errorf("Expected new expiry time %v, got %v", newExpiresAt, updatedSession.ExpiresAt)
	}

	// Test updating non-existent session
	err = store.UpdateSessionExpiry("non-existent-token", newExpiresAt, newRefreshToken)
	if err == nil {
		t.Error("Expected error when updating non-existent session")
	}
}
