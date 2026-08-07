package config

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration for the backend service.
type Config struct {
	AppPort string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// SMTP
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	SMTPFrom string

	// Public base URL of the frontend (for links in emails)
	FrontendURL string

	// Public base URL of this backend (used to build downloadable file URLs served to browsers)
	PublicBaseURL string

	// Gemini
	GeminiAPIKey string

	// Google OAuth client ID (audience to validate Google id_token against)
	GoogleClientID string

	// JWT
	JWTSecret     string
	JWTAccessTTL  string
	JWTRefreshTTL string
	CookieSecure  bool

	// AI Service
	AIServiceURL string

	// Redis / async queue
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string

	// Object storage (S3-compatible)
	StorageBucket    string
	StorageRegion    string
	StorageEndpoint  string
	StorageAccessKey string
	StorageSecretKey string

	// File uploads / signed download URLs
	MaxUploadBytes int64         // hard cap for multipart upload bodies
	UploadDir      string        // directory where uploaded files are stored
	FileURLSecret  []byte        // HMAC key for signing /uploads/ download URLs
	FileURLTTL     time.Duration // validity window of a signed download URL
}

// Load reads .env (if present) and then populates Config from environment
// variables, applying sensible defaults where values are absent.
func Load() (*Config, error) {
	// Best-effort: ignore the error so the app works when .env doesn't exist
	// (e.g. in containerised environments that inject vars directly).
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", ""),
		DBSSLMode:     getEnv("DB_SSL_MODE", "disable"),
		SMTPHost:      getEnv("SMTP_HOST", ""),
		SMTPPort:      getEnv("SMTP_PORT", "587"),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPass:      getEnv("SMTP_PASS", ""),
		SMTPFrom:      getEnv("SMTP_FROM", "no-reply@ai-interview.local"),
		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:5173"),
		PublicBaseURL: getEnv("PUBLIC_BASE_URL", "http://localhost:18080"),
		GeminiAPIKey:   getEnv("GEMINI_API_KEY", ""),
		GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTAccessTTL:  getEnv("JWT_ACCESS_TTL", "15m"),
		JWTRefreshTTL: getEnv("JWT_REFRESH_TTL", "168h"), // 7 days
		CookieSecure:  getEnvBool("COOKIE_SECURE", false),

		AIServiceURL: getEnv("AI_SERVICE_URL", "http://localhost:8000"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnv("REDIS_DB", "0"),

		StorageBucket:    getEnv("STORAGE_BUCKET", ""),
		StorageRegion:    getEnv("STORAGE_REGION", ""),
		StorageEndpoint:  getEnv("STORAGE_ENDPOINT", ""),
		StorageAccessKey: getEnv("STORAGE_ACCESS_KEY", ""),
		StorageSecretKey: getEnv("STORAGE_SECRET_KEY", ""),

		MaxUploadBytes: parseUploadMB(getEnv("MAX_UPLOAD_MB", "10")),
		UploadDir:      getEnv("UPLOAD_DIR", "uploads"),
		FileURLTTL:     parseTTL(getEnv("FILE_URL_TTL", "15m"), 15*time.Minute),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("config: JWT_SECRET must be set")
	}

	cfg.FileURLSecret = deriveFileURLSecret(getEnv("FILE_URL_SECRET", ""), cfg.JWTSecret)

	return cfg, nil
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

// fallback of 10 MB on malformed input.
func parseUploadMB(s string) int64 {
	mb, err := strconv.Atoi(s)
	if err != nil || mb < 1 || mb > 100 {
		mb = 10
	}
	return int64(mb) << 20
}

// parseTTL parses a duration string, returning fallback on malformed input.
func parseTTL(s string, fallback time.Duration) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return fallback
	}
	return d
}

// deriveFileURLSecret returns the HMAC key for signed download URLs. When
// FILE_URL_SECRET is unset it derives a stable key from JWT_SECRET with domain
// separation, so file-URL MACs can never be replayed as JWT MACs and URLs
// survive process restarts.
func deriveFileURLSecret(explicit, jwtSecret string) []byte {
	if explicit != "" {
		return []byte(explicit)
	}
	sum := sha256.Sum256([]byte("file-url-signing-v1|" + jwtSecret))
	return sum[:]
}

// getEnv returns the environment variable named by key, or fallback if unset/empty.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
