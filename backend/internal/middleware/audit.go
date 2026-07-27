package middleware

import (
	"context"
	"net/http"

	"backend/internal/service"
)

type ctxKeyAudit string

const CtxAuditLog = ctxKeyAudit("audit_log")

// AuditMiddleware captures request metadata and exposes an AuditHelper via context.
// Handlers call AuditHelper.Log(action, resourceType, resourceID, before, after).
func AuditMiddleware(auditSvc *service.AuditService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			helper := &AuditHelper{
				svc:     auditSvc,
				userID:  "",
				role:    "",
				ip:      extractIP(r),
				ua:      r.UserAgent(),
				request: r,
			}
			// Read user info from context if available (set by AuthMiddleware)
			if uid, ok := r.Context().Value(CtxUserID).(string); ok {
				helper.userID = uid
			}
			if role, ok := r.Context().Value(CtxUserRole).(string); ok {
				helper.role = role
			}

			ctx := context.WithValue(r.Context(), CtxAuditLog, helper)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuditHelper provides a convenience method for handlers to log audit events.
type AuditHelper struct {
	svc     *service.AuditService
	userID  string
	role    string
	ip      string
	ua      string
	request *http.Request
}

func (a *AuditHelper) Log(action, resourceType, resourceID, companyID string, before, after interface{}) {
	go a.svc.LogAction(a.request.Context(), service.AuditLogInput{
		CompanyID:    companyID,
		ActorUserID:  a.userID,
		ActorRole:    a.role,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		BeforeData:   before,
		AfterData:    after,
		IPAddress:    a.ip,
		UserAgent:    a.ua,
	})
}

// GetAuditHelper retrieves the AuditHelper from request context.
// Returns nil if AuditMiddleware was not applied to this route.
func GetAuditHelper(r *http.Request) *AuditHelper {
	h, _ := r.Context().Value(CtxAuditLog).(*AuditHelper)
	return h
}

func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return r.RemoteAddr
}
