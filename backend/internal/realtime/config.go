package realtime

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// LiveKitConfig holds the LiveKit credentials and endpoint used to mint room
// tokens. Loaded from env; in production the dev fallbacks are refused.
type LiveKitConfig struct {
	APIKey    string
	APISecret string
	URL       string
}

// devLiveKitFallback values are only acceptable outside production. They keep
// local dev working without real LiveKit credentials.
const (
	devLiveKitKey    = "devkey"
	devLiveKitSecret = "devsecret"
	devLiveKitURL    = "wss://livekit.example.com"
)

// LoadLiveKitConfig reads LiveKit settings from the environment. When APP_ENV is
// "production"/"prod" it requires real, non-fallback credentials and returns an
// error otherwise, so the gateway fails fast instead of minting tokens signed
// with a well-known dev secret.
func LoadLiveKitConfig() (LiveKitConfig, error) {
	cfg := LiveKitConfig{
		APIKey:    os.Getenv("LIVEKIT_API_KEY"),
		APISecret: os.Getenv("LIVEKIT_API_SECRET"),
		URL:       os.Getenv("LIVEKIT_URL"),
	}

	appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	isProd := appEnv == "production" || appEnv == "prod"

	if isProd {
		switch {
		case cfg.APIKey == "" || cfg.APIKey == devLiveKitKey:
			return cfg, fmt.Errorf("LIVEKIT_API_KEY must be set to a real value in production")
		case cfg.APISecret == "" || cfg.APISecret == devLiveKitSecret:
			return cfg, fmt.Errorf("LIVEKIT_API_SECRET must be set to a real value in production")
		case cfg.URL == "" || cfg.URL == devLiveKitURL:
			return cfg, fmt.Errorf("LIVEKIT_URL must be set to a real value in production")
		}
		return cfg, nil
	}

	// Development: fill blanks with dev fallbacks so local runs work.
	if cfg.APIKey == "" {
		cfg.APIKey = devLiveKitKey
	}
	if cfg.APISecret == "" {
		cfg.APISecret = devLiveKitSecret
	}
	if cfg.URL == "" {
		cfg.URL = devLiveKitURL
	}
	return cfg, nil
}

const (
	// writeWait is the maximum time allowed to write a message to a client.
	writeWait = 10 * time.Second

	// pongWait is the maximum time to wait for a pong reply from the client.
	pongWait = 60 * time.Second

	// pingPeriod is how often the server sends a ping. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// maxMessageSize limits the size of incoming WebSocket messages (bytes).
	maxMessageSize = 64 * 1024 // 64 KB

	// sendBufferSize is the channel buffer per connection.
	sendBufferSize = 256

)

var (
	// gracePeriodShort is how long the server waits before marking a disconnected
	// participant as "reconnecting" (and NOT broadcasting room:user_left).
	gracePeriodShort = 2 * time.Minute

	// gracePeriodRoom is how long an empty room is kept alive after all participants leave.
	gracePeriodRoom = 5 * time.Minute

	// gracePeriodEndInterview is the cool-down after interview:end before the room closes.
	gracePeriodEndInterview = 30 * time.Second

	// heartbeatOfflineThreshold: no heartbeat for this long → "offline".
	heartbeatOfflineThreshold = gracePeriodShort
)

// SetGracePeriodsForTest overrides grace period durations for fast automated testing.
func SetGracePeriodsForTest(short, room, endInterview time.Duration) {
	gracePeriodShort = short
	gracePeriodRoom = room
	gracePeriodEndInterview = endInterview
	heartbeatOfflineThreshold = short
}

const (
	// roomWaitingExpiry is how long a room stays in "waiting" before it auto-expires.
	roomWaitingExpiry = 30 * time.Minute

	// heartbeatReconnectingThreshold: no heartbeat for this long → "reconnecting".
	heartbeatReconnectingThreshold = 60 * time.Second

	// rateLimit_aiSuggestionPerInterval is the max AI suggestion requests per window.
	rateLimitAISuggestionPerInterval = 10

	// rateLimit_aiSuggestionInterval is the sliding-window duration for the AI rate limit.
	rateLimitAISuggestionInterval = 10 * time.Minute

	// rateLimit_chatPerMinute is the max chat messages a single participant can send per minute.
	rateLimitChatPerMinute = 10

	// rateLimit_mediaStatusPerMinute caps mic/camera toggle spam.
	rateLimitMediaStatusPerMinute = 30

	// aiRetryMax is the maximum number of automatic AI retries before escalating to "critical".
	aiRetryMax = 3
)
