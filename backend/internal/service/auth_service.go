package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"backend/internal/dto/request"
	"backend/internal/dto/response"
	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/jwt"
	"backend/internal/queue"
	"backend/internal/repository"
)

type Mailer interface {
	Send(to, subject, body string) error
	SendHTML(to, subject, textBody, htmlBody string) error
}

type OTPEmailEnqueuer interface {
	EnqueueSendOTPEmail(payload queue.SendOTPEmailPayload) error
}

type AuthService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtSecret        string
	otpSvc           *OTPService
	mailer           Mailer
	otpEnqueuer      OTPEmailEnqueuer
	googleClientID   string
}

func NewAuthService(userRepo *repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        jwtSecret,
	}
}

// WithOTP wires the OTP service and mailer used for registration and password-reset flows.
func (s *AuthService) WithOTP(otpSvc *OTPService, mailer Mailer) *AuthService {
	s.otpSvc = otpSvc
	s.mailer = mailer
	return s
}

func (s *AuthService) WithOTPEnqueuer(enqueuer OTPEmailEnqueuer) *AuthService {
	s.otpEnqueuer = enqueuer
	return s
}

func (s *AuthService) WithGoogle(clientID string) *AuthService {
	s.googleClientID = clientID
	return s
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ensureNotBlocked rejects users whose account has been blocked by an admin.
// Only "blocked" is rejected: "pending" recruiters MUST be able to log in to
// upload their verification document (RecruiterLayout overlay), and "inactive"
// users are already filtered out by the deleted_at IS NULL clause in the repo.
//
// Callers must check the returned pointer explicitly:
//
//	if e := ensureNotBlocked(user); e != nil { return nil, "", e }
//
// Never return the result directly through an error interface — a typed-nil
// *AppError wrapped in error is non-nil.
func ensureNotBlocked(u *models.User) *errors.AppError {
	if u != nil && u.Status == models.UserStatusBlocked {
		return errors.NewAccountBlocked("tài khoản của bạn đã bị khoá, vui lòng liên hệ quản trị viên")
	}
	return nil
}

// stripPort loại bỏ phần :port khỏi RemoteAddr (ví dụ "192.168.1.1:12345" → "192.168.1.1")
func stripPort(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr // fallback: trả nguyên nếu không parse được
	}
	return host
}

