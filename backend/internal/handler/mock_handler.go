package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/middleware"
	apierrors "backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/service"
)

type MockHandler struct {
	svc *service.MockService
}

func NewMockHandler(svc *service.MockService) *MockHandler {
	return &MockHandler{svc: svc}
}

func (h *MockHandler) ProtectedRoutes(r chi.Router) {
	r.Post("/", h.Create)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.GetMockInterview)
		r.Get("/messages", h.GetMessages)
		r.Post("/messages", h.AddMessage)
		r.Post("/finish", h.Finish)
	})
}

type CreateMockReq struct {
	TargetRole  string `json:"target_role"`
	TargetLevel string `json:"target_level"`
	CVFileID    string `json:"cv_file_id"`
}

func (h *MockHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)

	var req CreateMockReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeServiceError(w, apierrors.NewBadRequest("invalid JSON"), requestID)
		return
	}

	mi, err := h.svc.CreateMockInterview(r.Context(), userID, req.TargetRole, req.TargetLevel, req.CVFileID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusCreated, mi, nil, requestID)
}

func (h *MockHandler) GetMockInterview(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	mockID := chi.URLParam(r, "id")

	mi, err := h.svc.GetMockInterview(r.Context(), mockID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, mi, nil, requestID)
}

func (h *MockHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	mockID := chi.URLParam(r, "id")

	msgs, err := h.svc.GetMessages(r.Context(), mockID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, msgs, nil, requestID)
}

type AddMessageReq struct {
	Content string `json:"content"`
}

func (h *MockHandler) AddMessage(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	mockID := chi.URLParam(r, "id")

	var req AddMessageReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeServiceError(w, apierrors.NewBadRequest("invalid JSON"), requestID)
		return
	}

	err := h.svc.ProcessMessage(r.Context(), mockID, req.Content)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	// Fetch updated messages to return
	msgs, _ := h.svc.GetMessages(r.Context(), mockID)
	pkgresponse.JSON(w, http.StatusOK, msgs, nil, requestID)
}

func (h *MockHandler) Finish(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)
	mockID := chi.URLParam(r, "id")

	err := h.svc.FinishMockInterview(r.Context(), mockID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"status": "completed"}, nil, requestID)
}
