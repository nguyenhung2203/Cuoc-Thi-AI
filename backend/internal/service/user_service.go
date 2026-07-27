package service

import (
	"context"

	"backend/internal/dto/response"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/logger"
	"backend/internal/repository"
)

type UserService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewUserService(userRepo *repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository) *UserService {
	return &UserService{userRepo: userRepo, refreshTokenRepo: refreshTokenRepo}
}

func (s *UserService) ListPendingUsers(ctx context.Context) ([]response.UserMeResponse, error) {
	users, err := s.userRepo.ListPendingUsers(ctx)
	if err != nil {
		return nil, errors.NewInternal("failed to list pending users")
	}
	var res []response.UserMeResponse
	for _, u := range users {
		res = append(res, response.UserMeResponse{
			ID:                 u.ID,
			Email:              u.Email,
			FullName:           u.FullName,
			Role:               string(u.Role),
			Status:             string(u.Status),
			AvatarURL:          u.AvatarURL.String,
			VerificationFileID: u.VerificationFileID.String,
		})
	}
	return res, nil
}

func (s *UserService) ListAllUsers(ctx context.Context) ([]response.UserMeResponse, error) {
	users, err := s.userRepo.ListAllUsers(ctx)
	if err != nil {
		return nil, errors.NewInternal("failed to list users")
	}
	var res []response.UserMeResponse
	for _, u := range users {
		res = append(res, response.UserMeResponse{
			ID:                 u.ID,
			Email:              u.Email,
			FullName:           u.FullName,
			Role:               string(u.Role),
			Status:             string(u.Status),
			AvatarURL:          u.AvatarURL.String,
			VerificationFileID: u.VerificationFileID.String,
		})
	}
	return res, nil
}

func (s *UserService) ApproveUser(ctx context.Context, userID string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.NewNotFound("user not found")
	}
	if user.Status != "pending" {
		return errors.NewBadRequest("user is not pending")
	}
	return s.userRepo.UpdateStatus(ctx, userID, "active")
}

func (s *UserService) UpdateUserStatus(ctx context.Context, userID string, status string) error {
	// Validate before touching the DB so a bad status never reaches the repo.
	if status != "active" && status != "blocked" && status != "pending" {
		return errors.NewBadRequest("invalid status value")
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.NewNotFound("user not found")
	}
	if err := s.userRepo.UpdateStatus(ctx, userID, status); err != nil {
		return errors.NewInternal("failed to update user status")
	}
	s.revokeSessionsIfBlocked(ctx, userID, status)
	return nil
}

// revokeSessionsIfBlocked revokes every refresh token of a user who has just
// been blocked, so their session dies as soon as the access token expires.
// Best-effort: a revoke failure must not roll back the block itself (the
// RefreshToken path re-checks status as defense-in-depth).
func (s *UserService) revokeSessionsIfBlocked(ctx context.Context, userID, newStatus string) {
	if newStatus != "blocked" || s.refreshTokenRepo == nil {
		return
	}
	if _, err := s.refreshTokenRepo.RevokeAllByUserID(ctx, userID); err != nil {
		logger.Error("failed to revoke sessions for blocked user", "user_id", userID, "error", err)
	}
}

type DashboardStats struct {
	TotalUsers      int `json:"total_users"`
	TotalCompanies  int `json:"total_companies"`
	TotalInterviews int `json:"total_interviews"`
	TotalCandidates int `json:"total_candidates"`
	PendingUsers    int `json:"pending_users"`
}

func (s *UserService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	stats, err := s.userRepo.GetDashboardStats(ctx)
	if err != nil {
		return nil, err
	}
	return &DashboardStats{
		TotalUsers:      stats["total_users"],
		TotalCompanies:  stats["total_companies"],
		TotalInterviews: stats["total_interviews"],
		TotalCandidates: stats["total_candidates"],
		PendingUsers:    stats["pending_users"],
	}, nil
}

func (s *UserService) GetReports(ctx context.Context) (*response.AdminReports, error) {
	totalUsers, totalCandidates, totalRecruiters, totalCompanies, totalInterviews, tokensIn, tokensOut, monthlyItems, growthItems, err := s.userRepo.GetReportsData(ctx)
	if err != nil {
		return nil, err
	}

	var monthlyTokenUsage []response.MonthlyUsageItem
	for _, m := range monthlyItems {
		monthlyTokenUsage = append(monthlyTokenUsage, response.MonthlyUsageItem{
			Month: m.Month,
			Usage: m.Usage,
		})
	}

	var userGrowthTrend []response.GrowthItem
	for _, g := range growthItems {
		userGrowthTrend = append(userGrowthTrend, response.GrowthItem{
			Period: g.Period,
			Users:  g.Users,
		})
	}

	// Simulated server stability metrics
	uptimeHours := 142.5
	cpuUsage := 12.4
	ramUsage := 48.2
	redisMemory := 14.2
	serverLatency := 124

	tokenUsageByModel := map[string]int64{
		"gemini-2.5-flash":          tokensIn + tokensOut,
		"gemini-2.0-flash-live-001": 0,
	}

	return &response.AdminReports{
		TotalUsers:         totalUsers,
		TotalCandidates:    totalCandidates,
		TotalRecruiters:    totalRecruiters,
		TotalCompanies:     totalCompanies,
		TotalInterviews:    totalInterviews,
		TotalTokenUsage:    tokensIn + tokensOut,
		InputTokens:        tokensIn,
		OutputTokens:       tokensOut,
		SystemUptimeHours:  uptimeHours,
		CpuUsagePercent:    cpuUsage,
		RamUsagePercent:    ramUsage,
		RedisMemoryMb:      redisMemory,
		ServerLatencyMs:    serverLatency,
		TokenUsageByModel:  tokenUsageByModel,
		MonthlyTokenUsage:  monthlyTokenUsage,
		UserGrowthTrend:    userGrowthTrend,
	}, nil
}
