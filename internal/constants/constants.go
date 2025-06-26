package constants

// HTTP and Cookie constants
const (
	// SessionCookieName is the name of the session cookie
	SessionCookieName = "session"

	// CookiePath is the default path for cookies
	CookiePath = "/"
)

// Query parameter constants
const (
	// ActionQueryParam is the query parameter name for actions
	ActionQueryParam = "action"
)

// Action constants for snippet interactions
const (
	// ActionLike represents the like action
	ActionLike = "like"

	// ActionUnlike represents the unlike action
	ActionUnlike = "unlike"

	// ActionSave represents the save action
	ActionSave = "save"

	// ActionUnsave represents the unsave action
	ActionUnsave = "unsave"
)

// Default values
const (
	// DefaultLanguage is the default programming language for snippets
	DefaultLanguage = "plaintext"
)
