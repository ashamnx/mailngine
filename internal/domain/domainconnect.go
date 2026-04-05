package domain

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	dcProviderID = "mailngine.com"
	dcServiceID  = "mail"
)

// ProviderInfo represents a Domain Connect-capable DNS provider.
type ProviderInfo struct {
	ProviderName    string `json:"provider_name"`
	URLSyncUX       string `json:"url_sync_ux"`
	URLAPI          string `json:"url_api"`
	NameServerGroup string `json:"nameserver_group,omitempty"`
}

// DomainConnectResult is the response from checking Domain Connect availability.
type DomainConnectResult struct {
	Supported   bool   `json:"supported"`
	Provider    string `json:"provider,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
	Message     string `json:"message,omitempty"`
}

// DiscoverProvider checks if a domain's DNS provider supports Domain Connect
// by looking up the _domainconnect TXT record.
func DiscoverProvider(ctx context.Context, domainName string) (*ProviderInfo, error) {
	// Try the domain itself and parent domains (for subdomains)
	parts := strings.Split(domainName, ".")
	for i := 0; i < len(parts)-1; i++ {
		candidate := strings.Join(parts[i:], ".")
		txtRecords, err := net.LookupTXT("_domainconnect." + candidate)
		if err != nil || len(txtRecords) == 0 {
			continue
		}

		// Found a Domain Connect TXT record — query the provider's settings endpoint
		for _, txt := range txtRecords {
			// The TXT record contains the Domain Connect API URL
			apiURL := strings.TrimSpace(txt)
			if !strings.HasPrefix(apiURL, "http") {
				continue
			}

			info, err := fetchProviderSettings(ctx, apiURL)
			if err != nil {
				continue
			}
			return info, nil
		}
	}

	return nil, fmt.Errorf("no Domain Connect support found for %s", domainName)
}

// fetchProviderSettings queries the Domain Connect settings endpoint.
func fetchProviderSettings(ctx context.Context, apiURL string) (*ProviderInfo, error) {
	settingsURL := strings.TrimSuffix(apiURL, "/") + "/v2/domainTemplates/providers/" + dcProviderID + "/services/" + dcServiceID

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", settingsURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Even if the template isn't found, we can still use the base URL for the sync flow
	// The settings URL structure tells us the provider supports Domain Connect

	// Try to get provider name from the base settings endpoint
	baseSettingsURL := strings.TrimSuffix(apiURL, "/") + "/v2/domainTemplates/providers"
	baseReq, _ := http.NewRequestWithContext(ctx, "GET", baseSettingsURL, nil)
	baseResp, err := client.Do(baseReq)
	if err == nil {
		defer baseResp.Body.Close()
	}

	// Derive provider name from the API URL
	providerName := deriveProviderName(apiURL)

	// The sync UX URL is typically at the same host as the API
	parsed, _ := url.Parse(apiURL)
	syncUXURL := fmt.Sprintf("%s://%s/domain-connect/v2/domainTemplates/providers/%s/services/%s/apply",
		parsed.Scheme, parsed.Host, dcProviderID, dcServiceID)

	return &ProviderInfo{
		ProviderName: providerName,
		URLSyncUX:    syncUXURL,
		URLAPI:       apiURL,
	}, nil
}

// deriveProviderName guesses the provider name from the API URL.
func deriveProviderName(apiURL string) string {
	lower := strings.ToLower(apiURL)
	switch {
	case strings.Contains(lower, "cloudflare"):
		return "Cloudflare"
	case strings.Contains(lower, "godaddy"):
		return "GoDaddy"
	case strings.Contains(lower, "ionos"):
		return "IONOS"
	case strings.Contains(lower, "wordpress"):
		return "WordPress.com"
	case strings.Contains(lower, "namesilo"):
		return "NameSilo"
	case strings.Contains(lower, "plesk"):
		return "Plesk"
	case strings.Contains(lower, "vercel"):
		return "Vercel"
	default:
		return "DNS Provider"
	}
}

// GenerateRedirectURL builds the Domain Connect synchronous flow redirect URL.
func GenerateRedirectURL(provider *ProviderInfo, domainName string, dkimSelector, dkimPublicKey, redirectURI, state, signingKey string) (string, error) {
	// Build the sync URL with query parameters
	params := url.Values{}
	params.Set("domain", domainName)
	params.Set("providerName", "Mailngine")
	params.Set("mn_selector", dkimSelector)
	params.Set("dkim_public_key", dkimPublicKey)
	params.Set("redirect_uri", redirectURI)
	params.Set("state", state)

	fullURL := provider.URLSyncUX + "?" + params.Encode()

	// Sign the URL if signing key is provided
	if signingKey != "" {
		sig := signRequest(params.Encode(), signingKey)
		fullURL += "&sig=" + url.QueryEscape(sig)
	}

	return fullURL, nil
}

// signRequest generates an HMAC-SHA256 signature of the query parameters.
func signRequest(queryString, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(queryString))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// GenerateState creates a random state token for CSRF protection.
func GenerateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// CheckDomainConnectSupport is a quick check that returns whether a domain
// has a _domainconnect TXT record (without full provider discovery).
func CheckDomainConnectSupport(domainName string) bool {
	parts := strings.Split(domainName, ".")
	for i := 0; i < len(parts)-1; i++ {
		candidate := strings.Join(parts[i:], ".")
		records, err := net.LookupTXT("_domainconnect." + candidate)
		if err == nil && len(records) > 0 {
			return true
		}
	}
	return false
}

// Silence unused import warnings
var _ = json.Marshal
var _ = io.ReadAll
