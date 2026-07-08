package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"backend/internal/middleware"
	"backend/internal/models"
	apierrors "backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/repository"
)

type NotificationHandler struct {
	notifRepo *repository.NotificationRepository
}

func NewNotificationHandler(notifRepo *repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{notifRepo: notifRepo}
}

func (h *NotificationHandler) ProtectedRoutes(r chi.Router) {
	r.Get("/", h.ListNotifications)
	r.Put("/read-all", h.MarkAllNotificationsRead)
	r.Put("/{notification_id}/read", h.MarkNotificationRead)
}

func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
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
	offset := (page - 1) * pageSize

	items, err := h.notifRepo.ListByUser(r.Context(), userID, pageSize, offset)
	if err != nil {
		pkgresponse.Error(w, apierrors.NewInternal(err.Error()), requestID)
		return
	}

	unreadCount := 0
	count, err := h.notifRepo.CountUnreadByUserID(r.Context(), userID)
	if err == nil {
		unreadCount = count
	}

	if items == nil {
		items = []models.Notification{}
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{
		"notifications": items,
		"unread_count":  unreadCount,
	}, nil, requestID)
}

func (h *NotificationHandler) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	idStr := chi.URLParam(r, "notification_id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		pkgresponse.Error(w, apierrors.NewValidation("notification_id", []string{"must be a valid integer"}), requestID)
		return
	}

	if err := h.notifRepo.MarkRead(r.Context(), id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			pkgresponse.Error(w, apierrors.NewNotFound("notification"), requestID)
			return
		}
		pkgresponse.Error(w, apierrors.NewInternal(err.Error()), requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{
		"message": "notification marked as read",
	}, nil, requestID)
}

func (h *NotificationHandler) MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	requestID, _ := r.Context().Value(middleware.CtxRequestID).(string)

	if err := h.notifRepo.MarkAllAsRead(r.Context(), userID); err != nil {
		pkgresponse.Error(w, apierrors.NewInternal(err.Error()), requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{
		"message": "all notifications marked as read",
	}, nil, requestID)
}
