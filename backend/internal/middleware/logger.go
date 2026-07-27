package middleware

import (
	"net/http"
	"time"

	"backend/internal/pkg/logger"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a custom ResponseWriter to capture the status code
		rw := &responseWriter{w, http.StatusOK}
		
		next.ServeHTTP(rw, r)
		
		duration := time.Since(start)
		logger.Info("[%s] %s - %d %s", r.Method, r.URL.Path, rw.status, duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
