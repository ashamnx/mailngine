package email

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/emersion/go-msgauth/dkim"
)

// SignMessage signs an email message with a DKIM-Signature header using the
// provided domain, selector, and RSA private key in PEM format.
func SignMessage(message []byte, domain, selector, privateKeyPEM string) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("dkim: failed to decode PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Fall back to PKCS1 format.
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("dkim: failed to parse private key: %w", err)
		}
	}

	signer, ok := key.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("dkim: private key does not implement crypto.Signer")
	}

	opts := &dkim.SignOptions{
		Domain:   domain,
		Selector: selector,
		Signer:   signer,
		HeaderKeys: []string{
			"From", "To", "Subject", "Date", "Message-ID",
			"MIME-Version", "Content-Type", "Cc", "Reply-To",
		},
	}

	var signed bytes.Buffer
	if err := dkim.Sign(&signed, bytes.NewReader(message), opts); err != nil {
		return nil, fmt.Errorf("dkim: signing failed: %w", err)
	}

	return signed.Bytes(), nil
}
