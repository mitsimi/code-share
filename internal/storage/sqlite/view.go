package sqlite

import (
	"context"
	"database/sql"
	"strconv"

	db "mitsimi.dev/codeShare/internal/db/sqlc"
	"mitsimi.dev/codeShare/internal/repository"
)

var _ repository.ViewRepository = (*ViewRepository)(nil)

type ViewRepository struct {
	db *sql.DB
	q  *db.Queries
}

func NewViewRepository(dbConn *sql.DB) *ViewRepository {
	return &ViewRepository{
		db: dbConn,
		q:  db.New(dbConn),
	}
}

func (r *ViewRepository) CheckRecentView(ctx context.Context, snippetID, viewerIdentifier string) (*repository.ViewRecord, error) {
	result, err := r.q.CheckRecentView(ctx, db.CheckRecentViewParams{
		SnippetID:        snippetID,
		ViewerIdentifier: viewerIdentifier,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, repository.WrapError(err, "failed to check recent view")
	}

	// Handle the interface{} type from SQLite strftime calculation
	var secondsSinceLastView int64
	switch v := result.SecondsSinceLastView.(type) {
	case int64:
		secondsSinceLastView = v
	case int:
		secondsSinceLastView = int64(v)
	case float64:
		secondsSinceLastView = int64(v)
	case string:
		// Try to parse string number
		if parsed, parseErr := strconv.ParseInt(v, 10, 64); parseErr == nil {
			secondsSinceLastView = parsed
		} else {
			secondsSinceLastView = 0 // Default to 0 if parsing fails
		}
	default:
		secondsSinceLastView = 0 // Default to 0 for unknown types
	}

	viewRecord := &repository.ViewRecord{
		SnippetID:            result.SnippetID,
		ViewerIdentifier:     result.ViewerIdentifier,
		SecondsSinceLastView: secondsSinceLastView,
	}

	// Convert sql.NullTime to time.Time
	if result.LastViewedAt.Valid {
		viewRecord.LastViewedAt = result.LastViewedAt.Time
	}

	return viewRecord, nil
}

func (r *ViewRepository) RecordView(ctx context.Context, snippetID, viewerIdentifier, ipAddress string) error {
	err := r.q.RecordView(ctx, db.RecordViewParams{
		SnippetID:        snippetID,
		ViewerIdentifier: viewerIdentifier,
		IpAddress:        sql.NullString{String: ipAddress, Valid: ipAddress != ""},
	})
	if err != nil {
		return repository.WrapError(err, "failed to record view")
	}
	return nil
}

func (r *ViewRepository) IncrementViewCount(ctx context.Context, snippetID string) error {
	err := r.q.IncrementViews(ctx, snippetID)
	if err != nil {
		return repository.WrapError(err, "failed to increment view count")
	}
	return nil
}

func (r *ViewRepository) CleanupOldViews(ctx context.Context) error {
	err := r.q.CleanupOldViews(ctx)
	if err != nil {
		return repository.WrapError(err, "failed to cleanup old views")
	}
	return nil
}
