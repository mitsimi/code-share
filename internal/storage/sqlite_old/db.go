package storage

import (
	"context"
	"database/sql"
	"time"

	ddl "mitsimi.dev/codeShare/internal/db"
	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/storage"
)

var _ storage.StorageOLD = (*SQLiteStorage)(nil)

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

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}
