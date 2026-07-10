package config

import (
	"fmt"
	"os"

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

	// Gemini
	GeminiAPIKey string

	// JWT
	JWTSecret     string
	JWTAccessTTL  string
	JWTRefreshTTL string

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
		GeminiAPIKey:  getEnv("GEMINI_API_KEY", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTAccessTTL:  getEnv("JWT_ACCESS_TTL", "15m"),
		JWTRefreshTTL: getEnv("JWT_REFRESH_TTL", "168h"), // 7 days

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
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("config: JWT_SECRET must be set")
	}

	return cfg, nil
}

// getEnv returns the environment variable named by key, or fallback if unset/empty.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
