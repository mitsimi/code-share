package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	// Add SQLite configuration options
	dbConn, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=temp_store(MEMORY)&_pragma=mmap_size(30000000000)&_pragma=page_size(4096)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	dbConn.SetMaxOpenConns(1) // SQLite only supports one writer at a time
	dbConn.SetMaxIdleConns(1)
	dbConn.SetConnMaxLifetime(time.Hour)

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
func (s *SQLiteStorage) CreateUser(username, email, password string) (db.User, error) {
	// Hash the password
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return db.User{}, err
	}

	user, err := s.q.CreateUser(s.ctx, db.CreateUserParams{
		ID:           uuid.NewString(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

// GetUser gets a user by ID
func (s *SQLiteStorage) GetUser(id UserID) (db.User, error) {
	user, err := s.q.GetUser(s.ctx, string(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, errors.New("user not found")
		}
		return db.User{}, err
	}
	return user, nil
}

// GetUserByUsername gets a user by username
func (s *SQLiteStorage) GetUserByUsername(username string) (db.User, error) {
	user, err := s.q.GetUserByUsername(s.ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, errors.New("user not found")
		}
		return db.User{}, err
	}
	return user, nil
}

// GetUserByEmail gets a user by email
func (s *SQLiteStorage) GetUserByEmail(email string) (db.User, error) {
	user, err := s.q.GetUserByEmail(s.ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, errors.New("user not found")
		}
		return db.User{}, err
	}
	return user, nil
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
func (s *SQLiteStorage) CreateSession(userID string, token string, refreshToken string, expiresAt UnixTime) error {
	_, err := s.q.CreateSession(s.ctx, db.CreateSessionParams{
		ID:           uuid.NewString(),
		UserID:       userID,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
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
		ID:           session.ID,
		UserID:       session.UserID,
		Token:        session.Token,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
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

func (s *SQLiteStorage) GetSnippets(userID UserID) ([]models.Snippet, error) {
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

func (s *SQLiteStorage) GetSnippet(userID UserID, id string) (models.Snippet, error) {
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
	// Get the user by ID to verify it exists
	_, err := s.q.GetUser(s.ctx, snippet.Author)
	if err != nil {
		return "", errors.New("author not found")
	}

	result, err := s.q.CreateSnippet(s.ctx, db.CreateSnippetParams{
		ID:      uuid.NewString(),
		Title:   snippet.Title,
		Content: snippet.Content,
		Author:  snippet.Author, // Use the user ID directly
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

func (s *SQLiteStorage) ToggleLikeSnippet(userID UserID, id SnippetID, isLike bool) error {
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

// UpdateSessionExpiry updates the expiry time and refresh token of a session
func (s *SQLiteStorage) UpdateSessionExpiry(token string, expiresAt UnixTime, refreshToken string) error {
	// First get the session to verify it exists
	session, err := s.GetSession(token)
	if err != nil {
		return err
	}

	err = s.q.UpdateSessionExpiry(s.ctx, db.UpdateSessionExpiryParams{
		Token:        session.Token,
		ExpiresAt:    expiresAt,
		RefreshToken: refreshToken,
	})
	return err
}
