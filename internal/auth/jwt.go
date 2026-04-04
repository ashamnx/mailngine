package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT claims for Hello Mail authentication.
type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"uid"`
	OrgID  uuid.UUID `json:"oid"`
	Role   string    `json:"role"`
}

// JWTManager handles JWT token generation and validation.
type JWTManager struct {
	secret []byte
	expiry time.Duration
}

// NewJWTManager creates a new JWTManager with the given secret and expiry duration.
func NewJWTManager(secret string, expiry time.Duration) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
		expiry: expiry,
	}
}

// Generate creates a new signed JWT token for the given user, organization, and role.
func (m *JWTManager) Generate(userID, orgID uuid.UUID, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "hellomail",
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
		UserID: userID,
		OrgID:  orgID,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}

	return signed, nil
}

// Validate parses and validates a JWT token string, returning the claims if valid.
func (m *JWTManager) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse jwt: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid jwt claims")
	}

	return claims, nil
}

// Expiry returns the configured token expiry duration.
func (m *JWTManager) Expiry() time.Duration {
	return m.expiry
}
