package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	sqlcdb "github.com/hellomail/hellomail/internal/db/sqlcdb"
)

const cloudflareAPIBase = "https://api.cloudflare.com/client/v4"

// FindCFZoneID looks up the Cloudflare zone ID for a domain name.
// It handles subdomains by checking progressively higher levels (e.g., mail.example.com → example.com).
func FindCFZoneID(ctx context.Context, apiToken, domainName string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// Try the domain and parent domains (for subdomains like mail.example.com)
	parts := strings.Split(domainName, ".")
	for i := 0; i < len(parts)-1; i++ {
		candidate := strings.Join(parts[i:], ".")

		url := fmt.Sprintf("%s/zones?name=%s&per_page=1", cloudflareAPIBase, candidate)
		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+apiToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result struct {
			Success bool `json:"success"`
			Result  []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"result"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			continue
		}
		if result.Success && len(result.Result) > 0 {
			return result.Result[0].ID, nil
		}
	}

	return "", fmt.Errorf("zone not found for domain %s in Cloudflare account", domainName)
}

// cloudflareCreateRequest represents the request body for creating a DNS record
// via the Cloudflare API.
type cloudflareCreateRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	// Priority is used only for MX records.
	Priority *int `json:"priority,omitempty"`
}

// cloudflareResponse represents the top-level response from the Cloudflare API.
type cloudflareResponse struct {
	Success bool                `json:"success"`
	Errors  []cloudflareError   `json:"errors"`
}

// cloudflareError represents an error returned by the Cloudflare API.
type cloudflareError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// AutoCreateDNS creates DNS records in Cloudflare for the given zone using the
// Cloudflare API v4. Each record from the database is translated to a Cloudflare
// DNS record creation request.
func AutoCreateDNS(ctx context.Context, zoneID, apiToken string, records []sqlcdb.DnsRecord) error {
	client := &http.Client{Timeout: 30 * time.Second}

	for _, rec := range records {
		cfReq := buildCloudflareRequest(rec)

		body, err := json.Marshal(cfReq)
		if err != nil {
			return fmt.Errorf("marshaling request for %s record: %w", rec.Purpose, err)
		}

		url := fmt.Sprintf("%s/zones/%s/dns_records", cloudflareAPIBase, zoneID)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("creating request for %s record: %w", rec.Purpose, err)
		}

		req.Header.Set("Authorization", "Bearer "+apiToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("calling Cloudflare API for %s record: %w", rec.Purpose, err)
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("reading Cloudflare response for %s record: %w", rec.Purpose, err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			var cfResp cloudflareResponse
			if jsonErr := json.Unmarshal(respBody, &cfResp); jsonErr == nil && len(cfResp.Errors) > 0 {
				// Error code 81057 means "Record already exists", which is acceptable.
				if cfResp.Errors[0].Code == 81057 {
					continue
				}
				return fmt.Errorf("cloudflare API error for %s record: %s", rec.Purpose, cfResp.Errors[0].Message)
			}
			return fmt.Errorf("cloudflare API returned status %d for %s record", resp.StatusCode, rec.Purpose)
		}
	}

	return nil
}

// buildCloudflareRequest converts a DNS record from the database into a
// Cloudflare API create request.
func buildCloudflareRequest(rec sqlcdb.DnsRecord) cloudflareCreateRequest {
	cfReq := cloudflareCreateRequest{
		Type: rec.RecordType,
		Name: rec.Host,
		TTL:  3600,
	}

	switch rec.RecordType {
	case "MX":
		// Value format: "10 mx.hellomail.dev"
		parts := strings.Fields(rec.Value)
		if len(parts) >= 2 {
			priority := 10
			cfReq.Priority = &priority
			cfReq.Content = parts[1]
			// Parse priority from value if present
			if _, err := fmt.Sscanf(parts[0], "%d", &priority); err == nil {
				cfReq.Priority = &priority
			}
		} else {
			cfReq.Content = rec.Value
		}
	default:
		cfReq.Content = rec.Value
	}

	return cfReq
}
