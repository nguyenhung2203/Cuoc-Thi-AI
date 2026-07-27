package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	jwtpkg "backend/internal/pkg/jwt"
	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/response"
)

// AuthMiddleware extracts and validates the Bearer token from the Authorization header.
// On success it sets CtxUserID and CtxUserRole in the request context and attaches a
// request_id (from github.com/google/uuid) to both the context and response envelope.
// On failure it writes a 401 JSON response and stops the chain.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate a unique request ID for this request.
			requestID := uuid.NewString()

			// Attach request_id to the response header so clients can correlate.
			w.Header().Set("X-Request-ID", requestID)

			// Store request_id in context so handlers and downstream middleware can read it.
			ctx := context.WithValue(r.Context(), CtxRequestID, requestID)

			// Extract the Authorization header.
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, apierrors.NewUnauthorized("authorization header is required"), requestID)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				response.Error(w, apierrors.NewUnauthorized("authorization header must be 'Bearer <token>'"), requestID)
				return
			}
			token := strings.TrimSpace(parts[1])
			if token == "" {
				response.Error(w, apierrors.NewUnauthorized("bearer token must not be empty"), requestID)
				return
			}

			// Validate the JWT.
			claims, err := jwtpkg.ValidateToken(token, jwtSecret)
			if err != nil {
				response.Error(w, apierrors.NewUnauthorized("invalid or expired token"), requestID)
				return
			}

			// Inject identity into context.
			ctx = context.WithValue(ctx, CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, CtxUserRole, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RoleMiddleware restricts access to the specified roles.
func RoleMiddleware(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID, _ := r.Context().Value(CtxRequestID).(string)
			userRole, ok := r.Context().Value(CtxUserRole).(string)
			if !ok {
				response.Error(w, apierrors.NewForbidden("missing user role in context"), requestID)
				return
			}

			allowed := false
			for _, role := range roles {
				if userRole == role {
					allowed = true
					break
				}
			}

			if !allowed {
				response.Error(w, apierrors.NewForbidden("insufficient permissions to access this resource"), requestID)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
