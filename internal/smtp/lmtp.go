package smtp

import (
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/mail"
	"os"
	"strings"
	"time"

	gosmtp "github.com/emersion/go-smtp"
	"github.com/rs/zerolog"

	sqlcdb "github.com/mailngine/mailngine/internal/db/sqlcdb"
	"github.com/mailngine/mailngine/internal/inbox"
)

// LMTPServer receives inbound emails from Postfix and delivers them to the inbox.
type LMTPServer struct {
	server     *gosmtp.Server
	socketPath string
	queries    *sqlcdb.Queries
	inboxSvc   *inbox.Service
	logger     zerolog.Logger
}

// NewLMTPServer creates a new LMTP server that listens on the given Unix socket.
func NewLMTPServer(socketPath string, queries *sqlcdb.Queries, inboxSvc *inbox.Service, logger zerolog.Logger) *LMTPServer {
	l := &LMTPServer{
		socketPath: socketPath,
		queries:    queries,
		inboxSvc:   inboxSvc,
		logger:     logger.With().Str("component", "lmtp").Logger(),
	}

	s := gosmtp.NewServer(l)
	s.LMTP = true
	s.Domain = "mx.mailngine.com"
	s.ReadTimeout = 60 * time.Second
	s.WriteTimeout = 60 * time.Second
	s.MaxMessageBytes = 25 * 1024 * 1024 // 25MB
	s.MaxRecipients = 100

	l.server = s
	return l
}

// NewSession implements gosmtp.Backend.
func (l *LMTPServer) NewSession(_ *gosmtp.Conn) (gosmtp.Session, error) {
	return &lmtpSession{
		queries:  l.queries,
		inboxSvc: l.inboxSvc,
		logger:   l.logger,
	}, nil
}

// ListenAndServe starts the LMTP server on the Unix socket.
func (l *LMTPServer) ListenAndServe() error {
	// Remove stale socket file if it exists.
	os.Remove(l.socketPath)

	listener, err := net.Listen("unix", l.socketPath)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", l.socketPath, err)
	}

	// Postfix needs to be able to connect to the socket.
	if err := os.Chmod(l.socketPath, 0777); err != nil {
		listener.Close()
		return fmt.Errorf("chmod socket: %w", err)
	}

	l.logger.Info().Str("socket", l.socketPath).Msg("LMTP server listening")
	return l.server.Serve(listener)
}

// Close shuts down the LMTP server.
func (l *LMTPServer) Close() error {
	return l.server.Close()
}

// lmtpSession handles a single LMTP connection.
type lmtpSession struct {
	queries  *sqlcdb.Queries
	inboxSvc *inbox.Service
	logger   zerolog.Logger

	from       string
	recipients []string
}

func (s *lmtpSession) Reset() {
	s.from = ""
	s.recipients = nil
}

func (s *lmtpSession) Logout() error {
	return nil
}

func (s *lmtpSession) Mail(from string, _ *gosmtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *lmtpSession) Rcpt(to string, _ *gosmtp.RcptOptions) error {
	s.recipients = append(s.recipients, to)
	return nil
}

// Data is not used for LMTP — LMTPData is called instead.
func (s *lmtpSession) Data(_ io.Reader) error {
	return fmt.Errorf("unexpected Data call in LMTP mode")
}

