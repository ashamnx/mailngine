package mailngine

import "time"

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

// SendEmailParams contains the parameters for sending an email.
type SendEmailParams struct {
	From           string            `json:"from"`
	To             []string          `json:"to"`
	CC             []string          `json:"cc,omitempty"`
	BCC            []string          `json:"bcc,omitempty"`
	ReplyTo        string            `json:"reply_to,omitempty"`
	Subject        string            `json:"subject"`
	HTML           string            `json:"html,omitempty"`
	Text           string            `json:"text,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	Tags           []string          `json:"tags,omitempty"`
	TemplateID     string            `json:"template_id,omitempty"`
	TemplateData   map[string]string `json:"template_data,omitempty"`
	IdempotencyKey string            `json:"idempotency_key,omitempty"`
	ScheduledAt    *time.Time        `json:"scheduled_at,omitempty"`
}

// Domain represents a verified sending domain.
type Domain struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Status        string      `json:"status"`
	Region        string      `json:"region"`
	DkimSelector  string      `json:"dkim_selector"`
	OpenTracking  bool        `json:"open_tracking"`
	ClickTracking bool        `json:"click_tracking"`
	DNSRecords    []DNSRecord `json:"dns_records,omitempty"`
	VerifiedAt    *time.Time  `json:"verified_at,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// DNSRecord represents a DNS record required for domain verification.
type DNSRecord struct {
	ID         string     `json:"id"`
	DomainID   string     `json:"domain_id"`
	RecordType string     `json:"record_type"`
	Host       string     `json:"host"`
	Value      string     `json:"value"`
	Purpose    string     `json:"purpose"`
	Status     string     `json:"status"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// CreateDomainResponse is the response from creating a domain.
type CreateDomainResponse struct {
	Domain     Domain      `json:"domain"`
	DNSRecords []DNSRecord `json:"dns_records"`
}

// UpdateDomainParams contains the parameters for updating a domain.
type UpdateDomainParams struct {
	OpenTracking  *bool `json:"open_tracking,omitempty"`
	ClickTracking *bool `json:"click_tracking,omitempty"`
}

// Webhook represents a configured webhook endpoint.
type Webhook struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	IsActive  bool      `json:"is_active"`
	Secret    string    `json:"secret,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateWebhookParams contains the parameters for creating a webhook.
type CreateWebhookParams struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

// UpdateWebhookParams contains the parameters for updating a webhook.
type UpdateWebhookParams struct {
	URL      string   `json:"url"`
	Events   []string `json:"events"`
	IsActive bool     `json:"is_active"`
}

// Template represents an email template.
type Template struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Subject   string    `json:"subject"`
	HTMLBody  string    `json:"html_body"`
	TextBody  string    `json:"text_body,omitempty"`
	Variables []string  `json:"variables,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTemplateParams contains the parameters for creating a template.
type CreateTemplateParams struct {
	Name      string   `json:"name"`
	Subject   string   `json:"subject"`
	HTMLBody  string   `json:"html_body"`
	TextBody  string   `json:"text_body,omitempty"`
	Variables []string `json:"variables,omitempty"`
}

// UpdateTemplateParams contains the parameters for updating a template.
type UpdateTemplateParams struct {
	Name      string   `json:"name"`
	Subject   string   `json:"subject"`
	HTMLBody  string   `json:"html_body"`
	TextBody  string   `json:"text_body,omitempty"`
	Variables []string `json:"variables,omitempty"`
}

// PreviewTemplateParams contains the parameters for previewing a template.
type PreviewTemplateParams struct {
	Data map[string]string `json:"data"`
}

// TemplatePreview is the response from previewing a template.
type TemplatePreview struct {
	Subject  string `json:"subject"`
	HTMLBody string `json:"html_body"`
	TextBody string `json:"text_body"`
}

// APIKey represents an API key. The full Key value is only populated on creation.
type APIKey struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	Key        string     `json:"key,omitempty"`
	Permission string     `json:"permission"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// CreateAPIKeyParams contains the parameters for creating an API key.
type CreateAPIKeyParams struct {
	Name       string     `json:"name"`
	Permission string     `json:"permission,omitempty"`
	DomainID   *string    `json:"domain_id,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

// ListOptions contains pagination parameters for list endpoints.
type ListOptions struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
}

// ListMeta contains pagination metadata returned by list endpoints.
type ListMeta struct {
	Page    int  `json:"page"`
	PerPage int  `json:"per_page"`
	Total   int  `json:"total"`
	HasMore bool `json:"has_more"`
}

// ListResponse is a generic paginated list response.
type ListResponse[T any] struct {
	Data []T      `json:"data"`
	Meta ListMeta `json:"meta"`
}
