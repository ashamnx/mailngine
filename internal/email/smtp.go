package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

// SMTPSender delivers email messages to a local SMTP server (e.g., Postfix).
type SMTPSender struct {
	host string
	port string
}

// NewSMTPSender creates a new SMTPSender configured for the given host and port.
func NewSMTPSender(host, port string) *SMTPSender {
	return &SMTPSender{
		host: host,
		port: port,
	}
}

// Send delivers the composed message to the SMTP server.
// The from parameter is the envelope sender (MAIL FROM), and to contains all
// envelope recipients (RCPT TO), including CC and BCC addresses.
//
// For local relay (127.0.0.1), TLS certificate verification is skipped since
// Postfix's self-signed cert won't have a valid SAN for loopback addresses.
func (s *SMTPSender) Send(_ context.Context, from string, to []string, message []byte) error {
	addr := net.JoinHostPort(s.host, s.port)

	c, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("smtp dial %s: %w", addr, err)
	}
	defer c.Close()

	if err := c.Hello("mx.mailngine.com"); err != nil {
		return fmt.Errorf("smtp hello: %w", err)
	}

	// Upgrade to TLS if supported, skipping cert verification for local relay.
	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName:         s.host,
			InsecureSkipVerify: s.host == "127.0.0.1" || s.host == "localhost",
		}
		if err := c.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("smtp starttls: %w", err)
		}
	}

	if err := c.Mail(from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}

	for _, recipient := range to {
		if err := c.Rcpt(recipient); err != nil {
			return fmt.Errorf("smtp rcpt to %s: %w", recipient, err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}

	if _, err := w.Write(message); err != nil {
		return fmt.Errorf("smtp write message: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close data: %w", err)
	}

	return c.Quit()
}
