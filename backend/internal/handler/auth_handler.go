package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/dto/request"
	"backend/internal/dto/response"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/pkg/errors"
	pkgresponse "backend/internal/pkg/response"
	"backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	jwtSecret   string
	auditSvc    *service.AuditService
}

func NewAuthHandler(authService *service.AuthService, jwtSecret string, auditSvc *service.AuditService) *AuthHandler {
	return &AuthHandler{authService: authService, jwtSecret: jwtSecret, auditSvc: auditSvc}
}

func (h *AuthHandler) Routes(r chi.Router) {
	r.Post("/register", h.Register)
	r.Post("/register/send-otp", h.SendRegistrationOTP)
	r.Post("/login", h.Login)
	r.Post("/google-login", h.GoogleLogin)
	r.Post("/google", h.GoogleLogin)
	r.Post("/refresh", h.Refresh)
	r.Post("/forgot-password", h.ForgotPassword)
	r.Post("/verify-reset-otp", h.VerifyResetOTP)
	r.Post("/reset-password", h.ResetPassword)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Post("/logout", h.Logout)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Post("/logout-all", h.LogoutAll)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Get("/me", h.Me)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Put("/me", h.UpdateMe)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Get("/me/settings", h.GetSettings)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Put("/me/settings", h.UpdateSettings)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Put("/me/password", h.ChangePassword)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Post("/me/verify-document", h.VerifyDocument)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Delete("/me", h.DeleteAccount)
}

func setRefreshTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/v1/auth",
		MaxAge:   604800, // 7 days
	})
}

func clearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/v1/auth",
		MaxAge:   -1,
	})
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req request.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}

	ipAddress := getClientIP(r)
	userAgent := r.UserAgent()

	authResp, refreshToken, err := h.authService.Register(r.Context(), req, ipAddress, userAgent)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("registration failed"), "")
		}
		return
	}

	setRefreshTokenCookie(w, refreshToken)
	pkgresponse.JSON(w, http.StatusCreated, authResp, nil, "")
}

// SendRegistrationOTP emails a one-time code to verify a new registration email.
func (h *AuthHandler) SendRegistrationOTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		pkgresponse.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}
	if err := h.authService.SendRegistrationOTP(r.Context(), body.Email); err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to send otp"), "")
		}
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"status": "otp_sent"}, nil, "")
}

// ForgotPassword emails a password-reset OTP. Always returns 200 to avoid account enumeration.
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		pkgresponse.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}
	if err := h.authService.ForgotPassword(r.Context(), body.Email); err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to process request"), "")
		}
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"status": "otp_sent"}, nil, "")
}

// VerifyResetOTP validates a password-reset code without consuming it.
func (h *AuthHandler) VerifyResetOTP(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		pkgresponse.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}
	if err := h.authService.VerifyResetOTP(r.Context(), body.Email, body.OTP); err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to verify otp"), "")
		}
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"status": "otp_valid"}, nil, "")
}

// ResetPassword verifies the reset OTP and sets a new password.
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		pkgresponse.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}
	if err := h.authService.ResetPassword(r.Context(), body.Email, body.OTP, body.NewPassword); err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to reset password"), "")
		}
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"status": "password_reset"}, nil, "")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}

	ipAddress := getClientIP(r)
	userAgent := r.UserAgent()

	authResp, refreshToken, err := h.authService.Login(r.Context(), req, ipAddress, userAgent)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("login failed"), "")
		}
		return
	}

	setRefreshTokenCookie(w, refreshToken)

	// Audit log login
	if authResp != nil {
		h.auditSvc.LogAction(r.Context(), service.AuditLogInput{
			ActorUserID:  authResp.User.ID,
			ActorRole:    authResp.User.Role,
			Action:       "login",
			ResourceType: "user",
			ResourceID:   authResp.User.ID,
			IPAddress:    ipAddress,
			UserAgent:    userAgent,
		})
	}

	pkgresponse.JSON(w, http.StatusOK, authResp, nil, "")
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	var req request.GoogleLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, errors.NewBadRequest("invalid request body"), "")
		return
	}

	ipAddress := getClientIP(r)
	userAgent := r.UserAgent()

	authResp, refreshToken, err := h.authService.GoogleLogin(r.Context(), req, ipAddress, userAgent)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("google login failed"), "")
		}
		return
	}

	setRefreshTokenCookie(w, refreshToken)

	if authResp != nil {
		h.auditSvc.LogAction(r.Context(), service.AuditLogInput{
			ActorUserID:  authResp.User.ID,
			ActorRole:    authResp.User.Role,
			Action:       "google_login",
			ResourceType: "user",
			ResourceID:   authResp.User.ID,
			IPAddress:    ipAddress,
			UserAgent:    userAgent,
		})
	}

	pkgresponse.JSON(w, http.StatusOK, authResp, nil, "")
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		pkgresponse.Error(w, errors.NewUnauthorized("refresh token not found in cookies"), "")
		return
	}

	ipAddress := getClientIP(r)
	userAgent := r.UserAgent()

	tokens, newRefreshToken, err := h.authService.RefreshToken(r.Context(), cookie.Value, ipAddress, userAgent)
	if err != nil {
		clearRefreshTokenCookie(w)
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("refresh failed"), "")
		}
		return
	}

	setRefreshTokenCookie(w, newRefreshToken)
	pkgresponse.JSON(w, http.StatusOK, tokens, nil, "")
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		pkgresponse.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}

	me, err := h.authService.GetMe(r.Context(), userID)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to fetch profile"), "")
		}
		return
	}

	pkgresponse.JSON(w, http.StatusOK, me, nil, "")
}

