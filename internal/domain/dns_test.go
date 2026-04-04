package domain

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestGenerateDKIMKeyPair_NonEmpty(t *testing.T) {
	privPEM, pubBase64, err := GenerateDKIMKeyPair()
	if err != nil {
		t.Fatalf("GenerateDKIMKeyPair() error = %v", err)
	}

	if privPEM == "" {
		t.Error("private key PEM is empty")
	}
	if pubBase64 == "" {
		t.Error("public key base64 is empty")
	}
}

func TestGenerateDKIMKeyPair_ValidPEM(t *testing.T) {
	privPEM, _, err := GenerateDKIMKeyPair()
	if err != nil {
		t.Fatalf("GenerateDKIMKeyPair() error = %v", err)
	}

	if !strings.Contains(privPEM, "-----BEGIN RSA PRIVATE KEY-----") {
		t.Error("private key PEM does not contain expected header")
	}
	if !strings.Contains(privPEM, "-----END RSA PRIVATE KEY-----") {
		t.Error("private key PEM does not contain expected footer")
	}
}

func TestGenerateDKIMKeyPair_ValidBase64PublicKey(t *testing.T) {
	_, pubBase64, err := GenerateDKIMKeyPair()
	if err != nil {
		t.Fatalf("GenerateDKIMKeyPair() error = %v", err)
	}

	// Verify the public key is valid base64.
	_, err = base64.StdEncoding.DecodeString(pubBase64)
	if err != nil {
		t.Errorf("public key is not valid base64: %v", err)
	}
}

func TestPublicKeyFromPEM_RoundTrip(t *testing.T) {
	privPEM, originalPub, err := GenerateDKIMKeyPair()
	if err != nil {
		t.Fatalf("GenerateDKIMKeyPair() error = %v", err)
	}

	extractedPub, err := PublicKeyFromPEM(privPEM)
	if err != nil {
		t.Fatalf("PublicKeyFromPEM() error = %v", err)
	}

	if extractedPub != originalPub {
		t.Errorf("PublicKeyFromPEM() = %q, want %q", extractedPub, originalPub)
	}
}

func TestPublicKeyFromPEM_InvalidPEM(t *testing.T) {
	_, err := PublicKeyFromPEM("not-a-valid-pem")
	if err == nil {
		t.Fatal("PublicKeyFromPEM() expected error for invalid PEM, got nil")
	}
}

func TestGenerateDKIMKeyPair_UniqueKeys(t *testing.T) {
	priv1, _, err := GenerateDKIMKeyPair()
	if err != nil {
		t.Fatalf("first GenerateDKIMKeyPair() error = %v", err)
	}

	priv2, _, err := GenerateDKIMKeyPair()
	if err != nil {
		t.Fatalf("second GenerateDKIMKeyPair() error = %v", err)
	}

	if priv1 == priv2 {
		t.Error("GenerateDKIMKeyPair() produced identical key pairs")
	}
}
