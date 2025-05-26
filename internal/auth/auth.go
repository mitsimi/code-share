package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
)

const (
	AccessTokenExpiration  = 1 * time.Hour
	RefreshTokenExpiration = 7 * 24 * time.Hour
)

type TokenResponse struct {
	Token     string
	ExpiresAt int64
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plaintext equivalent
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(userID, secretKey string, isRefreshToken bool) (TokenResponse, error) {
	var expiration time.Duration
	if isRefreshToken {
		expiration = RefreshTokenExpiration
	} else {
		expiration = AccessTokenExpiration
	}

	var tokenType string
	if isRefreshToken {
		tokenType = "refresh"
	} else {
		tokenType = "access"
	}

	now := time.Now()
	claims := JWTClaims{
		UserID: userID,
		Exp:    now.Add(expiration).Unix(),
		Iat:    now.Unix(),
		Type:   tokenType,
	}

	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		Token:     signedToken,
		ExpiresAt: claims.Exp,
	}, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string, secretKey string) (JWTClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return JWTClaims{}, ErrExpiredToken
		}
		return JWTClaims{}, ErrInvalidToken
	}

	// Check if the token is valid
	if !token.Valid {
		return JWTClaims{}, ErrInvalidToken
	}

	// Get the claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return JWTClaims{}, ErrInvalidToken
	}

	return *claims, nil
}

// GenerateRandomToken generates a random token for session management
func GenerateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
