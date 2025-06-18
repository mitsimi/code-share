package dto

import "mitsimi.dev/codeShare/internal/domain"

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
	}
}

// UpdateUserInfoRequest represents the request body for updating a user's profile
type UpdateUserInfoRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UpdatePasswordRequest represents the request body for updating a user's password
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// UpdateAvatarRequest represents the request body for updating a user's avatar
type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatarUrl"`
}
