package auth

import (
	"strings"
	"testing"
)

func TestGenerateAPIKey_PrefixAndFormat(t *testing.T) {
	fullKey, prefix, hash, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("GenerateAPIKey() error = %v", err)
	}

	if !strings.HasPrefix(fullKey, APIKeyPrefix) {
		t.Errorf("fullKey = %q, want prefix %q", fullKey, APIKeyPrefix)
	}

	if !strings.HasPrefix(fullKey, prefix) {
		t.Errorf("prefix %q is not a prefix of fullKey %q", prefix, fullKey)
	}

	if len(prefix) != 12 {
		t.Errorf("prefix length = %d, want 12", len(prefix))
	}

	if hash == "" {
		t.Error("hash is empty")
	}

	// Verify that the hash is a valid hex string (64 chars for SHA-256).
	if len(hash) != 64 {
		t.Errorf("hash length = %d, want 64", len(hash))
	}
}

func TestHashAPIKey_Consistent(t *testing.T) {
	key := "hm_live_test-key-12345"

	hash1 := HashAPIKey(key)
	hash2 := HashAPIKey(key)

	if hash1 != hash2 {
		t.Errorf("HashAPIKey is not consistent: %q != %q", hash1, hash2)
	}
}

func TestHashAPIKey_DifferentKeysProduceDifferentHashes(t *testing.T) {
	hash1 := HashAPIKey("hm_live_key-aaa")
	hash2 := HashAPIKey("hm_live_key-bbb")

	if hash1 == hash2 {
		t.Error("HashAPIKey produced identical hashes for different keys")
	}
}

func TestGenerateAPIKey_Uniqueness(t *testing.T) {
	key1, _, _, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("first GenerateAPIKey() error = %v", err)
	}

	key2, _, _, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("second GenerateAPIKey() error = %v", err)
	}

	if key1 == key2 {
		t.Error("GenerateAPIKey() produced duplicate keys")
	}
}
