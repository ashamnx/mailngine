package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	oauthStatePrefix  = "oauth_state:"
	oauthStateTTL     = 5 * time.Minute
)

// GoogleOAuth handles Google OAuth2 authentication flows.
type GoogleOAuth struct {
	config *oauth2.Config
	cache  *redis.Client
}

// GoogleUserInfo represents the user information returned by Google's userinfo endpoint.
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

// NewGoogleOAuth creates a new GoogleOAuth instance configured with the given credentials.
func NewGoogleOAuth(clientID, clientSecret, redirectURL string, cache *redis.Client) *GoogleOAuth {
	return &GoogleOAuth{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
		cache: cache,
	}
}

// AuthCodeURL generates an OAuth2 authorization URL with a random state parameter.
// The state is stored in Valkey with a 5-minute TTL for CSRF protection.
func (g *GoogleOAuth) AuthCodeURL(ctx context.Context) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate oauth state: %w", err)
	}
	state := base64.RawURLEncoding.EncodeToString(b)

	key := oauthStatePrefix + state
	if err := g.cache.Set(ctx, key, "1", oauthStateTTL).Err(); err != nil {
		return "", fmt.Errorf("store oauth state: %w", err)
	}

	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// ValidateState checks that the given OAuth state parameter exists in Valkey and deletes it.
// This prevents CSRF attacks and ensures each state is used only once.
func (g *GoogleOAuth) ValidateState(ctx context.Context, state string) error {
	key := oauthStatePrefix + state
	result, err := g.cache.GetDel(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("invalid or expired oauth state")
		}
		return fmt.Errorf("validate oauth state: %w", err)
	}
	if result == "" {
		return fmt.Errorf("invalid or expired oauth state")
	}
	return nil
}

// Exchange exchanges an authorization code for Google user information.
// It performs the OAuth2 token exchange and fetches the user's profile.
func (g *GoogleOAuth) Exchange(ctx context.Context, code string) (*GoogleUserInfo, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange oauth code: %w", err)
	}

	client := g.config.Client(ctx, token)
	resp, err := client.Get(googleUserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("fetch google user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google userinfo returned status %d: %s", resp.StatusCode, body)
	}

	var info GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode google user info: %w", err)
	}

	if info.Email == "" {
		return nil, fmt.Errorf("google user info missing email")
	}

	return &info, nil
}
