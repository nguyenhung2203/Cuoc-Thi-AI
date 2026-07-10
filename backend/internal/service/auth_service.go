package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"backend/internal/dto/request"
	"backend/internal/dto/response"
	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/jwt"
	"backend/internal/repository"
)

type AuthService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtSecret        string
	emailSvc         *EmailService
}

func NewAuthService(userRepo *repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository, jwtSecret string, emailSvc *EmailService) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        jwtSecret,
		emailSvc:         emailSvc,
	}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
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

func generateOTP() string {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "123456" // fallback
	}
	return fmt.Sprintf("%06d", n.Int64())
}

func (s *AuthService) Register(ctx context.Context, req request.RegisterRequest, ipAddress, userAgent string) (*response.AuthResponse, string, error) {
	ipAddress = stripPort(ipAddress)

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", errors.NewInternal("failed to hash password")
	}

	role := models.UserRole(req.Role)
	if role != models.RoleAdmin && role != models.RoleRecruiter && role != models.RoleCandidate {
		role = models.RoleCandidate // default
	}

	user := &models.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		Role:         role,
		Status:       models.UserStatusPending, // Đổi thành pending
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", errors.NewConflict("email already exists")
	}

	// Generate OTP
	otpCode := generateOTP()
	expiresAt := time.Now().Add(15 * time.Minute)
	err = s.userRepo.CreateOTP(ctx, user.Email, otpCode, "register", expiresAt)
	if err != nil {
		fmt.Println("Failed to create OTP:", err)
	} else {
		// Send email in background
		go func() {
			if s.emailSvc != nil {
				_ = s.emailSvc.SendOTPEmail(user.Email, otpCode, "register")
			}
		}()
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
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      string(user.Role),
		AvatarURL: user.AvatarURL.String,
		Companies: companies,
	}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.NewNotFound("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.NewUnauthorized("invalid old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternal("failed to hash new password")
	}

	if err := s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return errors.NewInternal("failed to update password")
	}

	return nil
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

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Do not return error if user not found to prevent email enumeration
		return nil
	}

	otpCode := generateOTP()
	expiresAt := time.Now().Add(15 * time.Minute)

	err = s.userRepo.CreateOTP(ctx, user.Email, otpCode, "forgot_password", expiresAt)
	if err != nil {
		return errors.NewInternal("failed to generate reset token")
	}

	go func() {
		if s.emailSvc != nil {
			err := s.emailSvc.SendOTPEmail(email, otpCode, "forgot_password")
			if err != nil {
				fmt.Println("Gửi email thất bại:", err)
			}
		}
	}()

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, email, otp, newPassword string) error {
	otpRow, err := s.userRepo.GetOTP(ctx, email, otp, "forgot_password")
	if err != nil {
		return errors.NewUnauthorized("invalid or expired OTP")
	}

	if otpRow.Used {
		return errors.NewUnauthorized("OTP has already been used")
	}

	if time.Now().After(otpRow.ExpiresAt) {
		return errors.NewUnauthorized("OTP has expired")
	}
	
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.NewUnauthorized("user not found")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternal("failed to hash password")
	}

	if err := s.userRepo.UpdatePassword(ctx, user.ID, string(hash)); err != nil {
		return errors.NewInternal("failed to update password")
	}

	if err := s.userRepo.MarkOTPUsed(ctx, otpRow.ID); err != nil {
		fmt.Println("Error marking token as used:", err.Error())
	}

	// Logout all devices after password reset for security
	_, _ = s.LogoutAll(ctx, user.ID)

	return nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, email, otp string) error {
	otpRow, err := s.userRepo.GetOTP(ctx, email, otp, "register")
	if err != nil {
		return errors.NewUnauthorized("invalid or expired OTP")
	}

	if otpRow.Used {
		return errors.NewUnauthorized("OTP has already been used")
	}

	if time.Now().After(otpRow.ExpiresAt) {
		return errors.NewUnauthorized("OTP has expired")
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.NewUnauthorized("user not found")
	}

	if err := s.userRepo.MarkEmailVerified(ctx, user.ID); err != nil {
		return errors.NewInternal("failed to verify email")
	}

	if err := s.userRepo.MarkOTPUsed(ctx, otpRow.ID); err != nil {
		fmt.Println("Error marking token as used:", err.Error())
	}

	return nil
}

func (s *AuthService) ResendOTP(ctx context.Context, email, purpose string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil // Prevent email enumeration
	}

	otpCode := generateOTP()
	expiresAt := time.Now().Add(15 * time.Minute)

	err = s.userRepo.CreateOTP(ctx, user.Email, otpCode, purpose, expiresAt)
	if err != nil {
		return errors.NewInternal("failed to generate OTP")
	}

	go func() {
		if s.emailSvc != nil {
			err := s.emailSvc.SendOTPEmail(email, otpCode, purpose)
			if err != nil {
				fmt.Println("Gửi email thất bại:", err)
			}
		}
	}()

	return nil
}

