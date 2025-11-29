package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing(t *testing.T) {
	password := "my-secret-password"

	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	t.Run("correct password", func(t *testing.T) {
		assert.True(t, CheckPasswordHash(password, hashedPassword))
	})

	t.Run("incorrect password", func(t *testing.T) {
		assert.False(t, CheckPasswordHash("wrong-password", hashedPassword))
	})
}

func TestTokenGenerationAndValidation(t *testing.T) {
	secretKey := "test-secret-key"
	userID := "test-user-id"

	t.Run("valid token", func(t *testing.T) {
		tokenResponse, err := GenerateToken(userID, secretKey, false)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenResponse.Token)

		claims, err := ValidateToken(tokenResponse.Token, secretKey)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.False(t, claims.IsRefresh)
	})

	t.Run("invalid secret", func(t *testing.T) {
		tokenResponse, err := GenerateToken(userID, secretKey, false)
		assert.NoError(t, err)

		_, err = ValidateToken(tokenResponse.Token, "wrong-secret")
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("expired token", func(t *testing.T) {
		// Override expiration for testing
		originalAccessTokenExpiration := AccessTokenExpiration
		AccessTokenExpiration = 1 * time.Second
		defer func() { AccessTokenExpiration = originalAccessTokenExpiration }()

		tokenResponse, err := GenerateToken(userID, secretKey, false)
		assert.NoError(t, err)

		// Wait for the token to expire
		time.Sleep(2 * time.Second)

		_, err = ValidateToken(tokenResponse.Token, secretKey)
		assert.ErrorIs(t, err, ErrExpiredToken)
	})
}

func TestGenerateRandomToken(t *testing.T) {
	token1, err := GenerateRandomToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, token1)

	token2, err := GenerateRandomToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, token2)

	assert.NotEqual(t, token1, token2)
}
