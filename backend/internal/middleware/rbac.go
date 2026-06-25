package middleware

import (
	"net/http"

	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/response"
)

// RequireRole returns a middleware that allows only users whose role appears in
// the allowed list. The role is read from CtxUserRole, which must have been set
// by AuthMiddleware earlier in the chain.
//
// Usage:
//
//	r.With(middleware.RequireRole("admin", "recruiter")).Get("/jobs", handler)
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	// Build a set for O(1) lookup.
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := getRequestIDFromCtx(r)

			userRole, _ := r.Context().Value(CtxUserRole).(string)
			if _, ok := allowed[userRole]; !ok {
				response.Error(w, apierrors.NewForbidden("insufficient permissions"), requestID)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireCompanyRole restricts access based on company role (owner, admin, member, viewer).
// Must be used AFTER CompanyScopeMiddleware.
func RequireCompanyRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := getRequestIDFromCtx(r)
			role, _ := r.Context().Value(CtxCompanyRole).(string)

			// System admin always allowed
			sysRole, _ := r.Context().Value(CtxUserRole).(string)
			if sysRole == "admin" {
				next.ServeHTTP(w, r)
				return
			}

			if _, ok := allowed[role]; !ok {
				response.Error(w, apierrors.NewForbidden("insufficient company permissions"), requestID)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
