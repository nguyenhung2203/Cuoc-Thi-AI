package service

import (
	"context"

	"backend/internal/dto/response"
	"backend/internal/pkg/errors"
	"backend/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
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
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.NewNotFound("user not found")
	}
	if status != "active" && status != "blocked" && status != "pending" {
		return errors.NewBadRequest("invalid status value")
	}
	return s.userRepo.UpdateStatus(ctx, userID, status)
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
