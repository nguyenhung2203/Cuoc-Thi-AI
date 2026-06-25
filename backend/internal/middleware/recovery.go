package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"backend/internal/pkg/errors"
	"backend/internal/pkg/logger"
	"backend/internal/pkg/response"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered: %v\n%s", err, debug.Stack())
				
				requestID := ""
				if reqID, ok := r.Context().Value(CtxRequestID).(string); ok {
					requestID = reqID
				}
				
				appErr := errors.NewInternal("an unexpected error occurred")
				log.Printf("PANIC recovered: %v", err)
				response.Error(w, appErr, requestID)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
