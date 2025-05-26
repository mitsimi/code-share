package storage

import (
	"context"
	"database/sql"
	"errors"

	"mitsimi.dev/codeShare/internal/auth"
	ddl "mitsimi.dev/codeShare/internal/db"
	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/models"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

var _ Storage = (*SQLiteStorage)(nil)

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
	_, err := dbConn.ExecContext(ctx, ddl.SchemaDDL)
	return err
}

// CreateUser creates a new user
func (s *SQLiteStorage) CreateUser(username, email, password string) (UserID, error) {
	// Hash the password
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return "", err
	}

	user, err := s.q.CreateUser(s.ctx, db.CreateUserParams{
		ID:           uuid.NewString(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
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
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
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
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

// GetUserByEmail gets a user by email
func (s *SQLiteStorage) GetUserByEmail(email string) (models.User, error) {
	user, err := s.q.GetUserByEmail(s.ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return models.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

// Login authenticates a user and returns their ID
func (s *SQLiteStorage) Login(email, password string) (UserID, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return "", auth.ErrInvalidCredentials
	}

	if !auth.CheckPasswordHash(password, user.PasswordHash) {
		return "", auth.ErrInvalidCredentials
	}

	return user.ID, nil
}

// CreateSession creates a new session for a user
func (s *SQLiteStorage) CreateSession(userID string, token string, expiresAt UnixTime) error {
	_, err := s.q.CreateSession(s.ctx, db.CreateSessionParams{
		ID:        uuid.NewString(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	return err
}

// GetSession gets a session by token
func (s *SQLiteStorage) GetSession(token string) (models.Session, error) {
	session, err := s.q.GetSession(s.ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Session{}, errors.New("session not found")
		}
		return models.Session{}, err
	}

	return models.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
	}, nil
}

// DeleteSession deletes a session by token
func (s *SQLiteStorage) DeleteSession(token string) error {
	return s.q.DeleteSession(s.ctx, token)
}

// DeleteExpiredSessions deletes all expired sessions
func (s *SQLiteStorage) DeleteExpiredSessions() error {
	return s.q.DeleteExpiredSessions(s.ctx)
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
			IsLiked:   snippet.IsLiked == 1,
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
		IsLiked:   snippet.IsLiked == 1,
	}, nil
}

func (s *SQLiteStorage) CreateSnippet(snippet models.Snippet) (SnippetID, error) {
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

func (s *SQLiteStorage) DeleteSnippet(id SnippetID) error {
	return s.q.DeleteSnippet(s.ctx, id)
}

func (s *SQLiteStorage) ToggleLikeSnippet(id SnippetID, isLike bool) error {
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
		if err := qtx.IncrementLikesCount(s.ctx, id); err != nil {
			return err
		}
	} else {
		if err := qtx.UnlikeSnippet(s.ctx, db.UnlikeSnippetParams{
			SnippetID: id,
			UserID:    userID,
		}); err != nil {
			return err
		}
		if err := qtx.DecrementLikesCount(s.ctx, id); err != nil {
			return err
		}
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
	for _, sampleSnippet := range sampleSnippets {
		// Check if snippet already exists using the actual author's ID
		_, err := qtx.GetSnippet(s.ctx, db.GetSnippetParams{
			UserID: userIDs[sampleSnippet.Author],
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
