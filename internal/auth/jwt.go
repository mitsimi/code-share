package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

var _ jwt.Claims = (*JWTClaims)(nil)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID           string `json:"sub"`
	IsRefresh        bool   `json:"is_refresh"`
	RegisteredClaims jwt.RegisteredClaims
}

func (c JWTClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return c.RegisteredClaims.GetExpirationTime()
}

func (c JWTClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return c.RegisteredClaims.GetNotBefore()
}

func (c JWTClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return c.RegisteredClaims.GetIssuedAt()
}

func (c JWTClaims) GetIssuer() (string, error) {
	return c.RegisteredClaims.GetIssuer()
}

func (c JWTClaims) GetSubject() (string, error) {
	return c.UserID, nil
}

func (c JWTClaims) GetAudience() (jwt.ClaimStrings, error) {
	return c.RegisteredClaims.GetAudience()
}
