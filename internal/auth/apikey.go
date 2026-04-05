package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// APIKeyPrefix is the standard prefix for all Mailngine API keys.
const APIKeyPrefix = "mn_live_"

// GenerateAPIKey creates a new API key with a cryptographically random value.
// It returns the full key (to be shown once), the prefix (for identification),
// and the SHA-256 hash (for storage and lookup).
func GenerateAPIKey() (fullKey, prefix, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", "", fmt.Errorf("generate random bytes: %w", err)
	}

	// Base64 URL encoding without padding for a URL-safe, compact representation.
	encoded := base64.RawURLEncoding.EncodeToString(b)

	fullKey = APIKeyPrefix + encoded
	prefix = fullKey[:12]
	hash = HashAPIKey(fullKey)

	return fullKey, prefix, hash, nil
}

// HashAPIKey computes the SHA-256 hex digest of an API key for secure storage.
func HashAPIKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}
