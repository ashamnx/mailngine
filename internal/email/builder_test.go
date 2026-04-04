package email

import (
	"strings"
	"testing"
)

func TestBuildMIMEMessage_TextOnly(t *testing.T) {
	msg, err := BuildMIMEMessage(
		"sender@example.com",
		"Test Subject",
		"Hello, plain text!",
		"",
		[]string{"to@example.com"},
		nil, nil, "",
		nil,
		"<msg-001@example.com>",
	)
	if err != nil {
		t.Fatalf("BuildMIMEMessage() error = %v", err)
	}

	body := string(msg)

	if !strings.Contains(body, "Content-Type: text/plain; charset=utf-8") {
		t.Error("text-only message should have text/plain Content-Type")
	}
	if strings.Contains(body, "multipart/alternative") {
		t.Error("text-only message should not be multipart/alternative")
	}
	if !strings.Contains(body, "Hello, plain text!") {
		t.Error("message body should contain the text content")
	}
}

func TestBuildMIMEMessage_HTMLOnly(t *testing.T) {
	msg, err := BuildMIMEMessage(
		"sender@example.com",
		"HTML Subject",
		"",
		"<h1>Hello</h1>",
		[]string{"to@example.com"},
		nil, nil, "",
		nil,
		"<msg-002@example.com>",
	)
	if err != nil {
		t.Fatalf("BuildMIMEMessage() error = %v", err)
	}

	body := string(msg)

	if !strings.Contains(body, "Content-Type: text/html; charset=utf-8") {
		t.Error("HTML-only message should have text/html Content-Type")
	}
	if strings.Contains(body, "multipart/alternative") {
		t.Error("HTML-only message should not be multipart/alternative")
	}
	if !strings.Contains(body, "<h1>Hello</h1>") {
		t.Error("message body should contain the HTML content")
	}
}

func TestBuildMIMEMessage_MultipartAlternative(t *testing.T) {
	msg, err := BuildMIMEMessage(
		"sender@example.com",
		"Multipart Subject",
		"Plain text version",
		"<p>HTML version</p>",
		[]string{"to@example.com"},
		nil, nil, "",
		nil,
		"<msg-003@example.com>",
	)
	if err != nil {
		t.Fatalf("BuildMIMEMessage() error = %v", err)
	}

	body := string(msg)

	if !strings.Contains(body, "multipart/alternative") {
		t.Error("message with both text and HTML should be multipart/alternative")
	}
	if !strings.Contains(body, "text/plain; charset=utf-8") {
		t.Error("multipart message should contain text/plain part")
	}
	if !strings.Contains(body, "text/html; charset=utf-8") {
		t.Error("multipart message should contain text/html part")
	}
	if !strings.Contains(body, "Plain text version") {
		t.Error("multipart message should contain the text body")
	}
	if !strings.Contains(body, "<p>HTML version</p>") {
		t.Error("multipart message should contain the HTML body")
	}
}

func TestBuildMIMEMessage_StandardHeaders(t *testing.T) {
	msg, err := BuildMIMEMessage(
		"sender@example.com",
		"Header Test",
		"body",
		"",
		[]string{"alice@example.com", "bob@example.com"},
		[]string{"cc@example.com"},
		[]string{"bcc@example.com"},
		"reply@example.com",
		map[string]string{"X-Custom": "value"},
		"<msg-004@example.com>",
	)
	if err != nil {
		t.Fatalf("BuildMIMEMessage() error = %v", err)
	}

	body := string(msg)

	tests := []struct {
		name     string
		contains string
	}{
		{"From header", "From: sender@example.com"},
		{"To header", "To: alice@example.com, bob@example.com"},
		{"Cc header", "Cc: cc@example.com"},
		{"Subject header", "Subject: Header Test"},
		{"Date header", "Date: "},
		{"Message-ID header", "Message-ID: <msg-004@example.com>"},
		{"MIME-Version header", "MIME-Version: 1.0"},
		{"Reply-To header", "Reply-To: reply@example.com"},
		{"Custom header", "X-Custom: value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(body, tt.contains) {
				t.Errorf("message should contain %q", tt.contains)
			}
		})
	}

	// BCC should NOT appear in headers per RFC 5322.
	if strings.Contains(body, "Bcc:") || strings.Contains(body, "bcc@example.com") {
		t.Error("BCC recipients should not appear in message headers")
	}
}

func TestBuildMIMEMessage_CriticalHeadersNotOverridden(t *testing.T) {
	msg, err := BuildMIMEMessage(
		"real@example.com",
		"Real Subject",
		"body",
		"",
		[]string{"to@example.com"},
		nil, nil, "",
		map[string]string{
			"From":    "attacker@evil.com",
			"Subject": "Spoofed Subject",
		},
		"<msg-005@example.com>",
	)
	if err != nil {
		t.Fatalf("BuildMIMEMessage() error = %v", err)
	}

	body := string(msg)

	if strings.Contains(body, "attacker@evil.com") {
		t.Error("custom headers should not override From")
	}
	if strings.Contains(body, "Spoofed Subject") {
		t.Error("custom headers should not override Subject")
	}
}