// UpdateMe updates the current user's profile (full_name, avatar_url).
func (h *AuthHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		pkgresponse.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}

	var body struct {
		FullName  string `json:"full_name"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		pkgresponse.Error(w, errors.NewValidation("payload", []string{"invalid json payload"}), "")
		return
	}

	if err := h.authService.UpdateProfile(r.Context(), userID, body.FullName, body.AvatarURL); err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to update profile"), "")
		}
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"status": "updated"}, nil, "")
}

// GetSettings returns the current user's settings JSON.
func (h *AuthHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		pkgresponse.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}
	settings, err := h.authService.GetSettings(r.Context(), userID)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to fetch settings"), "")
		}
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]interface{}{"settings": settings}, nil, "")
}

// UpdateSettings replaces the current user's settings JSON.
func (h *AuthHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		pkgresponse.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}

	// Accept either {"settings": {...}} or a raw object body.
	var wrapper struct {
		Settings json.RawMessage `json:"settings"`
	}
	raw, _ := io.ReadAll(r.Body)
	var payload json.RawMessage
	if err := json.Unmarshal(raw, &wrapper); err == nil && len(wrapper.Settings) > 0 {
		payload = wrapper.Settings
	} else {
		payload = raw
	}
	if len(payload) == 0 {
		payload = json.RawMessage([]byte("{}"))
	}

	if err := h.authService.UpdateSettings(r.Context(), userID, models.JSONB(payload)); err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to save settings"), "")
		}
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"status": "saved"}, nil, "")
}

// ChangePassword lets the authenticated user change their own password.
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		pkgresponse.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}

	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		pkgresponse.Error(w, errors.NewValidation("payload", []string{"invalid json payload"}), "")
		return
	}

	if err := h.authService.ChangePassword(r.Context(), userID, body.CurrentPassword, body.NewPassword); err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to change password"), "")
		}
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"status": "password_changed"}, nil, "")
}

// VerifyDocument allows a pending user to submit their verification document.
func (h *AuthHandler) VerifyDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		pkgresponse.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}

	var body struct {
		FileID string `json:"file_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		pkgresponse.Error(w, errors.NewValidation("payload", []string{"invalid json payload"}), "")
		return
	}

	if err := h.authService.VerifyDocument(r.Context(), userID, body.FileID); err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to submit verification document"), "")
		}
		return
	}
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"status": "document_submitted"}, nil, "")
}

// DeleteAccount soft-deletes the authenticated user's account and clears their session.
func (h *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || userID == "" {
		pkgresponse.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}
	if err := h.authService.DeleteAccount(r.Context(), userID); err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to delete account"), "")
		}
		return
	}
	clearRefreshTokenCookie(w)
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"status": "account_deleted"}, nil, "")
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.CtxUserID).(string)
	userRole, _ := r.Context().Value(middleware.CtxUserRole).(string)

	cookie, err := r.Cookie("refresh_token")
	if err == nil {
		_ = h.authService.Logout(r.Context(), cookie.Value)
	}

	// Audit log logout
	if userID != "" {
		h.auditSvc.LogAction(r.Context(), service.AuditLogInput{
			ActorUserID:  userID,
			ActorRole:    userRole,
			Action:       "logout",
			ResourceType: "user",
			ResourceID:   userID,
			IPAddress:    getClientIP(r),
			UserAgent:    r.UserAgent(),
		})
	}

	clearRefreshTokenCookie(w)
	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"}, nil, "")
}

func (h *AuthHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok || uid == "" {
		pkgresponse.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}

	userRole, _ := r.Context().Value(middleware.CtxUserRole).(string)
	revokedCount, _ := h.authService.LogoutAll(r.Context(), uid)

	// Audit log logout-all
	h.auditSvc.LogAction(r.Context(), service.AuditLogInput{
		ActorUserID:  uid,
		ActorRole:    userRole,
		Action:       "logout_all",
		ResourceType: "user",
		ResourceID:   uid,
		IPAddress:    getClientIP(r),
		UserAgent:    r.UserAgent(),
		AfterData:    map[string]int{"revoked_sessions": revokedCount},
	})

	clearRefreshTokenCookie(w)
	pkgresponse.JSON(w, http.StatusOK, response.LogoutAllResponse{
		Message:         "All sessions revoked. Please login again.",
		RevokedSessions: revokedCount,
	}, nil, "")
}
