package repository

// Container holds all repository instances for dependency injection
type Container struct {
	Snippets  SnippetRepository
	Likes     LikeRepository
	Bookmarks BookmarkRepository
	Users     UserRepository
	Sessions  SessionRepository
	Views     ViewRepository
}

// NewContainer creates a new repository container with all repositories
func NewContainer(
	snippets SnippetRepository,
	likes LikeRepository,
	bookmarks BookmarkRepository,
	users UserRepository,
	sessions SessionRepository,
	views ViewRepository,
) *Container {
	return &Container{
		Snippets:  snippets,
		Likes:     likes,
		Bookmarks: bookmarks,
		Users:     users,
		Sessions:  sessions,
		Views:     views,
	}
}
