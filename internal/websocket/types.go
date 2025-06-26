package ws

// Message types
type MessageType string

const (
	MessageTypeError          MessageType = "error"
	MessageTypeSuccess        MessageType = "success"
	MessageTypeSubscribe      MessageType = "subscribe"
	MessageTypeUnsubscribe    MessageType = "unsubscribe"
	MessageTypeUserActions    MessageType = "user_actions"
	MessageTypeSnippetUpdates MessageType = "snippet_updates"
	MessageTypeListUpdates    MessageType = "list_updates"
)

// Subscription types
type SubscriptionType string

const (
	SubTypeUserActions    SubscriptionType = "user_actions"
	SubTypeSnippetUpdates SubscriptionType = "snippet_updates"
	SubTypeListUpdates    SubscriptionType = "list_updates"
)

// WebSocket message structure
type WebSocketMessage struct {
	Type      MessageType `json:"type"`
	Data      any         `json:"data"`
	SnippetID *string     `json:"snippet_id,omitempty"`
	UserID    *string     `json:"user_id,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// Subscription request
type SubscriptionRequest struct {
	Type      SubscriptionType `json:"type"`
	SnippetID *string          `json:"snippet_id,omitempty"`
}

// User actions data - syncs liked/saved status across devices
type UserActionData struct {
	Action    string `json:"action"` // "like", "unlike", "save", "unsave"
	SnippetID string `json:"snippet_id"`
	Value     bool   `json:"value"` // true for like/save, false for unlike/unsave
}

// Snippet updates data - for single snippet view (includes content + stats)
type SnippetUpdateData struct {
	SnippetID  string `json:"snippet_id"`
	UpdateType string `json:"update_type"` // "content", "stats", "both"

	// Content changes (optional)
	Title    *string `json:"title,omitempty"`
	Content  *string `json:"content,omitempty"`
	Language *string `json:"language,omitempty"`

	// Stats changes (optional)
	ViewCount *int `json:"view_count,omitempty"`
	LikeCount *int `json:"like_count,omitempty"`
}

// List updates data - for list view (content changes only)
type ListUpdateData struct {
	SnippetID string  `json:"snippet_id"`
	Title     *string `json:"title,omitempty"`
	Content   *string `json:"content,omitempty"`
	Language  *string `json:"language,omitempty"`
}

// Broadcast message with targeting
type BroadcastMessage struct {
	Message WebSocketMessage
	Target  BroadcastTarget
}

type BroadcastTarget struct {
	Type      BroadcastTargetType
	UserID    *string
	SnippetID *string
}

type BroadcastTargetType string

const (
	BroadcastTargetTypeUser           BroadcastTargetType = "user"
	BroadcastTargetTypeSnippetUpdates BroadcastTargetType = "snippet_updates"
	BroadcastTargetTypeListUpdates    BroadcastTargetType = "list_updates"
)
