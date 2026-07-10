package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"backend/internal/middleware"
	apierrors "backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/service"
)

type MockHandler struct {
	mockSvc *service.MockService
}

func NewMockHandler(mockSvc *service.MockService) *MockHandler {
	return &MockHandler{mockSvc: mockSvc}
}

func (h *MockHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	var req service.CreateMockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid json", []string{err.Error()}), requestID)
		return
	}
	if req.TargetRole == "" {
		pkgresponse.Error(w, apierrors.NewValidation("target_role", []string{"required"}), requestID)
		return
	}

	m, err := h.mockSvc.Create(r.Context(), userID, req)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusCreated, m, nil, requestID)
}

func (h *MockHandler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, err := h.mockSvc.ListByUser(r.Context(), userID, page, pageSize)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, items, nil, requestID)
}

func (h *MockHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)
	id := chi.URLParam(r, "id")

	m, err := h.mockSvc.GetByID(r.Context(), id, userID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, m, nil, requestID)
}

func (h *MockHandler) Start(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)
	id := chi.URLParam(r, "id")

	if err := h.mockSvc.Start(r.Context(), id, userID); err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "started"}, nil, requestID)
}

func (h *MockHandler) End(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)
	id := chi.URLParam(r, "id")

	if err := h.mockSvc.End(r.Context(), id, userID); err != nil {
		writeServiceError(w, err, requestID)
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "ended"}, nil, requestID)
}

func (h *MockHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)
	id := chi.URLParam(r, "id")

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid json", []string{err.Error()}), requestID)
		return
	}
	if req.Content == "" {
		pkgresponse.Error(w, apierrors.NewValidation("content", []string{"required"}), requestID)
		return
	}

	userMsg, aiMsg, err := h.mockSvc.SendMessage(r.Context(), id, userID, req.Content)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{
		"user_message": userMsg,
		"ai_message":   aiMsg,
	}, nil, requestID)
}

// SaveLiveTranscript persists a spoken (Gemini Live) mock conversation as
// mock_interview_messages so the report generator has material to work with.
func (h *MockHandler) SaveLiveTranscript(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)
	id := chi.URLParam(r, "id")

	var req struct {
		Turns []struct {
			Role string `json:"role"`
			Text string `json:"text"`
		} `json:"turns"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("invalid json", []string{err.Error()}), requestID)
		return
	}

	turns := make([]service.LiveTurn, 0, len(req.Turns))
	for _, t := range req.Turns {
		if t.Text == "" {
			continue
		}
		turns = append(turns, service.LiveTurn{Role: t.Role, Text: t.Text})
	}

	if err := h.mockSvc.SaveLiveTranscript(r.Context(), id, userID, turns); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{"saved": len(turns)}, nil, requestID)
}

func (h *MockHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)
	id := chi.URLParam(r, "id")

	report, err := h.mockSvc.GetReport(r.Context(), id, userID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, report, nil, requestID)
}

func (h *MockHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)
	id := chi.URLParam(r, "id")

	messages, err := h.mockSvc.GetMessages(r.Context(), id, userID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, messages, nil, requestID)
}
