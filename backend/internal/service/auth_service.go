package service

import (
	"context"
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
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, req request.RegisterRequest) (*response.AuthTokens, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternal("failed to hash password")
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
		Status:       models.UserStatusActive,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.NewConflict("email already exists")
	}

	accessToken, refreshToken, err := jwt.GenerateTokenPair(
		user.ID, user.Email, string(user.Role), s.jwtSecret,
		15*time.Minute, 7*24*time.Hour,
	)
	if err != nil {
		return nil, errors.NewInternal("failed to generate tokens")
	}

	return &response.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req request.LoginRequest) (*response.AuthTokens, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewUnauthorized("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.NewUnauthorized("invalid email or password")
	}

	accessToken, refreshToken, err := jwt.GenerateTokenPair(
		user.ID, user.Email, string(user.Role), s.jwtSecret,
		15*time.Minute, 7*24*time.Hour,
	)
	if err != nil {
		return nil, errors.NewInternal("failed to generate tokens")
	}

	return &response.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
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

func (s *AuthService) RefreshToken(ctx context.Context, req request.RefreshRequest) (*response.AuthTokens, error) {
	claims, err := jwt.ValidateToken(req.RefreshToken, s.jwtSecret)
	if err != nil {
		return nil, errors.NewUnauthorized("invalid or expired refresh token")
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.NewUnauthorized("user not found")
	}

	accessToken, refreshToken, err := jwt.GenerateTokenPair(
		user.ID, user.Email, string(user.Role), s.jwtSecret,
		15*time.Minute, 7*24*time.Hour,
	)
	if err != nil {
		return nil, errors.NewInternal("failed to generate tokens")
	}

	return &response.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
