package email

import (
	"context"
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
// Context is accepted for future use (e.g., deadline-aware dialing) but the
// current net/smtp.SendMail implementation does not support cancellation.
func (s *SMTPSender) Send(_ context.Context, from string, to []string, message []byte) error {
	addr := net.JoinHostPort(s.host, s.port)

	// net/smtp.SendMail handles the full SMTP conversation:
	// EHLO/HELO, MAIL FROM, RCPT TO, DATA, and QUIT.
	// No authentication is needed for local Postfix relay on port 25.
	if err := smtp.SendMail(addr, nil, from, to, message); err != nil {
		return fmt.Errorf("smtp send to %s: %w", addr, err)
	}

	return nil
}
