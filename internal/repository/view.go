package repository

import (
	"context"
	"time"
)

// ViewRecord represents a view tracking record
type ViewRecord struct {
	SnippetID            string
	ViewerIdentifier     string
	IPAddress            string
	LastViewedAt         time.Time
	SecondsSinceLastView int64
}

// ViewRepository defines the interface for view tracking operations
type ViewRepository interface {
	// CheckRecentView returns the last view record for a snippet by a viewer
	CheckRecentView(ctx context.Context, snippetID, viewerIdentifier string) (*ViewRecord, error)

	// RecordView records a view attempt (always executes)
	RecordView(ctx context.Context, snippetID, viewerIdentifier, ipAddress string) error

	// IncrementViewCount increments the public view counter for a snippet
	IncrementViewCount(ctx context.Context, snippetID string) error

	// CleanupOldViews removes old view tracking records
	CleanupOldViews(ctx context.Context) error
}