// createAndSaveRefreshToken sinh refresh token, hash SHA-256, lưu vào DB và trả token raw
func (s *AuthService) createAndSaveRefreshToken(
	ctx context.Context,
	userID, email, role, familyID, ipAddress, userAgent string,
) (string, error) {
	_, refreshToken, err := jwt.GenerateTokenPair(
		userID, email, role, s.jwtSecret,
		15*time.Minute, 7*24*time.Hour,
	)
	if err != nil {
		return "", errors.NewInternal("failed to generate tokens")
	}

	tokenModel := &models.RefreshToken{
		UserID:    userID,
		FamilyID:  familyID,
		TokenHash: hashToken(refreshToken),
		IsRevoked: false,
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.refreshTokenRepo.Create(ctx, tokenModel); err != nil {
		return "", errors.NewInternal("failed to save refresh token")
	}

	return refreshToken, nil
}

func (s *AuthService) Register(ctx context.Context, req request.RegisterRequest, ipAddress, userAgent string) (*response.AuthResponse, string, error) {
	ipAddress = stripPort(ipAddress)

	// Require a verified email OTP before creating the account.
	if s.otpSvc != nil && s.otpSvc.Enabled() {
		ok, verifyErr := s.otpSvc.Verify(ctx, otpPurposeRegister, req.Email, req.OTP)
		if verifyErr != nil {
			return nil, "", errors.NewInternal("failed to verify otp")
		}
		if !ok {
			return nil, "", errors.NewValidation("otp", []string{"mã xác nhận không hợp lệ hoặc đã hết hạn"})
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", errors.NewInternal("failed to hash password")
	}

	// Self-registration may only create recruiter or candidate accounts.
	// admin (or any other value) is never accepted from client input — the
	// public /register endpoint must not be a path to privilege escalation.
	role := models.UserRole(req.Role)
	if role != models.RoleRecruiter && role != models.RoleCandidate {
		role = models.RoleCandidate // default
	}

	status := models.UserStatusActive
	if role == models.RoleRecruiter {
		status = models.UserStatusPending
	}

	user := &models.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		Role:         role,
		Status:       status,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", errors.NewConflict("email already exists")
	}

	accessToken, _, err := jwt.GenerateTokenPair(
		user.ID, user.Email, string(user.Role), s.jwtSecret,
		15*time.Minute, 7*24*time.Hour,
	)
	if err != nil {
		return nil, "", errors.NewInternal("failed to generate tokens")
	}

	familyID := uuid.NewString()
	refreshToken, err := s.createAndSaveRefreshToken(ctx, user.ID, user.Email, string(user.Role), familyID, ipAddress, userAgent)
	if err != nil {
		return nil, "", err
	}

	resp := &response.AuthResponse{
		User: response.AuthUser{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     string(user.Role),
			Status:   string(user.Status),
		},
		AccessToken: accessToken,
	}

	return resp, refreshToken, nil
}

func (s *AuthService) Login(ctx context.Context, req request.LoginRequest, ipAddress, userAgent string) (*response.AuthResponse, string, error) {
	ipAddress = stripPort(ipAddress)

	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", errors.NewUnauthorized("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, "", errors.NewUnauthorized("invalid email or password")
	}

	// Status check runs AFTER the bcrypt match so callers without valid
	// credentials cannot enumerate which accounts are blocked.
	if e := ensureNotBlocked(user); e != nil {
		return nil, "", e
	}

	accessToken, _, err := jwt.GenerateTokenPair(
		user.ID, user.Email, string(user.Role), s.jwtSecret,
		15*time.Minute, 7*24*time.Hour,
	)
	if err != nil {
		return nil, "", errors.NewInternal("failed to generate tokens")
	}

	familyID := uuid.NewString()
	refreshToken, err := s.createAndSaveRefreshToken(ctx, user.ID, user.Email, string(user.Role), familyID, ipAddress, userAgent)
	if err != nil {
		return nil, "", err
	}

	resp := &response.AuthResponse{
		User: response.AuthUser{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     string(user.Role),
			Status:   string(user.Status),
		},
		AccessToken: accessToken,
	}

	return resp, refreshToken, nil
}

// GoogleLogin handles real authentication/registration via Google OAuth API tokens
func (s *AuthService) GoogleLogin(ctx context.Context, req request.GoogleLoginRequest, ipAddress, userAgent string) (*response.AuthResponse, string, error) {
	ipAddress = stripPort(ipAddress)

	// Verify the Google id_token server-side. Never trust client-supplied email.
	info, err := verifyGoogleIDToken(ctx, req.IDToken, s.googleClientID)
	if err != nil {
		return nil, "", errors.NewUnauthorized("invalid Google credential")
	}
	// Use the verified identity from the token, not the request payload.
	email := info.Email
	fullName := info.Name
	if fullName == "" {
		fullName = req.FullName
	}
	avatar := info.Picture
	if avatar == "" {
		avatar = req.Avatar
	}

	// Look up user by verified email in PostgreSQL
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// User not found -> Auto create a real user record in database.
		// Only recruiter/candidate may be self-assigned; admin can never be
		// obtained through this public path (privilege escalation guard).
		role := models.UserRole(req.Role)
		if role != models.RoleRecruiter && role != models.RoleCandidate {
			role = models.RoleCandidate
		}
		status := models.UserStatusActive
		if role == models.RoleRecruiter {
			status = models.UserStatusPending
		}
		user = &models.User{
			ID:           uuid.NewString(),
			Email:        email,
			PasswordHash: "google_oauth_user",
			FullName:     fullName,
			Role:         role,
			Status:       status,
		}
		if avatar != "" {
			user.AvatarURL = sql.NullString{String: avatar, Valid: true}
		}
		if errCreate := s.userRepo.Create(ctx, user); errCreate != nil {
			return nil, "", errors.NewInternal("failed to create user from Google OAuth: " + errCreate.Error())
		}
	} else if user.FullName == "" && fullName != "" {
		_ = s.userRepo.UpdateProfile(ctx, user.ID, fullName, avatar)
	}

	// Google OAuth must not be a side door around an admin block.
	if e := ensureNotBlocked(user); e != nil {
		return nil, "", e
	}

	accessToken, _, err := jwt.GenerateTokenPair(
		user.ID, user.Email, string(user.Role), s.jwtSecret,
		15*time.Minute, 7*24*time.Hour,
	)
	if err != nil {
		return nil, "", errors.NewInternal("failed to generate tokens")
	}

	familyID := uuid.NewString()
	refreshToken, err := s.createAndSaveRefreshToken(ctx, user.ID, user.Email, string(user.Role), familyID, ipAddress, userAgent)
	if err != nil {
		return nil, "", err
	}

	resp := &response.AuthResponse{
		User: response.AuthUser{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     string(user.Role),
			Status:   string(user.Status),
		},
		AccessToken: accessToken,
	}

	return resp, refreshToken, nil
}

