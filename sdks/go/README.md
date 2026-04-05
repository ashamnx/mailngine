# mailngine-go

Official Go SDK for the [Mailngine](https://mailngine.com) email API.

## Installation

```bash
go get github.com/mailngine/mailngine-go
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mailngine/mailngine-go/mailngine"
)

func main() {
	client := mailngine.New("mn_live_...")

	ctx := context.Background()

	// Send an email
	email, err := client.Emails.Send(ctx, &mailngine.SendEmailParams{
		From:    "hello@example.com",
		To:      []string{"user@example.com"},
		Subject: "Hello!",
		HTML:    "<h1>Welcome</h1>",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent email: %s\n", email.ID)

	// Get email status
	email, err = client.Emails.Get(ctx, email.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %s\n", email.Status)

	// List emails with pagination
	list, err := client.Emails.List(ctx, &mailngine.ListOptions{
		Page:    1,
		PerPage: 20,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total emails: %d\n", list.Meta.Total)
}
```

## Resources

### Emails

```go
client.Emails.Send(ctx, params)     // Send an email
client.Emails.Get(ctx, id)          // Get email by ID
client.Emails.List(ctx, opts)       // List emails (paginated)
```

### Domains

```go
client.Domains.Create(ctx, name)    // Register a new domain
client.Domains.Get(ctx, id)         // Get domain by ID
client.Domains.List(ctx)            // List all domains
client.Domains.Verify(ctx, id)      // Trigger DNS verification
client.Domains.Update(ctx, id, p)   // Update domain settings
client.Domains.Delete(ctx, id)      // Delete a domain
```

### Webhooks

```go
client.Webhooks.Create(ctx, params) // Create a webhook
client.Webhooks.Get(ctx, id)        // Get webhook by ID
client.Webhooks.List(ctx)           // List all webhooks
client.Webhooks.Update(ctx, id, p)  // Update a webhook
client.Webhooks.Delete(ctx, id)     // Delete a webhook
```

### Templates

```go
client.Templates.Create(ctx, params)     // Create a template
client.Templates.Get(ctx, id)            // Get template by ID
client.Templates.List(ctx)               // List all templates
client.Templates.Update(ctx, id, params) // Update a template
client.Templates.Delete(ctx, id)         // Delete a template
client.Templates.Preview(ctx, id, data)  // Preview with data
```

### API Keys

```go
client.APIKeys.Create(ctx, params)  // Create an API key
client.APIKeys.List(ctx)            // List all API keys
client.APIKeys.Revoke(ctx, id)      // Revoke an API key
```

## Error Handling

All API errors are returned as `*mailngine.APIError`. Use the provided helper
functions to check for specific error conditions:

```go
email, err := client.Emails.Get(ctx, "nonexistent")
if mailngine.IsNotFound(err) {
	// Handle 404
}
if mailngine.IsRateLimited(err) {
	// Handle 429 - the SDK already retries these automatically
}
if mailngine.IsValidationError(err) {
	// Handle 400
}

// Access the full error details
var apiErr *mailngine.APIError
if errors.As(err, &apiErr) {
	fmt.Printf("Code: %s, Message: %s\n", apiErr.Code, apiErr.Message)
}
```

## Configuration

```go
// Custom base URL (self-hosted)
client := mailngine.New("mn_live_...",
	mailngine.WithBaseURL("https://api.your-instance.com"),
)

// Custom HTTP client
client := mailngine.New("mn_live_...",
	mailngine.WithHTTPClient(&http.Client{
		Timeout: 60 * time.Second,
	}),
)
```

## Automatic Retries

The SDK automatically retries requests that receive a 429 (rate limited) or 5xx
(server error) response. It retries up to 3 times with exponential backoff
(1s, 2s, 4s). Non-retryable errors (4xx other than 429) are returned immediately.
