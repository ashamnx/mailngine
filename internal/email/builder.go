package email

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"
	"time"
)

// BuildMIMEMessage constructs an RFC 5322 compliant MIME message from the given parameters.
// If both textBody and htmlBody are provided, a multipart/alternative message is built.
// If only one is provided, a simple single-part message is built.
func BuildMIMEMessage(from, subject, textBody, htmlBody string, to, cc, bcc []string, replyTo string, headers map[string]string, messageID string) ([]byte, error) {
	var buf bytes.Buffer

	// Write standard headers.
	writeHeader(&buf, "From", from)
	writeHeader(&buf, "To", strings.Join(to, ", "))
	if len(cc) > 0 {
		writeHeader(&buf, "Cc", strings.Join(cc, ", "))
	}
	// BCC recipients are intentionally omitted from headers per RFC 5322.
	writeHeader(&buf, "Subject", subject)
	writeHeader(&buf, "Date", time.Now().UTC().Format(time.RFC1123Z))
	writeHeader(&buf, "Message-ID", messageID)
	writeHeader(&buf, "MIME-Version", "1.0")

	if replyTo != "" {
		writeHeader(&buf, "Reply-To", replyTo)
	}

	// Write custom headers.
	for k, v := range headers {
		// Prevent overriding critical headers.
		normalized := strings.ToLower(k)
		switch normalized {
		case "from", "to", "cc", "bcc", "subject", "date", "message-id", "mime-version":
			continue
		}
		writeHeader(&buf, k, v)
	}

	hasText := textBody != ""
	hasHTML := htmlBody != ""

	switch {
	case hasText && hasHTML:
		if err := writeMultipartAlternative(&buf, textBody, htmlBody); err != nil {
			return nil, fmt.Errorf("build multipart message: %w", err)
		}
	case hasHTML:
		writeHeader(&buf, "Content-Type", "text/html; charset=utf-8")
		writeHeader(&buf, "Content-Transfer-Encoding", "quoted-printable")
		buf.WriteString("\r\n")
		buf.WriteString(htmlBody)
	case hasText:
		writeHeader(&buf, "Content-Type", "text/plain; charset=utf-8")
		writeHeader(&buf, "Content-Transfer-Encoding", "quoted-printable")
		buf.WriteString("\r\n")
		buf.WriteString(textBody)
	default:
		// Empty body.
		writeHeader(&buf, "Content-Type", "text/plain; charset=utf-8")
		buf.WriteString("\r\n")
	}

	return buf.Bytes(), nil
}

// writeHeader writes a single MIME header line.
// It strips CR and LF characters from the value to prevent header injection.
func writeHeader(w io.Writer, key, value string) {
	sanitized := strings.NewReplacer("\r", "", "\n", "").Replace(value)
	fmt.Fprintf(w, "%s: %s\r\n", key, sanitized)
}

// writeMultipartAlternative writes a multipart/alternative body with text/plain and text/html parts.
func writeMultipartAlternative(buf *bytes.Buffer, textBody, htmlBody string) error {
	boundary := generateBoundary()

	writeHeader(buf, "Content-Type", fmt.Sprintf("multipart/alternative; boundary=%q", boundary))
	buf.WriteString("\r\n")

	w := multipart.NewWriter(buf)
	if err := w.SetBoundary(boundary); err != nil {
		return fmt.Errorf("set boundary: %w", err)
	}

	// Text part.
	textHeader := make(textproto.MIMEHeader)
	textHeader.Set("Content-Type", "text/plain; charset=utf-8")
	textHeader.Set("Content-Transfer-Encoding", "quoted-printable")
	textPart, err := w.CreatePart(textHeader)
	if err != nil {
		return fmt.Errorf("create text part: %w", err)
	}
	if _, err := textPart.Write([]byte(textBody)); err != nil {
		return fmt.Errorf("write text part: %w", err)
	}

	// HTML part.
	htmlHeader := make(textproto.MIMEHeader)
	htmlHeader.Set("Content-Type", "text/html; charset=utf-8")
	htmlHeader.Set("Content-Transfer-Encoding", "quoted-printable")
	htmlPart, err := w.CreatePart(htmlHeader)
	if err != nil {
		return fmt.Errorf("create html part: %w", err)
	}
	if _, err := htmlPart.Write([]byte(htmlBody)); err != nil {
		return fmt.Errorf("write html part: %w", err)
	}

	return w.Close()
}

// generateBoundary creates a random MIME boundary string.
func generateBoundary() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to a timestamp-based boundary if crypto/rand fails.
		return fmt.Sprintf("mailngine-%d", time.Now().UnixNano())
	}
	return "mailngine-" + hex.EncodeToString(b)
}
