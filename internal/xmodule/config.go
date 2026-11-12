package xmodule

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Version is the build version. Overridable via ldflags.
var Version = "0.1.0-dev"

// Config is the runtime configuration, read from the environment.
type Config struct {
	Listen     string
	DSN        string
	AppName    string
	AppURL     string
	AppLocale  string
	AuthSecret string

	SessionTTL    time.Duration
	CookieSecure  bool
	TrustedOrigin []string

	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	SMTPSecure   bool

	RateLimitPerMinute int
	LogLevel           string
}

func LoadConfig() Config {
	return Config{
		Listen:     envOr("LISTEN", "127.0.0.1:8080"),
		DSN:        strings.TrimSpace(os.Getenv("DATABASE_URL")),
		AppName:    envOr("APP_NAME", "Go React Boilerplate"),
		AppURL:     envOr("APP_URL", "http://localhost:5173"),
		AppLocale:  envOr("APP_LOCALE", "id"),
		AuthSecret: strings.TrimSpace(os.Getenv("AUTH_SECRET")),

		SessionTTL:    time.Duration(envInt("SESSION_TTL_HOURS", 168)) * time.Hour,
		TrustedOrigin: splitList(os.Getenv("TRUSTED_ORIGINS")),

		SMTPHost:     strings.TrimSpace(os.Getenv("SMTP_HOST")),
		SMTPPort:     envInt("SMTP_PORT", 1025),
		SMTPUser:     strings.TrimSpace(os.Getenv("SMTP_USER")),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:     envOr("MAIL_FROM", "hello@example.com"),
		SMTPSecure:   envBool("SMTP_SECURE"),

		RateLimitPerMinute: envInt("RATE_LIMIT_PER_MINUTE", 20),
		LogLevel:           envOr("LOG_LEVEL", "info"),
	}
}

func (c *Config) Validate() error {
	if c.DSN == "" {
		return errors.New("DATABASE_URL is required")
	}
	if c.AppURL == "" {
		return errors.New("APP_URL is required")
	}
	parsed, err := url.Parse(c.AppURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("APP_URL=%q: must be an absolute URL", c.AppURL)
	}
	if c.AppLocale != "en" && c.AppLocale != "id" {
		return fmt.Errorf("APP_LOCALE=%q: must be en or id", c.AppLocale)
	}
	if c.SessionTTL <= 0 {
		return errors.New("SESSION_TTL_HOURS must be positive")
	}
	return nil
}

// ResolveCookieSecure derives the cookie Secure flag from the public URL
// unless COOKIE_SECURE is set explicitly.
func (c *Config) ResolveCookieSecure() {
	if raw := strings.TrimSpace(os.Getenv("COOKIE_SECURE")); raw != "" {
		c.CookieSecure = raw == "1" || strings.EqualFold(raw, "true")
		return
	}
	parsed, err := url.Parse(c.AppURL)
	c.CookieSecure = err == nil && parsed.Scheme == "https"
}

// LoadOrCreateSecret resolves the HMAC secret: environment override, existing
// local file, or a freshly generated one with 0600 permissions. Production
// deployments should set AUTH_SECRET.
func LoadOrCreateSecret(override string) (string, error) {
	if strings.TrimSpace(override) != "" {
		if len(override) < 32 {
			return "", errors.New("AUTH_SECRET must be at least 32 characters")
		}
		return override, nil
	}
	path, err := secretFile()
	if err != nil {
		return "", err
	}
	if data, err := os.ReadFile(path); err == nil && strings.TrimSpace(string(data)) != "" {
		return strings.TrimSpace(string(data)), nil
	}
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	secret := hex.EncodeToString(buf)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(secret+"\n"), 0o600); err != nil {
		return "", err
	}
	return secret, nil
}

func secretFile() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "go-react-boilerplate", "secret"), nil
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func envBool(key string) bool {
	raw := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	return raw == "1" || raw == "true"
}

func splitList(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if value := strings.TrimSpace(part); value != "" {
			out = append(out, strings.TrimRight(value, "/"))
		}
	}
	return out
}
