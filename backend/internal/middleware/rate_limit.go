package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	apierrors "backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Minute), 100)
		visitors[ip] = limiter
	}
	return limiter
}

func AuthRateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // Simplified, should parse X-Forwarded-For in prod
		limiter := getVisitor(ip)
		if !limiter.Allow() {
			pkgresponse.Error(w, apierrors.NewRateLimited("Too many requests"), "")
			return
		}
		next.ServeHTTP(w, r)
	})
}
