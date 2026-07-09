package service

import (
	"context"

	"backend/internal/dto/response"

	"backend/internal/pkg/errors"
	"backend/internal/repository"
	"github.com/google/uuid"
)

type CandidatePortalService struct {
	repo          *repository.CandidatePortalRepository
	userRepo      *repository.UserRepository
	candidateRepo *repository.CandidateRepository
	jobRepo       *repository.JobRepository
}

func NewCandidatePortalService(
	repo *repository.CandidatePortalRepository,
	userRepo *repository.UserRepository,
	candidateRepo *repository.CandidateRepository,
	jobRepo *repository.JobRepository,
) *CandidatePortalService {
	return &CandidatePortalService{
		repo:          repo,
		userRepo:      userRepo,
		candidateRepo: candidateRepo,
		jobRepo:       jobRepo,
	}
}

func (s *CandidatePortalService) GetDashboardStats(ctx context.Context, userID string) (*response.CandidatePortalDashboardStats, error) {
	interviews, err := s.repo.GetInterviewsByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewInternal("failed to get interviews")
	}

	upcomingCount := 0
	for _, iv := range interviews {
		if iv.Status == "scheduled" || iv.Status == "active" {
			upcomingCount++
		}
	}

	mockCount, err := s.repo.GetCompletedMockCount(ctx, userID)
	if err != nil {
		return nil, errors.NewInternal("failed to get mock count")
	}

	avgScore, err := s.repo.GetAverageMockScore(ctx, userID)
	if err != nil {
		return nil, errors.NewInternal("failed to get avg score")
	}

	// Simple completeness logic
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewInternal("failed to get user")
	}
	completeness := 50
	if user.AvatarURL.Valid && user.AvatarURL.String != "" {
		completeness += 25
	}
	// Note: Proper CV check would require checking candidate rows or user's central CV

	return &response.CandidatePortalDashboardStats{
		UpcomingInterviews: upcomingCount,
		CompletedMockTests: mockCount,
		AverageMockScore:   avgScore,
		ProfileCompleteness: completeness,
	}, nil
}

func (s *CandidatePortalService) GetInterviews(ctx context.Context, userID string) ([]response.CandidatePortalInterview, error) {
	interviews, err := s.repo.GetInterviewsByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewInternal("failed to get interviews")
	}
	
	// Inject join link for candidates
	for i := range interviews {
		if interviews[i].InviteTokenHash != "" {
			interviews[i].JoinLink = "/interview-consent?token=" + interviews[i].InviteTokenHash
		}
	}

	return interviews, nil
}

func (s *CandidatePortalService) GetProfile(ctx context.Context, userID string) (*response.CandidatePortalProfile, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewNotFound("user not found")
	}

	cvUrl := ""
	cvName := ""
	fileID, fileName, storageKey, _ := s.repo.GetUserLatestCV(ctx, userID)
	if fileName != "" && fileID != "" {
		cvUrl = "http://localhost:18080/uploads/" + storageKey
		cvName = fileName
	}

	return &response.CandidatePortalProfile{
		UserID:    user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		AvatarURL: user.AvatarURL.String,
		CVUrl:     cvUrl,
		CVName:    cvName,
	}, nil
}

// UploadCV handles updating the central CV for the user
func (s *CandidatePortalService) UploadCV(ctx context.Context, userID string, originalName string, cvFileID string) error {
	return s.repo.UpdateUserCV(ctx, userID, cvFileID, originalName)
}

func (s *CandidatePortalService) ApplyForJob(ctx context.Context, userID, jobID, cvFileID, cvOriginalName string) error {
	// 1. Get Job to verify it exists and is open
	job, err := s.jobRepo.GetByIDAnyCompany(ctx, jobID)
	if err != nil || job == nil {
		return errors.NewNotFound("job not found")
	}
	if job.Status != "open" {
		return errors.NewBadRequest("job is not open for applications")
	}

	// 2. Get User
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.NewNotFound("user not found")
	}

	// 3. Create or update Candidate record for this company
	candidateID := uuid.NewString()
	err = s.repo.ApplyForJob(ctx, candidateID, job.CompanyID, user.ID, user.FullName, user.Email, "", cvFileID, cvOriginalName, job.ID)
	if err != nil {
		return errors.NewInternal("failed to apply for job")
	}
	
	return nil
}
