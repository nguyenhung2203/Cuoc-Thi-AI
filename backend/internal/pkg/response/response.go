package response

import (
	"encoding/json"
	"net/http"

	"backend/internal/pkg/errors"
)

// Meta holds pagination metadata for list responses.
type Meta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ErrorBody is the error section of a response envelope.
type ErrorBody struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// Envelope is the standard API response wrapper.
type Envelope struct {
	Success   bool       `json:"success"`
	Data      any        `json:"data,omitempty"`
	Meta      *Meta      `json:"meta,omitempty"`
	RequestID string     `json:"request_id,omitempty"`
	Error     *ErrorBody `json:"error,omitempty"`
}

// JSON writes a successful response envelope with the given status code.
func JSON(w http.ResponseWriter, status int, data any, meta *Meta, requestID string) {
	env := Envelope{
		Success:   true,
		Data:      data,
		Meta:      meta,
		RequestID: requestID,
	}
	write(w, status, env)
}

// Error writes an error response envelope derived from an AppError.
func Error(w http.ResponseWriter, appErr *errors.AppError, requestID string) {
	env := Envelope{
		Success:   false,
		RequestID: requestID,
		Error: &ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
	}
	write(w, appErr.HTTPStatus, env)
}

// write is the shared JSON serialiser; it sets Content-Type before writing the body.
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
