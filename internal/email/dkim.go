package email

// SignMessage signs an email message with a DKIM-Signature header using the
// provided domain, selector, and RSA private key in PEM format.
//
// TODO: Implement full DKIM signing using github.com/emersion/go-msgauth/dkim.
// The interface is defined now so that the SMTP delivery pipeline can call it
// without changes once the signing logic is implemented.
func SignMessage(message []byte, domain, selector, privateKeyPEM string) ([]byte, error) {
	// Placeholder: return the message unchanged until DKIM signing is integrated.
	return message, nil
}
