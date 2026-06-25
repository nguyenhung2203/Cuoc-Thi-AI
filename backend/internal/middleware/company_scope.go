package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	apierrors "backend/internal/pkg/errors"
	"backend/internal/pkg/response"
)

// CompanyScopeMiddleware verifies the authenticated user belongs to the company
// identified by the {company_id} URL parameter.
//
// Prerequisites:
//   - AuthMiddleware must run first (CtxUserID and CtxUserRole must be set).
//
// Behaviour:
//   - If the user's role is "admin", the membership check is skipped and
//     company_id is set in context directly.
//   - Otherwise, a direct DB query against company_members confirms active membership.
//   - On success, CtxCompanyID is set in the request context.
//   - On failure, a 403 JSON response is written and the chain is stopped.
func CompanyScopeMiddleware(db *sqlx.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := getRequestIDFromCtx(r)

			companyID := chi.URLParam(r, "company_id")
			if companyID == "" {
				response.Error(w, apierrors.NewForbidden("company_id is required"), requestID)
				return
			}

			userID, ok := r.Context().Value(CtxUserID).(string)
			if !ok || userID == "" {
				response.Error(w, apierrors.NewUnauthorized("unauthenticated"), requestID)
				return
			}

			userRole, _ := r.Context().Value(CtxUserRole).(string)

			ctx := r.Context()
			// Admins bypass membership check — they have platform-wide access.
			if userRole == "admin" {
				ctx = context.WithValue(ctx, CtxCompanyRole, "admin")
			} else {
				const q = `
					SELECT role FROM company_members
					WHERE company_id = $1::uuid
					  AND user_id = $2::uuid
					  AND status = 'active'
				`

				var role string
				err := db.QueryRowContext(r.Context(), q, companyID, userID).Scan(&role)
				if err != nil {
					response.Error(w, apierrors.NewForbidden("you are not a member of this company"), requestID)
					return
				}
				ctx = context.WithValue(r.Context(), CtxCompanyRole, role)
			}

			ctx = context.WithValue(ctx, CtxCompanyID, companyID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// getRequestIDFromCtx reads the request_id string from context (set by AuthMiddleware).
// Returns an empty string if not present so callers never receive a nil interface.
func getRequestIDFromCtx(r *http.Request) string {
	id, _ := r.Context().Value(CtxRequestID).(string)
	return id
}
