package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/middleware"
	"backend/internal/pkg/response"
	"backend/internal/service"
)

type NotificationHandler struct {
	svc *service.NotificationService
}

func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) ProtectedRoutes(r chi.Router) {
	r.Get("/", h.ListNotifications)
	r.Put("/read-all", h.MarkAllAsRead)
	r.Put("/{notification_id}/read", h.MarkAsRead)
}

func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok {
		userID = "" // should not happen if under ProtectedRoutes, but just in case
	}
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	items, err := h.svc.ListByUser(r.Context(), userID, 50, 0)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	response.JSON(w, http.StatusOK, items, nil, requestID)
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	notificationID := chi.URLParam(r, "notification_id")
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	err := h.svc.MarkAsRead(r.Context(), notificationID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "marked as read"}, nil, requestID)
}

func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok {
		userID = ""
	}
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	err := h.svc.MarkAllAsRead(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "all marked as read"}, nil, requestID)
}
