package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/dto/request"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/response"
	"backend/internal/service"
	"backend/internal/middleware"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtSecret   string
}

func NewAuthHandler(authService *service.AuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{authService: authService, jwtSecret: jwtSecret}
}

func (h *AuthHandler) Routes(r chi.Router) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Post("/logout", h.Logout)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Get("/me", h.Me)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req request.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}

	tokens, err := h.authService.Register(r.Context(), req)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(w, appErr, "")
		} else {
			response.Error(w, errors.NewInternal("registration failed"), "")
		}
		return
	}

	response.JSON(w, http.StatusCreated, tokens, nil, "")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}

	tokens, err := h.authService.Login(r.Context(), req)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(w, appErr, "")
		} else {
			response.Error(w, errors.NewInternal("login failed"), "")
		}
		return
	}

	response.JSON(w, http.StatusOK, tokens, nil, "")
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req request.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}

	tokens, err := h.authService.RefreshToken(r.Context(), req)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(w, appErr, "")
		} else {
			response.Error(w, errors.NewInternal("refresh failed"), "")
		}
		return
	}

	response.JSON(w, http.StatusOK, tokens, nil, "")
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		response.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}

	me, err := h.authService.GetMe(r.Context(), userID)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			response.Error(w, appErr, "")
		} else {
			response.Error(w, errors.NewInternal("failed to fetch profile"), "")
		}
		return
	}

	response.JSON(w, http.StatusOK, me, nil, "")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Stub for JWT blocklist
	response.JSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"}, nil, "")
}
