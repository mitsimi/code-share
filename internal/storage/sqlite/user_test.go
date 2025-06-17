package storage

import (
	"testing"

	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/storage"
)

func TestSQLiteCreateAndGetUser(t *testing.T) {
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
	user, err := store.GetUserByID(storage.UserID(userID.ID))
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

func TestSQLiteGetUserByUsername(t *testing.T) {
	store, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a user
	username := "testuser"
	email := "testuser@example.com"
	password := "password123"
	createdUser, err := store.CreateUser(username, email, password)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test getting user by username
	queriedUser, err := store.GetUserByUsername(username)
	if err != nil {
		t.Fatalf("Failed to get user by username: %v", err)
	}

	// Verify user data
	if queriedUser.ID != createdUser.ID {
		t.Errorf("Expected user ID %q, got %q", createdUser.ID, queriedUser.ID)
	}
	if queriedUser.Username != username {
		t.Errorf("Expected username %q, got %q", username, queriedUser.Username)
	}

	// Test getting non-existent user
	_, err = store.GetUserByUsername("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent user")
	}
}

func TestSQLiteUpdateUser(t *testing.T) {
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

	// Update user information
	updatedUser := db.User{
		ID:       user.ID,
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	_, err = store.UpdateUser(updatedUser.ID, updatedUser.Username, updatedUser.Email)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Get the updated user
	gotUser, err := store.GetUserByID(storage.UserID(user.ID))
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	// Verify the update
	if gotUser.Username != updatedUser.Username {
		t.Errorf("Expected username %q, got %q", updatedUser.Username, gotUser.Username)
	}
	if gotUser.Email != updatedUser.Email {
		t.Errorf("Expected email %q, got %q", updatedUser.Email, gotUser.Email)
	}
	if gotUser.ID != user.ID {
		t.Errorf("Expected user ID %q, got %q", user.ID, gotUser.ID)
	}
}
