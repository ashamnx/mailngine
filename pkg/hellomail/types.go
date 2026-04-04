package hellomail

import "time"

// SendEmailRequest represents a request to send an email.
type SendEmailRequest struct {
	From           string            `json:"from" validate:"required,email"`
	To             []string          `json:"to" validate:"required,min=1,dive,email"`
	CC             []string          `json:"cc,omitempty" validate:"omitempty,dive,email"`
	BCC            []string          `json:"bcc,omitempty" validate:"omitempty,dive,email"`
	ReplyTo        string            `json:"reply_to,omitempty" validate:"omitempty,email"`
	Subject        string            `json:"subject" validate:"required,max=998"`
	HTML           string            `json:"html,omitempty"`
	Text           string            `json:"text,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	Tags           []string          `json:"tags,omitempty"`
	TemplateID     string            `json:"template_id,omitempty"`
	TemplateData   map[string]string `json:"template_data,omitempty"`
	IdempotencyKey string            `json:"idempotency_key,omitempty"`
	ScheduledAt    *time.Time        `json:"scheduled_at,omitempty"`
}

// Email represents a sent email.
type Email struct {
	ID          string     `json:"id"`
	From        string     `json:"from"`
	To          []string   `json:"to"`
	Subject     string     `json:"subject"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
}

// Domain represents a verified sending domain.
type Domain struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Status        string     `json:"status"`
	DNSRecords    []DNSRecord `json:"dns_records,omitempty"`
	OpenTracking  bool       `json:"open_tracking"`
	ClickTracking bool       `json:"click_tracking"`
	CreatedAt     time.Time  `json:"created_at"`
	VerifiedAt    *time.Time `json:"verified_at,omitempty"`
}

// DNSRecord represents a DNS record that must be configured for domain verification.
type DNSRecord struct {
	ID         string     `json:"id"`
	RecordType string     `json:"record_type"`
	Host       string     `json:"host"`
	Value      string     `json:"value"`
	Purpose    string     `json:"purpose"`
	Status     string     `json:"status"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

// Webhook represents a configured webhook endpoint.
type Webhook struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// APIKey represents an API key (prefix only, never the full key after creation).
type APIKey struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	Permission string     `json:"permission"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// Organization represents a tenant organization.
type Organization struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Plan         string    `json:"plan"`
	MonthlyLimit int       `json:"monthly_limit"`
	CreatedAt    time.Time `json:"created_at"`
}

// User represents an authenticated user.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
