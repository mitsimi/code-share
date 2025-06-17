package auth

// TODO: Implement role-based access control (RBAC) for users
// IsAdmin checks if the user has admin privileges
func IsAdmin(userID string) bool {
	// This function should check if the user is an admin.
	// For now, we will assume that user ID "admin" is the only admin.
	return userID == "admin"
}