// LMTPData processes the incoming message for each recipient.
func (s *lmtpSession) LMTPData(r io.Reader, status gosmtp.StatusCollector) error {
	msg, err := mail.ReadMessage(r)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to parse message")
		return err
	}

	// Extract headers.
	incoming := &inbox.IncomingMessage{
		MessageID:  msg.Header.Get("Message-Id"),
		InReplyTo:  msg.Header.Get("In-Reply-To"),
		References: msg.Header.Get("References"),
		Subject:    msg.Header.Get("Subject"),
		ReceivedAt: time.Now().UTC(),
	}

	// Parse From address.
	if fromList, err := msg.Header.AddressList("From"); err == nil && len(fromList) > 0 {
		incoming.From = fromList[0].Address
		incoming.FromName = fromList[0].Name
	} else {
		incoming.From = s.from
	}

	// Parse To addresses.
	if toList, err := msg.Header.AddressList("To"); err == nil {
		for _, addr := range toList {
			incoming.To = append(incoming.To, addr.Address)
		}
	}

	// Parse CC addresses.
	if ccList, err := msg.Header.AddressList("Cc"); err == nil {
		for _, addr := range ccList {
			incoming.CC = append(incoming.CC, addr.Address)
		}
	}

	// Parse Date header.
	if date, err := msg.Header.Date(); err == nil {
		incoming.ReceivedAt = date.UTC()
	}

	// Extract body content.
	textBody, htmlBody := extractBody(msg)
	incoming.TextBody = textBody
	incoming.HTMLBody = htmlBody

	// Generate snippet from text body.
	if textBody != "" {
		incoming.Snippet = truncate(textBody, 200)
	} else if htmlBody != "" {
		incoming.Snippet = truncate(stripHTML(htmlBody), 200)
	}

	ctx := context.Background()

	// Deliver to each recipient.
	for _, rcpt := range s.recipients {
		// Extract domain from recipient address.
		parts := strings.SplitN(rcpt, "@", 2)
		if len(parts) != 2 {
			status.SetStatus(rcpt, &gosmtp.SMTPError{
				Code:    550,
				Message: "invalid recipient address",
			})
			continue
		}
		domainName := parts[1]

		// Look up the domain to get org_id and domain_id.
		domain, err := s.queries.GetDomainByNameForInbound(ctx, domainName)
		if err != nil {
			s.logger.Warn().Err(err).Str("domain", domainName).Str("recipient", rcpt).Msg("domain not found for inbound")
			status.SetStatus(rcpt, &gosmtp.SMTPError{
				Code:    550,
				Message: fmt.Sprintf("unknown domain: %s", domainName),
			})
			continue
		}

		// Deliver to inbox.
		_, err = s.inboxSvc.ReceiveMessage(ctx, domain.OrgID, domain.ID, incoming)
		if err != nil {
			s.logger.Error().Err(err).Str("recipient", rcpt).Msg("failed to receive message")
			status.SetStatus(rcpt, &gosmtp.SMTPError{
				Code:    451,
				Message: "temporary delivery failure",
			})
			continue
		}

		s.logger.Info().
			Str("from", incoming.From).
			Str("to", rcpt).
			Str("subject", incoming.Subject).
			Msg("inbound message delivered")
		status.SetStatus(rcpt, nil)
	}

	return nil
}

// extractBody reads the message body and returns text and HTML parts.
func extractBody(msg *mail.Message) (text, html string) {
	contentType := msg.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "text/plain"
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		// Fall back to reading the whole body as text.
		body, _ := io.ReadAll(io.LimitReader(msg.Body, 1<<20))
		return string(body), ""
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		boundary := params["boundary"]
		if boundary == "" {
			body, _ := io.ReadAll(io.LimitReader(msg.Body, 1<<20))
			return string(body), ""
		}
		return extractMultipart(msg.Body, boundary)
	}

	body, _ := io.ReadAll(io.LimitReader(msg.Body, 1<<20))
	bodyStr := string(body)

	if strings.HasPrefix(mediaType, "text/html") {
		return "", bodyStr
	}
	return bodyStr, ""
}

// extractMultipart recursively extracts text and HTML parts from a multipart message.
func extractMultipart(r io.Reader, boundary string) (text, html string) {
	mr := multipart.NewReader(r, boundary)
	for {
		part, err := mr.NextPart()
		if err != nil {
			break
		}

		partType := part.Header.Get("Content-Type")
		mediaType, params, _ := mime.ParseMediaType(partType)

		if strings.HasPrefix(mediaType, "multipart/") {
			if b := params["boundary"]; b != "" {
				t, h := extractMultipart(part, b)
				if text == "" {
					text = t
				}
				if html == "" {
					html = h
				}
			}
			continue
		}

		body, _ := io.ReadAll(io.LimitReader(part, 1<<20))
		switch {
		case strings.HasPrefix(mediaType, "text/plain") && text == "":
			text = string(body)
		case strings.HasPrefix(mediaType, "text/html") && html == "":
			html = string(body)
		}
	}
	return
}

// truncate returns the first n characters of s.
func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n])
}

// stripHTML removes HTML tags from a string (basic implementation for snippets).
func stripHTML(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			result.WriteRune(r)
		}
	}
	return strings.TrimSpace(result.String())
}
