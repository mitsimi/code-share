package sqlite

import (
	"context"
	"database/sql"
	"time"

	ddl "mitsimi.dev/codeShare/internal/db"
	db "mitsimi.dev/codeShare/internal/db/sqlc"
)

// Storage implements the storage interface using SQLite
type Storage struct {
	db *sql.DB
	q  *db.Queries
}

// New creates a new SQLite storage instance
func New(dbPath string) (*Storage, error) {
	// Add SQLite configuration options
	dbConn, err := sql.Open("sqlite3", dbPath+"?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=temp_store(MEMORY)&_pragma=mmap_size(30000000000)&_pragma=page_size(4096)&_pragma=busy_timeout(5000)")
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

	return &Storage{
		db: dbConn,
		q:  db.New(dbConn),
	}, nil
}

// createTables creates all necessary tables in the database
func createTables(dbConn *sql.DB) error {
	ctx := context.Background()
	_, err := dbConn.ExecContext(ctx, ddl.SchemaDDL)
	return err
}

// Close closes the database connection
func (s *Storage) Close() error {
	return s.db.Close()
}

// DB returns the underlying database connection
func (s *Storage) DB() *sql.DB {
	return s.db
}
