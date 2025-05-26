package auth

import "testing"

// Add this test to verify token generation works correctly
func TestTokenGeneration(t *testing.T) {
	secretKey := "test-secret-key"
	userID := "test-user-id"

	// Generate two access tokens
	token1, err := GenerateToken(userID, secretKey, false)
	if err != nil {
		t.Fatalf("Failed to generate token 1: %v", err)
	}

	token2, err := GenerateToken(userID, secretKey, false)
	if err != nil {
		t.Fatalf("Failed to generate token 2: %v", err)
	}

	if token1.Token == token2.Token {
		t.Error("Generated tokens should be different")
	}

	t.Logf("Token 1: %s", token1.Token)
	t.Logf("Token 2: %s", token2.Token)
}
