package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/cors"
)

func CorsMiddleware() func(http.Handler) http.Handler {
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		// KHÔNG dùng "*" khi AllowCredentials=true (vi phạm CORS spec, trình duyệt sẽ block Cookie)
		allowedOrigins = "http://localhost:3000,http://localhost:5173,http://localhost:13000"
	}
	origins := strings.Split(allowedOrigins, ",")

	return cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
