package livekit

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// Config holds the LiveKit credentials and endpoint used to mint and verify
// room tokens. It is the single source of truth for both the API server and
// the realtime gateway — do not read LIVEKIT_* env vars anywhere else.
type Config struct {
	APIKey    string
	APISecret string
	URL       string
}

// Dev fallback values, only acceptable outside production. They keep local
// dev working without real LiveKit credentials.
const (
	DevKey    = "devkey"
	DevSecret = "devsecret"
	DevURL    = "wss://livekit.example.com"
)

// Load reads LiveKit settings from the environment. When APP_ENV is
// "production"/"prod" it requires real, non-fallback credentials and returns
// an error otherwise, so binaries fail fast instead of minting tokens signed
// with a well-known dev secret.
func Load() (Config, error) {
	cfg := Config{
		APIKey:    os.Getenv("LIVEKIT_API_KEY"),
		APISecret: os.Getenv("LIVEKIT_API_SECRET"),
		URL:       os.Getenv("LIVEKIT_URL"),
	}

	appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	isProd := appEnv == "production" || appEnv == "prod"

	if isProd {
		switch {
		case cfg.APIKey == "" || cfg.APIKey == DevKey:
			return cfg, fmt.Errorf("LIVEKIT_API_KEY must be set to a real value in production")
		case cfg.APISecret == "" || cfg.APISecret == DevSecret:
			return cfg, fmt.Errorf("LIVEKIT_API_SECRET must be set to a real value in production")
		case cfg.URL == "" || cfg.URL == DevURL:
			return cfg, fmt.Errorf("LIVEKIT_URL must be set to a real value in production")
		}
		return cfg, nil
	}

	// Development: fill blanks with dev fallbacks so local runs work.
	if cfg.APIKey == "" {
		cfg.APIKey = DevKey
	}
	if cfg.APISecret == "" {
		cfg.APISecret = DevSecret
	}
	if cfg.URL == "" {
		cfg.URL = DevURL
	}
	return cfg, nil
}

var (
	currentOnce sync.Once
	currentCfg  Config
)

// Current returns the process-wide LiveKit config, resolved once. Both
// binaries call Load() at startup and exit on error, so by the time Current()
// is used the production guard has already run; the error here is ignored on
// purpose (dev fallbacks are filled in by Load).
func Current() Config {
	currentOnce.Do(func() {
		currentCfg, _ = Load()
	})
	return currentCfg
}
