package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	apierrors "backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"golang.org/x/time/rate"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

const (
	defaultAuthRate  = 10
	loginRate        = 30
	defaultAuthBurst = 10
	loginBurst       = 30
)

func getVisitor(ip string, requestsPerMinute, burst int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	key := fmt.Sprintf("%s:%d:%d", ip, requestsPerMinute, burst)
	limiter, exists := visitors[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(requestsPerMinute)), burst)
		visitors[key] = limiter
	}
	return limiter
}

func rateLimitMiddleware(requestsPerMinute, burst int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := getVisitor(r.RemoteAddr, requestsPerMinute, burst)
			if !limiter.Allow() {
				pkgresponse.Error(w, apierrors.NewRateLimited("Too many requests"), "")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AuthRateLimitMiddleware(next http.Handler) http.Handler {
	return rateLimitMiddleware(defaultAuthRate, defaultAuthBurst)(next)
}

func LoginRateLimitMiddleware(next http.Handler) http.Handler {
	return rateLimitMiddleware(loginRate, loginBurst)(next)
}
