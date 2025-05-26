package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var _ jwt.Claims = (*JWTClaims)(nil)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID string `json:"sub"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
	Type   string `json:"type"`
}

func (c JWTClaims) Valid() error {
	if c.Exp < time.Now().Unix() {
		return ErrExpiredToken
	}
	return nil
}

func (c JWTClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.Exp, 0)), nil
}

func (c JWTClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.Iat, 0)), nil
}

func (c JWTClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.Iat, 0)), nil
}

func (c JWTClaims) GetSubject() (string, error) {
	return c.UserID, nil
}

func (c JWTClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

func (c JWTClaims) GetIssuer() (string, error) {
	return "", nil
}

func (c JWTClaims) IsRefresh() bool {
	return c.Type == "refresh"
}