// UpdateProfile updates the current user's editable profile fields.
func (s *AuthService) UpdateProfile(ctx context.Context, userID, fullName, avatarURL string) error {
	if fullName == "" {
		return errors.NewValidation("full_name", []string{"full_name is required"})
	}
	return s.userRepo.UpdateProfile(ctx, userID, fullName, avatarURL)
}

// ChangePassword verifies the current password and sets a new one.
func (s *AuthService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.NewValidation("new_password", []string{"new password must be at least 6 characters"})
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.NewNotFound("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.NewValidation("current_password", []string{"current password is incorrect"})
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternal("failed to hash password")
	}
	return s.userRepo.UpdatePassword(ctx, userID, string(hash))
}

// GetSettings returns the current user's settings JSON blob ({} if unset).
func (s *AuthService) GetSettings(ctx context.Context, userID string) (models.JSONB, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewNotFound("user not found")
	}
	if len(user.Settings) == 0 {
		return models.JSONB([]byte("{}")), nil
	}
	return user.Settings, nil
}

// UpdateSettings replaces the current user's settings JSON blob.
func (s *AuthService) UpdateSettings(ctx context.Context, userID string, settings models.JSONB) error {
	if len(settings) == 0 {
		settings = models.JSONB([]byte("{}"))
	}
	return s.userRepo.UpdateSettings(ctx, userID, settings)
}

// VerifyDocument saves the uploaded business license file ID.
func (s *AuthService) VerifyDocument(ctx context.Context, userID, fileID string) error {
	if fileID == "" {
		return errors.NewValidation("file_id", []string{"file_id is required"})
	}
	// We don't change status to active here, admin will do it.
	return s.userRepo.UpdateVerificationFile(ctx, userID, fileID)
}

