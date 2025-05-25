package storage

import (
	ddl "codeShare/internal/db"
	db "codeShare/internal/db/sqlc"
	"codeShare/internal/models"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

// SQLiteStorage implements Storage interface with SQLite
type SQLiteStorage struct {
	db  *sql.DB
	q   *db.Queries
	ctx context.Context
}

// NewSQLiteStorage creates a new SQLite storage
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	dbConn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	if err := createTables(dbConn); err != nil {
		return nil, err
	}

	return &SQLiteStorage{
		db:  dbConn,
		q:   db.New(dbConn),
		ctx: context.Background(),
	}, nil
}

func createTables(dbConn *sql.DB) error {
	ctx := context.Background()

	_, err := dbConn.ExecContext(ctx, ddl.DDL)
	return err
}

// CreateUser creates a new user
func (s *SQLiteStorage) CreateUser(username string) (string, error) {
	user, err := s.q.CreateUser(s.ctx, db.CreateUserParams{
		ID:       uuid.NewString(),
		Username: username,
	})
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

// GetUser gets a user by ID
func (s *SQLiteStorage) GetUser(id string) (models.User, error) {
	user, err := s.q.GetUser(s.ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return models.User{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}

// GetUserByUsername gets a user by username
func (s *SQLiteStorage) GetUserByUsername(username string) (models.User, error) {
	user, err := s.q.GetUserByUsername(s.ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return models.User{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *SQLiteStorage) GetSnippets() ([]models.Snippet, error) {
	// For now, we'll use a dummy user ID since we don't have authentication
	userID := "current_user"

	snippets, err := s.q.GetSnippets(s.ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]models.Snippet, len(snippets))
	for i, snippet := range snippets {
		result[i] = models.Snippet{
			ID:        snippet.ID,
			Title:     snippet.Title,
			Content:   snippet.Content,
			Author:    snippet.AuthorUsername.String, // Handle sql.NullString
			CreatedAt: snippet.CreatedAt,
			UpdatedAt: snippet.UpdatedAt,
			Likes:     int(snippet.Likes),
			UserLikes: map[string]bool{
				userID: snippet.IsLiked == 1,
			},
		}
	}

	return result, nil
}

func (s *SQLiteStorage) GetSnippet(id string) (models.Snippet, error) {
	// For now, we'll use a dummy user ID since we don't have authentication
	userID := "current_user"

	snippet, err := s.q.GetSnippet(s.ctx, db.GetSnippetParams{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Snippet{}, errors.New("snippet not found")
		}
		return models.Snippet{}, err
	}

	return models.Snippet{
		ID:        snippet.ID,
		Title:     snippet.Title,
		Content:   snippet.Content,
		Author:    snippet.AuthorUsername.String, // Handle sql.NullString
		CreatedAt: snippet.CreatedAt,
		UpdatedAt: snippet.UpdatedAt,
		Likes:     int(snippet.Likes),
		UserLikes: map[string]bool{
			userID: snippet.IsLiked == 1,
		},
	}, nil
}

func (s *SQLiteStorage) CreateSnippet(snippet models.Snippet) (string, error) {
	// Get the user ID from the username
	user, err := s.q.GetUserByUsername(s.ctx, snippet.Author)
	if err != nil {
		return "", errors.New("author not found")
	}

	result, err := s.q.CreateSnippet(s.ctx, db.CreateSnippetParams{
		ID:      uuid.NewString(),
		Title:   snippet.Title,
		Content: snippet.Content,
		Author:  user.ID, // Use the user ID instead of username
	})
	if err != nil {
		return "", err
	}

	return result.ID, nil
}

func (s *SQLiteStorage) UpdateSnippet(snippet models.Snippet) error {
	_, err := s.q.UpdateSnippet(s.ctx, db.UpdateSnippetParams{
		Title:   snippet.Title,
		Content: snippet.Content,
		ID:      snippet.ID,
	})
	return err
}

func (s *SQLiteStorage) DeleteSnippet(id string) error {
	return s.q.DeleteSnippet(s.ctx, id)
}

func (s *SQLiteStorage) ToggleLikeSnippet(id string, isLike bool) error {
	// For now, we'll use a dummy user ID since we don't have authentication
	userID := "current_user"

	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.q.WithTx(tx)

	// Like or unlike the snippet
	if isLike {
		if err := qtx.LikeSnippet(s.ctx, db.LikeSnippetParams{
			SnippetID: id,
			UserID:    userID,
		}); err != nil {
			return err
		}
	} else {
		if err := qtx.UnlikeSnippet(s.ctx, db.UnlikeSnippetParams{
			SnippetID: id,
			UserID:    userID,
		}); err != nil {
			return err
		}
	}

	// Update the likes count
	if err := qtx.UpdateLikesCount(s.ctx, db.UpdateLikesCountParams{
		SnippetID: id,
		ID:        id,
	}); err != nil {
		return err
	}

	return tx.Commit()
}

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
	for username := range users {
		// Check if user already exists
		existingUser, err := qtx.GetUserByUsername(s.ctx, username)
		if err == nil {
			userIDs[username] = existingUser.ID
			continue
		}

		// Create new user
		user, err := qtx.CreateUser(s.ctx, db.CreateUserParams{
			ID:       uuid.NewString(),
			Username: username,
		})
		if err != nil {
			return err
		}
		userIDs[username] = user.ID
	}

	// Create sample snippets
	for _, sampleSnippet := range sampleSnippets {
		// Check if snippet already exists
		_, err := qtx.GetSnippet(s.ctx, db.GetSnippetParams{
			UserID: "current_user",
			ID:     sampleSnippet.ID,
		})
		if err == nil {
			continue // Skip if snippet already exists
		}

		// Create new snippet
		_, err = qtx.CreateSnippet(s.ctx, db.CreateSnippetParams{
			ID:      sampleSnippet.ID,
			Title:   sampleSnippet.Title,
			Content: sampleSnippet.Content,
			Author:  userIDs[sampleSnippet.Author],
		})
		if err != nil {
			return err
		}

		// Add likes if any
		for userID, liked := range sampleSnippet.UserLikes {
			if liked {
				err = qtx.LikeSnippet(s.ctx, db.LikeSnippetParams{
					SnippetID: sampleSnippet.ID,
					UserID:    userID,
				})
				if err != nil {
					return err
				}
			}
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

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
