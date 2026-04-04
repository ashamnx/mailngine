// Package domain provides domain management functionality including
// DNS record generation, DKIM key management, and domain verification.
package domain

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

const dkimKeySize = 2048

// GenerateDKIMKeyPair generates an RSA 2048-bit keypair for DKIM signing.
// It returns the private key in PEM format and the public key in base64-encoded
// DER format suitable for use in a DKIM TXT DNS record.
func GenerateDKIMKeyPair() (privateKeyPEM string, publicKeyBase64 string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, dkimKeySize)
	if err != nil {
		return "", "", fmt.Errorf("generating RSA key: %w", err)
	}

	// Encode private key to PEM
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}
	privateKeyPEM = string(pem.EncodeToMemory(privBlock))

	// Encode public key to base64 DER
	pubDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("marshaling public key: %w", err)
	}
	publicKeyBase64 = base64.StdEncoding.EncodeToString(pubDER)

	return privateKeyPEM, publicKeyBase64, nil
}

// PublicKeyFromPEM extracts the RSA public key from a PEM-encoded private key
// and returns it as a base64-encoded DER string for use in DKIM DNS records.
func PublicKeyFromPEM(privatePEM string) (string, error) {
	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		return "", errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("parsing private key: %w", err)
	}

	pubDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", fmt.Errorf("marshaling public key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(pubDER), nil
}
