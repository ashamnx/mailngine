package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port        string
	Env         string
	DatabaseURL string
	ValkeyURL   string

	Spaces SpacesConfig

	Google GoogleConfig
	JWT    JWTConfig

	DKIMMasterKey string
	FrontendURL   string

	Postfix PostfixConfig

	WorkerShutdownTimeout time.Duration
}

type SpacesConfig struct {
	Endpoint  string
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type JWTConfig struct {
	Secret string
	Expiry time.Duration
}

type PostfixConfig struct {
	SMTPHost string
	SMTPPort string
}

func Load() (*Config, error) {
	jwtExpiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY: %w", err)
	}

	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
		DatabaseURL: requireEnv("DATABASE_URL"),
		ValkeyURL:   requireEnv("VALKEY_URL"),

		Spaces: SpacesConfig{
			Endpoint:  getEnv("SPACES_ENDPOINT", ""),
			Region:    getEnv("SPACES_REGION", "sgp1"),
			Bucket:    getEnv("SPACES_BUCKET", "mailngine"),
			AccessKey: getEnv("SPACES_ACCESS_KEY", ""),
			SecretKey: getEnv("SPACES_SECRET_KEY", ""),
		},

		Google: GoogleConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/v1/auth/google/callback"),
		},

		JWT: JWTConfig{
			Secret: requireEnv("JWT_SECRET"),
			Expiry: jwtExpiry,
		},

		DKIMMasterKey: getEnv("DKIM_MASTER_KEY", ""),
		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:4200"),

		Postfix: PostfixConfig{
			SMTPHost: getEnv("POSTFIX_SMTP_HOST", "127.0.0.1"),
			SMTPPort: getEnv("POSTFIX_SMTP_PORT", "25"),
		},
	}

	workerTimeout, err := time.ParseDuration(getEnv("WORKER_SHUTDOWN_TIMEOUT", "60s"))
	if err != nil {
		return nil, fmt.Errorf("invalid WORKER_SHUTDOWN_TIMEOUT: %w", err)
	}
	cfg.WorkerShutdownTimeout = workerTimeout

	// Reject weak JWT secrets in non-development environments.
	if cfg.Env != "development" && len(cfg.JWT.Secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters in production")
	}

	return cfg, nil
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return v
}
