package storage

import (
	"github.com/google/uuid"
	"mitsimi.dev/codeShare/internal/auth"
	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/storage"
)

// Seed populates the database with sample data
func (s *SQLiteStorage) Seed() error {
	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.q.WithTx(tx)

	// Create sample users if they don't exist
	users := map[string]string{
		"John Doe":                 "john@example.com",
		"Alice Johnson":            "alice@example.com",
		"Paula Cyan":               "paula@example.com",
		"Ivy Purple":               "ivy@example.com",
		"Grace Yellow":             "grace@example.com",
		"Guido van Rossum":         "guido@example.com",
		"Anders Hejlsberg":         "anders@example.com",
		"Graydon Hoare":            "graydon@example.com",
		"Solomon Hykes":            "solomon@example.com",
		"Edgar Codd":               "edgar@example.com",
		"JetBrains Team":           "jetbrains@example.com",
		"Lee Byron":                "lee@example.com",
		"Chris Lattner":            "chris@example.com",
		"José Valim":               "jose@example.com",
		"Simon Peyton Jones":       "simon@example.com",
		"Martin Odersky":           "martin@example.com",
		"Taylor Otwell":            "taylor@example.com",
		"David Heinemeier Hansson": "dhh@example.com",
	}

	// Create users and store their IDs
	userIDs := make(map[string]string)
	for username, email := range users {
		// Check if user already exists
		existingUser, err := qtx.GetUserByUsername(s.ctx, username)
		if err == nil {
			userIDs[username] = existingUser.ID
			continue
		}

		// Create new user with a default password
		passwordHash, err := auth.HashPassword("password123") // Default password for seeded users
		if err != nil {
			return err
		}

		user, err := qtx.CreateUser(s.ctx, db.CreateUserParams{
			ID:           uuid.NewString(),
			Username:     username,
			Email:        email,
			PasswordHash: passwordHash,
		})
		if err != nil {
			return err
		}
		userIDs[username] = user.ID
	}

	// Create sample snippets
	for _, sampleSnippet := range storage.SampleSnippets {
		// Check if snippet already exists using the actual author's ID
		_, err := qtx.GetSnippet(s.ctx, db.GetSnippetParams{
			UserID:    userIDs[sampleSnippet.Author],
			SnippetID: sampleSnippet.ID,
		})
		if err == nil {
			continue // Skip if snippet already exists
		}

		// Create new snippet
		_, err = qtx.CreateSnippet(s.ctx, db.CreateSnippetParams{
			ID:       sampleSnippet.ID,
			Title:    sampleSnippet.Title,
			Content:  sampleSnippet.Content,
			Language: sampleSnippet.Language,
			Author:   userIDs[sampleSnippet.Author],
		})
		if err != nil {
			return err
		}

		// Update likes count
		err = qtx.UpdateLikesCount(s.ctx, db.UpdateLikesCountParams{
			SnippetID: sampleSnippet.ID,
			ID:        sampleSnippet.ID,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
