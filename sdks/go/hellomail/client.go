// Package hellomail provides a Go client for the Hello Mail email API.
//
// Create a client with your API key:
//
//	client := hellomail.New("hm_live_...")
//
// Then use the resource fields to interact with the API:
//
//	email, err := client.Emails.Send(ctx, &hellomail.SendEmailParams{
//	    From:    "hello@example.com",
//	    To:      []string{"user@example.com"},
//	    Subject: "Hello!",
//	    HTML:    "<h1>Welcome</h1>",
//	})
package hellomail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default Hello Mail API base URL.
	DefaultBaseURL = "https://api.hellomail.dev"

	// Version is the SDK version.
	Version = "0.1.0"

	defaultTimeout    = 30 * time.Second
	maxRetries        = 3
	initialRetryDelay = 1 * time.Second
)

// Client is the Hello Mail API client. Use New to create one.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	// Emails provides access to the email sending and retrieval API.
	Emails *EmailsResource

	// Domains provides access to the domain management API.
	Domains *DomainsResource

	// Webhooks provides access to the webhook management API.
	Webhooks *WebhooksResource

	// Templates provides access to the template management API.
	Templates *TemplatesResource

	// APIKeys provides access to the API key management API.
	APIKeys *APIKeysResource
}

// Option configures the Client.
type Option func(*Client)

// WithBaseURL sets a custom base URL for the API. This is useful for testing
// or when using a self-hosted Hello Mail instance.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client. Use this to configure proxies,
// custom TLS settings, or other transport-level options.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// New creates a new Hello Mail API client with the given API key.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Emails = &EmailsResource{client: c}
	c.Domains = &DomainsResource{client: c}
	c.Webhooks = &WebhooksResource{client: c}
	c.Templates = &TemplatesResource{client: c}
	c.APIKeys = &APIKeysResource{client: c}

	return c
}

// envelope is the standard API response wrapper.
type envelope struct {
	Data  json.RawMessage `json:"data"`
	Error *APIError       `json:"error,omitempty"`
	Meta  *ListMeta       `json:"meta,omitempty"`
}

// do executes an API request with automatic retry on 429 and 5xx responses.
//
// If body is non-nil, it is marshaled to JSON and sent as the request body.
// If result is non-nil, the response data field is unmarshaled into it.
func (c *Client) do(ctx context.Context, method, path string, body, result any) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("hellomail: marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			delay := initialRetryDelay << (attempt - 1) // 1s, 2s, 4s
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}

			// Reset body reader for retry
			if body != nil {
				data, _ := json.Marshal(body)
				reqBody = bytes.NewReader(data)
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
		if err != nil {
			return fmt.Errorf("hellomail: create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("User-Agent", "hellomail-go/"+Version)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("hellomail: send request: %w", err)
			// Network errors are retryable
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("hellomail: read response: %w", err)
			continue
		}

		// Retry on 429 (rate limited) or 5xx (server error)
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			var env envelope
			if err := json.Unmarshal(respBody, &env); err == nil && env.Error != nil {
				env.Error.StatusCode = resp.StatusCode
				lastErr = env.Error
			} else {
				lastErr = &APIError{
					StatusCode: resp.StatusCode,
					Code:       http.StatusText(resp.StatusCode),
					Message:    string(respBody),
				}
			}
			continue
		}

		// Non-retryable error
		if resp.StatusCode >= 400 {
			var env envelope
			if err := json.Unmarshal(respBody, &env); err == nil && env.Error != nil {
				env.Error.StatusCode = resp.StatusCode
				return env.Error
			}
			return &APIError{
				StatusCode: resp.StatusCode,
				Code:       http.StatusText(resp.StatusCode),
				Message:    string(respBody),
			}
		}

		// Success - unmarshal if result is expected
		if result != nil {
			var env envelope
			if err := json.Unmarshal(respBody, &env); err != nil {
				return fmt.Errorf("hellomail: unmarshal response: %w", err)
			}
			if err := json.Unmarshal(env.Data, result); err != nil {
				return fmt.Errorf("hellomail: unmarshal data: %w", err)
			}
		}

		return nil
	}

	return lastErr
}

// doWithMeta is like do but also returns pagination metadata.
func (c *Client) doWithMeta(ctx context.Context, method, path string, body, result any) (*ListMeta, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("hellomail: marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			delay := initialRetryDelay << (attempt - 1)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}

			if body != nil {
				data, _ := json.Marshal(body)
				reqBody = bytes.NewReader(data)
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
		if err != nil {
			return nil, fmt.Errorf("hellomail: create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("User-Agent", "hellomail-go/"+Version)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("hellomail: send request: %w", err)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("hellomail: read response: %w", err)
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			var env envelope
			if err := json.Unmarshal(respBody, &env); err == nil && env.Error != nil {
				env.Error.StatusCode = resp.StatusCode
				lastErr = env.Error
			} else {
				lastErr = &APIError{
					StatusCode: resp.StatusCode,
					Code:       http.StatusText(resp.StatusCode),
					Message:    string(respBody),
				}
			}
			continue
		}

		if resp.StatusCode >= 400 {
			var env envelope
			if err := json.Unmarshal(respBody, &env); err == nil && env.Error != nil {
				env.Error.StatusCode = resp.StatusCode
				return nil, env.Error
			}
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Code:       http.StatusText(resp.StatusCode),
				Message:    string(respBody),
			}
		}

		if result != nil {
			var env envelope
			if err := json.Unmarshal(respBody, &env); err != nil {
				return nil, fmt.Errorf("hellomail: unmarshal response: %w", err)
			}
			if err := json.Unmarshal(env.Data, result); err != nil {
				return nil, fmt.Errorf("hellomail: unmarshal data: %w", err)
			}
			return env.Meta, nil
		}

		return nil, nil
	}

	return nil, lastErr
}
