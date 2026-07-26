package realtime

import (
	"os"
	"strings"
	"time"

	"backend/internal/livekit"
)

// LiveKitConfig aliases the shared LiveKit config so existing gateway code
// (cmd/realtime, server.go) keeps compiling unchanged.
type LiveKitConfig = livekit.Config

// LoadLiveKitConfig delegates to the shared livekit.Load — one source of
// truth for credentials and the production fail-fast guard.
func LoadLiveKitConfig() (LiveKitConfig, error) {
	return livekit.Load()
}

// allowDevBroadcast is the pure decision for dev-only chat mirroring across
// rooms: it requires an explicit opt-in flag AND a non-production APP_ENV.
// Production always returns false regardless of the flag.
func allowDevBroadcast(appEnv, flagValue string) bool {
	env := strings.ToLower(strings.TrimSpace(appEnv))
	if env == "production" || env == "prod" {
		return false
	}
	return strings.ToLower(strings.TrimSpace(flagValue)) == "true"
}

// devBroadcastAllRooms reads APP_ENV + REALTIME_DEV_BROADCAST_ALL from the
// environment. Default OFF — chat never leaves its own room unless a dev
// explicitly opts in. Reads env on every call so t.Setenv works in tests.
func devBroadcastAllRooms() bool {
	return allowDevBroadcast(os.Getenv("APP_ENV"), os.Getenv("REALTIME_DEV_BROADCAST_ALL"))
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
