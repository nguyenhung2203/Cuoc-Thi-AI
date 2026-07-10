package service

import (
	"context"
	"encoding/json"
	"log"

	"backend/internal/dto/response"

	"backend/internal/pkg/errors"
	"backend/internal/repository"
	"backend/internal/models"
	"github.com/google/uuid"
)

type CandidatePortalService struct {
	repo          *repository.CandidatePortalRepository
	userRepo      *repository.UserRepository
	candidateRepo *repository.CandidateRepository
	jobRepo       *repository.JobRepository
	aiSvc         *AIService
	notifRepo     *repository.NotificationRepository
}

func NewCandidatePortalService(
	repo *repository.CandidatePortalRepository,
	userRepo *repository.UserRepository,
	candidateRepo *repository.CandidateRepository,
	jobRepo *repository.JobRepository,
	aiSvc *AIService,
	notifRepo *repository.NotificationRepository,
) *CandidatePortalService {
	return &CandidatePortalService{
		repo:          repo,
		userRepo:      userRepo,
		candidateRepo: candidateRepo,
		jobRepo:       jobRepo,
		aiSvc:         aiSvc,
		notifRepo:     notifRepo,
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

// UpdateProfile updates the candidate's editable profile fields (name, avatar).
func (s *CandidatePortalService) UpdateProfile(ctx context.Context, userID, fullName, avatarURL string) error {
	if fullName == "" {
		return errors.NewValidation("full_name", []string{"full_name is required"})
	}
	return s.userRepo.UpdateProfile(ctx, userID, fullName, avatarURL)
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

	profile := &response.CandidatePortalProfile{
		UserID:    user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		AvatarURL: user.AvatarURL.String,
		CVFileID:  fileID,
		CVUrl:     cvUrl,
		CVName:    cvName,
	}

	// Attach AI-parsed CV data if available.
	if parsedJSON, err := s.repo.GetParsedCV(ctx, userID); err == nil && parsedJSON != "" {
		var parsed interface{}
		if json.Unmarshal([]byte(parsedJSON), &parsed) == nil {
			profile.ParsedData = parsed
		}
	}

	return profile, nil
}

// UploadCV updates the central CV for the user, then parses it with AI (best-effort).
// storageKey is the on-disk key of the freshly uploaded file so we can read/parse it.
func (s *CandidatePortalService) UploadCV(ctx context.Context, userID, originalName, cvFileID, storageKey string) error {
	if err := s.repo.UpdateUserCV(ctx, userID, cvFileID, originalName); err != nil {
		return err
	}

	// Parse the CV with AI and persist the result. Best-effort: upload still
	// succeeds even if parsing fails (e.g. AI unavailable or non-PDF content).
	if s.aiSvc != nil && storageKey != "" {
		parsedJSON, summary, perr := s.aiSvc.ExtractAndParseCV(ctx, storageKey)
		if perr != nil {
			log.Printf("portal: CV parse failed for user %s: %v", userID, perr)
			return nil
		}
		if serr := s.repo.SaveParsedCV(ctx, userID, parsedJSON, summary); serr != nil {
			log.Printf("portal: failed to save parsed CV for user %s: %v", userID, serr)
		}
	}
	return nil
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
	
	// Notify Recruiter (Job Creator) and Candidate
	if s.notifRepo != nil {
		// Notify Recruiter
		if job.CreatedBy != "" {
			_ = s.notifRepo.Create(ctx, &models.Notification{
				UserID:  job.CreatedBy,
				Title:   "Ứng viên mới",
				Message: "Ứng viên " + user.FullName + " vừa ứng tuyển vào vị trí: " + job.Title,
				Type:    "new_applicant",
				Link:    "/candidates",
			})
		}
		// Notify Candidate
		_ = s.notifRepo.Create(ctx, &models.Notification{
			UserID:  user.ID,
			Title:   "Ứng tuyển thành công",
			Message: "Bạn đã ứng tuyển thành công vào vị trí: " + job.Title,
			Type:    "job_applied",
			Link:    "/job-board",
		})
	}

	return nil
}
