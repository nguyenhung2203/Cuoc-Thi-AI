package errors

import (
	"errors"
	"net/http"
)

// Error code constants
const (
	UNAUTHORIZED     = "UNAUTHORIZED"
	FORBIDDEN        = "FORBIDDEN"
	NOT_FOUND        = "NOT_FOUND"
	VALIDATION_ERROR = "VALIDATION_ERROR"
	CONFLICT         = "CONFLICT"
	RATE_LIMITED     = "RATE_LIMITED"
	AI_SERVICE_ERROR = "AI_SERVICE_ERROR"
	REALTIME_ERROR   = "REALTIME_ERROR"
	INTERNAL_ERROR   = "INTERNAL_ERROR"
	BAD_REQUEST      = "BAD_REQUEST"
	// ACCOUNT_BLOCKED must stay in sync with the frontend check in
	// frontend/src/services/api.service.js (redirect to /403?reason=account_blocked).
	ACCOUNT_BLOCKED   = "ACCOUNT_BLOCKED"
	PAYLOAD_TOO_LARGE = "PAYLOAD_TOO_LARGE"
)

// AppError is a structured application error with an HTTP status code.
type AppError struct {
	Code       string
	Message    string
	Details    []string
	HTTPStatus int
}

// Error implements the error interface.
func (e *AppError) Error() string {
	return e.Message
}

// IsAppError checks whether err is (or wraps) an *AppError and returns it.
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// NewUnauthorized returns a 401 AppError.
func NewUnauthorized(msg string) *AppError {
	return &AppError{
		Code:       UNAUTHORIZED,
		Message:    msg,
		HTTPStatus: http.StatusUnauthorized,
	}
}

// NewForbidden returns a 403 AppError.
func NewForbidden(msg string) *AppError {
	return &AppError{
		Code:       FORBIDDEN,
		Message:    msg,
		HTTPStatus: http.StatusForbidden,
	}
}

// NewNotFound returns a 404 AppError.
func NewNotFound(msg string) *AppError {
	return &AppError{
		Code:       NOT_FOUND,
		Message:    msg,
		HTTPStatus: http.StatusNotFound,
	}
}

// NewValidation returns a 422 AppError with field-level details.
func NewValidation(msg string, details []string) *AppError {
	return &AppError{
		Code:       VALIDATION_ERROR,
		Message:    msg,
		Details:    details,
		HTTPStatus: http.StatusUnprocessableEntity,
	}
}

// NewConflict returns a 409 AppError.
func NewConflict(msg string) *AppError {
	return &AppError{
		Code:       CONFLICT,
		Message:    msg,
		HTTPStatus: http.StatusConflict,
	}
}

// NewInternal returns a 500 AppError.
func NewInternal(msg string) *AppError {
	return &AppError{
		Code:       INTERNAL_ERROR,
		Message:    msg,
		HTTPStatus: http.StatusInternalServerError,
	}
}

// NewBadRequest returns a 400 AppError.
func NewBadRequest(msg string) *AppError {
	return &AppError{
		Code:       BAD_REQUEST,
		Message:    msg,
		HTTPStatus: http.StatusBadRequest,
	}
}

// NewRateLimited returns a 429 AppError.
func NewRateLimited(msg string) *AppError {
	return &AppError{
		Code:       RATE_LIMITED,
		Message:    msg,
		HTTPStatus: http.StatusTooManyRequests,
	}
}

// NewAccountBlocked returns a 403 AppError for blocked user accounts.
func NewAccountBlocked(msg string) *AppError {
	return &AppError{
		Code:       ACCOUNT_BLOCKED,
		Message:    msg,
		HTTPStatus: http.StatusForbidden,
	}
}

// NewPayloadTooLarge returns a 413 AppError for oversized request bodies.
func NewPayloadTooLarge(msg string) *AppError {
	return &AppError{
		Code:       PAYLOAD_TOO_LARGE,
		Message:    msg,
		HTTPStatus: http.StatusRequestEntityTooLarge,
	}
}

// NewAIServiceError returns a 502 AppError for upstream AI failures.
func NewAIServiceError(msg string) *AppError {
	return &AppError{
		Code:       AI_SERVICE_ERROR,
		Message:    msg,
		HTTPStatus: http.StatusBadGateway,
	}
}