func (s *AuthService) GetMe(ctx context.Context, userID string) (*response.UserMeResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewNotFound("user not found")
	}

	companyRows, err := s.userRepo.FindUserCompanies(ctx, userID)
	if err != nil {
		return nil, errors.NewInternal("failed to fetch user companies")
	}

	companies := make([]response.CompanyRole, 0, len(companyRows))
	for _, row := range companyRows {
		companies = append(companies, response.CompanyRole{
			ID:   row.CompanyID,
			Name: row.CompanyName,
			Role: row.Role,
		})
	}

	return &response.UserMeResponse{
		ID:                 user.ID,
		Email:              user.Email,
		FullName:           user.FullName,
		Role:               string(user.Role),
		Status:             string(user.Status),
		AvatarURL:          user.AvatarURL.String,
		VerificationFileID: user.VerificationFileID.String,
		Companies:          companies,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, ipAddress, userAgent string) (*response.AuthTokens, string, error) {
	ipAddress = stripPort(ipAddress)

	claims, err := jwt.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, "", errors.NewUnauthorized("invalid or expired refresh token")
	}

	tokenHash := hashToken(refreshToken)
	tokenRecord, err := s.refreshTokenRepo.FindByHash(ctx, tokenHash)
	if err != nil {
		return nil, "", errors.NewUnauthorized("refresh token not found")
	}

	if tokenRecord.IsRevoked {
		// Replay attack detected — revoke toàn bộ family
		_ = s.refreshTokenRepo.RevokeFamily(ctx, tokenRecord.FamilyID)
		return nil, "", errors.NewUnauthorized("token revoked, please login again")
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return nil, "", errors.NewUnauthorized("refresh token expired")
	}

	// Mark old token as revoked
	if err := s.refreshTokenRepo.MarkAsRevoked(ctx, tokenRecord.ID); err != nil {
		return nil, "", errors.NewInternal("failed to revoke old token")
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, "", errors.NewUnauthorized("user not found")
	}

	// Defense-in-depth: blocking an admin-side account also revokes its
	// refresh tokens, but if that revoke ever fails (or status was changed
	// directly in the DB) the rotation path must still refuse a session.
	if e := ensureNotBlocked(user); e != nil {
		return nil, "", e
	}

	newAccessToken, _, err := jwt.GenerateTokenPair(
		user.ID, user.Email, string(user.Role), s.jwtSecret,
		15*time.Minute, 7*24*time.Hour,
	)
	if err != nil {
		return nil, "", errors.NewInternal("failed to generate tokens")
	}

	newRefreshToken, err := s.createAndSaveRefreshToken(ctx, user.ID, user.Email, string(user.Role), tokenRecord.FamilyID, ipAddress, userAgent)
	if err != nil {
		return nil, "", err
	}

	return &response.AuthTokens{
		AccessToken: newAccessToken,
	}, newRefreshToken, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	tokenRecord, err := s.refreshTokenRepo.FindByHash(ctx, tokenHash)
	if err != nil {
		return nil // Ignore if token not found, already logged out
	}

	return s.refreshTokenRepo.RevokeFamily(ctx, tokenRecord.FamilyID)
}

func (s *AuthService) LogoutAll(ctx context.Context, userID string) (int, error) {
	return s.refreshTokenRepo.RevokeAllByUserID(ctx, userID)
}

// DeleteAccount soft-deletes the user and revokes all their refresh tokens.
func (s *AuthService) DeleteAccount(ctx context.Context, userID string) error {
	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		return errors.NewNotFound("user not found")
	}
	if err := s.userRepo.SoftDelete(ctx, userID); err != nil {
		return errors.NewInternal("failed to delete account")
	}
	// Best-effort: revoke all sessions so the account can no longer be used.
	_, _ = s.refreshTokenRepo.RevokeAllByUserID(ctx, userID)
	return nil
}

const (
	otpPurposeRegister = "register"
	otpPurposeReset    = "reset"
)

// SendRegistrationOTP issues an OTP for a not-yet-registered email and emails it.
func (s *AuthService) SendRegistrationOTP(ctx context.Context, email string) error {
	if s.otpSvc == nil || !s.otpSvc.Enabled() {
		return errors.NewInternal("otp service is not configured")
	}
	email = NormalizeEmail(email)
	if email == "" {
		return errors.NewValidation("email", []string{"email is required"})
	}
	if allowed, err := s.otpSvc.AllowSend(ctx, otpPurposeRegister, email); err != nil {
		return errors.NewInternal("failed to check otp rate limit")
	} else if !allowed {
		return errors.NewValidation("email", []string{"vui lòng chờ trước khi yêu cầu mã mới"})
	}
	if _, err := s.userRepo.FindByEmail(ctx, email); err == nil {
		return errors.NewConflict("email already exists")
	}

	_, generationID, err := s.otpSvc.Generate(ctx, otpPurposeRegister, email)
	if err != nil {
		return errors.NewInternal("failed to generate otp")
	}
	if s.otpEnqueuer == nil {
		return errors.NewInternal("email queue is not configured")
	}
	if err := s.otpEnqueuer.EnqueueSendOTPEmail(queue.SendOTPEmailPayload{To: email, Purpose: otpPurposeRegister, GenerationID: generationID, TemplateType: "Xác nhận đăng ký tài khoản", CreatedAt: time.Now()}); err != nil {
		_ = s.otpSvc.Rollback(ctx, otpPurposeRegister, email, generationID)
		return errors.NewInternal("failed to queue otp email")
	}
	return nil
}

