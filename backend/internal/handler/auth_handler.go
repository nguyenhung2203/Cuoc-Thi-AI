package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"backend/internal/dto/request"
	"backend/internal/dto/response"
	"backend/internal/middleware"
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
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.Post("/forgot-password", h.ForgotPassword)
	r.Post("/reset-password", h.ResetPassword)
	r.Post("/verify-email", h.VerifyEmail)
	r.Post("/resend-otp", h.ResendOTP)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Post("/logout", h.Logout)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Post("/logout-all", h.LogoutAll)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Get("/me", h.Me)
	r.With(middleware.AuthMiddleware(h.jwtSecret)).Post("/change-password", h.ChangePassword)
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

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.CtxUserID).(string)
	if !ok {
		pkgresponse.Error(w, errors.NewUnauthorized("unauthorized"), "")
		return
	}

	var req request.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, errors.NewValidation("invalid request format", nil), "")
		return
	}

	if err := h.authService.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			pkgresponse.Error(w, appErr, "")
		} else {
			pkgresponse.Error(w, errors.NewInternal("failed to change password"), "")
		}
		return
	}

	ipAddress := getClientIP(r)
	userAgent := r.UserAgent()

	// Audit log
	h.auditSvc.LogAction(r.Context(), service.AuditLogInput{
		ActorUserID:  userID,
		Action:       "change_password",
		ResourceType: "user",
		ResourceID:   userID,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
	})

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "Password changed successfully"}, nil, "")
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

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)

	var req request.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, errors.NewValidation("invalid json body", []string{err.Error()}), requestID)
		return
	}

	// Basic validation could be done here (e.g. validator.Validate(&req))

	if err := h.authService.ForgotPassword(r.Context(), req.Email); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "If an account with this email exists, a password reset link has been sent."}, nil, requestID)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)

	var req request.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, errors.NewValidation("invalid json body", []string{err.Error()}), requestID)
		return
	}

	if err := h.authService.ResetPassword(r.Context(), req.Email, req.OTP, req.NewPassword); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "Password has been successfully reset. You can now log in."}, nil, requestID)
}

func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)

	var req request.VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, errors.NewValidation("invalid json body", []string{err.Error()}), requestID)
		return
	}

	if err := h.authService.VerifyEmail(r.Context(), req.Email, req.OTP); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "Email has been successfully verified."}, nil, requestID)
}

func (h *AuthHandler) ResendOTP(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r)

	var req struct {
		Email   string `json:"email" validate:"required,email"`
		Purpose string `json:"purpose" validate:"required"` // 'register' or 'forgot_password'
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, errors.NewValidation("invalid json body", []string{err.Error()}), requestID)
		return
	}

	if err := h.authService.ResendOTP(r.Context(), req.Email, req.Purpose); err != nil {
		writeServiceError(w, err, requestID)
		return
	}

	pkgresponse.JSON(w, http.StatusOK, map[string]string{"message": "OTP has been resent successfully."}, nil, requestID)
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
