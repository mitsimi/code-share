package dto

type AuthResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refreshToken"`
	User         UserResponse `json:"user"`
	ExpiresAt    int64        `json:"expiresAt"`
}

type RegistrationRequest struct {
	Username string `json:"username" validate:"required,min=2,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}