// ForgotPassword issues a password-reset OTP if the account exists.
// It always returns nil so callers don't leak whether an email is registered.
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	if s.otpSvc == nil || !s.otpSvc.Enabled() {
		return errors.NewInternal("otp service is not configured")
	}
	email = NormalizeEmail(email)
	if email == "" {
		return errors.NewValidation("email", []string{"email is required"})
	}
	if allowed, err := s.otpSvc.AllowSend(ctx, otpPurposeReset, email); err != nil {
		return errors.NewInternal("failed to check otp rate limit")
	} else if !allowed {
		return nil
	}
	if _, err := s.userRepo.FindByEmail(ctx, email); err != nil {
		return nil
	}

	_, generationID, err := s.otpSvc.Generate(ctx, otpPurposeReset, email)
	if err != nil {
		return errors.NewInternal("failed to generate otp")
	}
	if s.otpEnqueuer == nil {
		return errors.NewInternal("email queue is not configured")
	}
	if err := s.otpEnqueuer.EnqueueSendOTPEmail(queue.SendOTPEmailPayload{To: email, Purpose: otpPurposeReset, GenerationID: generationID, TemplateType: "Đặt lại mật khẩu tài khoản", CreatedAt: time.Now()}); err != nil {
		_ = s.otpSvc.Rollback(ctx, otpPurposeReset, email, generationID)
		return errors.NewInternal("failed to queue otp email")
	}
	return nil
}

// VerifyResetOTP checks a password-reset code without consuming other purposes.
// Kept separate so the FE can validate the OTP step before showing the new-password form.
// Note: this does NOT consume the code — ResetPassword does the final consume.
func (s *AuthService) VerifyResetOTP(ctx context.Context, email, code string) error {
	if s.otpSvc == nil || !s.otpSvc.Enabled() {
		return errors.NewInternal("otp service is not configured")
	}
	ok, err := s.otpSvc.Peek(ctx, otpPurposeReset, email, code)
	if err != nil {
		return errors.NewInternal("failed to verify otp")
	}
	if !ok {
		return errors.NewValidation("otp", []string{"mã xác nhận không hợp lệ hoặc đã hết hạn"})
	}
	return nil
}

// ResetPassword verifies the reset OTP, sets a new password, and revokes all sessions.
func (s *AuthService) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	if s.otpSvc == nil || !s.otpSvc.Enabled() {
		return errors.NewInternal("otp service is not configured")
	}
	if len(newPassword) < 6 {
		return errors.NewValidation("new_password", []string{"mật khẩu mới phải có ít nhất 6 ký tự"})
	}
	ok, err := s.otpSvc.Verify(ctx, otpPurposeReset, email, code)
	if err != nil {
		return errors.NewInternal("failed to verify otp")
	}
	if !ok {
		return errors.NewValidation("otp", []string{"mã xác nhận không hợp lệ hoặc đã hết hạn"})
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.NewNotFound("user not found")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternal("failed to hash password")
	}
	if err := s.userRepo.UpdatePassword(ctx, user.ID, string(hash)); err != nil {
		return errors.NewInternal("failed to update password")
	}
	// Revoke all existing sessions so a compromised account can't stay logged in.
	_, _ = s.refreshTokenRepo.RevokeAllByUserID(ctx, user.ID)
	return nil
}
