package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"backend/internal/dto/response"

	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/pkg/utils"
	"backend/internal/repository"
	"github.com/google/uuid"
)

type CandidatePortalService struct {
	repo          *repository.CandidatePortalRepository
	userRepo      *repository.UserRepository
	candidateRepo *repository.CandidateRepository
	jobRepo       *repository.JobRepository
	aiSvc         *AIService
	orchestrator  *AIOrchestratorService
	matchCache    *MatchCacheService
	notifRepo     *repository.NotificationRepository
	fileSvc       *FileService
	// matchEnqueuer, when set, schedules an async recompute of all of a user's
	// application fit scores (wired to the asynq queue in main.go). When nil,
	// recompute runs inline in a background goroutine.
	matchEnqueuer func(userID string) error
}

// SetMatchEnqueuer wires the async match-recompute dispatcher (asynq). Optional:
// without it, CV re-upload recomputes matches inline.
func (s *CandidatePortalService) SetMatchEnqueuer(fn func(userID string) error) {
	s.matchEnqueuer = fn
}

func NewCandidatePortalService(
	repo *repository.CandidatePortalRepository,
	userRepo *repository.UserRepository,
	candidateRepo *repository.CandidateRepository,
	jobRepo *repository.JobRepository,
	aiSvc *AIService,
	orchestrator *AIOrchestratorService,
	matchCache *MatchCacheService,
	notifRepo *repository.NotificationRepository,
	fileSvc *FileService,
) *CandidatePortalService {
	return &CandidatePortalService{
		repo:          repo,
		userRepo:      userRepo,
		candidateRepo: candidateRepo,
		jobRepo:       jobRepo,
		aiSvc:         aiSvc,
		orchestrator:  orchestrator,
		matchCache:    matchCache,
		notifRepo:     notifRepo,
		fileSvc:       fileSvc,
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

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewInternal("failed to get user")
	}
	completeness := 50
	if user.AvatarURL.Valid && user.AvatarURL.String != "" {
		completeness += 25
	}
	if fileID, _, _, cvErr := s.repo.GetUserLatestCV(ctx, userID); cvErr == nil && fileID != "" {
		completeness += 25
	}

	return &response.CandidatePortalDashboardStats{
		UpcomingInterviews:  upcomingCount,
		CompletedMockTests:  mockCount,
		AverageMockScore:    avgScore,
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
		cvUrl, _ = s.fileSvc.SignedURL(storageKey)
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
	if status, err := s.repo.GetUserCVParseStatus(ctx, userID); err == nil {
		profile.CVParseStatus = status
	}

	// Attach AI-parsed CV data if available.
	if parsedJSON, err := s.repo.GetParsedCV(ctx, userID); err == nil && parsedJSON != "" {
		var parsed interface{}
		if json.Unmarshal([]byte(parsedJSON), &parsed) == nil {
			profile.ParsedData = parsed
		}
	}
	if score, err := s.repo.GetLatestInterviewScore(ctx, userID); err == nil {
		profile.InterviewScore = score
	}

	return profile, nil
}

// UploadCV saves the CV file metadata even if AI parse fails (best-effort analysis).
func (s *CandidatePortalService) UploadCV(ctx context.Context, userID, originalName, cvFileID, storageKey string) (*response.CandidatePortalCVUpload, error) {
	if storageKey == "" {
		return nil, errors.NewInternal("CV storage unavailable")
	}
	if err := s.repo.UpdateUserCV(ctx, userID, cvFileID, originalName); err != nil {
		return nil, err
	}

	result := &response.CandidatePortalCVUpload{
		Message:     "CV uploaded",
		FileName:    originalName,
		CVFileID:    cvFileID,
		ParseStatus: "pending",
	}

	if s.aiSvc == nil {
		_ = s.repo.UpdateCVAIStatus(ctx, userID, "failed", fmt.Errorf("AI service unavailable"))
		result.ParseStatus = "failed"
		result.Message = "CV uploaded; analysis unavailable"
		return result, nil
	}

	if err := s.repo.UpdateCVAIStatus(ctx, userID, "processing", nil); err != nil {
		result.ParseStatus = "failed"
		result.Message = "CV uploaded; could not start analysis"
		return result, nil
	}

	parsedJSON, summary, err := s.aiSvc.ExtractAndParseCV(ctx, storageKey)
	if err != nil {
		_ = s.repo.UpdateCVAIStatus(ctx, userID, "failed", err)
		result.ParseStatus = "failed"
		result.Message = "CV uploaded; analysis failed"
		return result, nil
	}
	var parsed interface{}
	if err := json.Unmarshal([]byte(parsedJSON), &parsed); err != nil {
		_ = s.repo.UpdateCVAIStatus(ctx, userID, "failed", err)
		result.ParseStatus = "failed"
		result.Message = "CV uploaded; analysis returned invalid data"
		return result, nil
	}
	saved, err := s.repo.SaveParsedCV(ctx, userID, cvFileID, parsedJSON, summary)
	if err != nil {
		_ = s.repo.UpdateCVAIStatus(ctx, userID, "failed", err)
		result.ParseStatus = "failed"
		result.Message = "CV uploaded; failed to save analysis"
		return result, nil
	}
	// saved=false can mean no candidate rows yet (first portal upload) or a race with a newer file.
	// Still return parsed data so the UI can show analysis immediately.
	if saved {
		if s.matchCache != nil {
			s.matchCache.InvalidateUser(ctx, userID)
		}
		if s.matchEnqueuer != nil {
			_ = s.matchEnqueuer(userID)
		}
	}
	result.ParseStatus = "ready"
	result.ParsedData = parsed
	result.Message = "CV uploaded and analyzed"
	return result, nil
}

// ReparseCV re-runs AI parse on the user's latest owned CV and persists the
// result on the file row (and candidate rows when present).
func (s *CandidatePortalService) ReparseCV(ctx context.Context, userID string) (*response.CandidatePortalCVUpload, error) {
	fileID, fileName, storageKey, err := s.repo.GetUserLatestCV(ctx, userID)
	if err != nil {
		return nil, errors.NewInternal("failed to load CV")
	}
	if fileID == "" || storageKey == "" {
		return nil, errors.NewBadRequest("no CV uploaded")
	}

	result := &response.CandidatePortalCVUpload{
		Message:     "CV re-analyzed",
		FileName:    fileName,
		CVFileID:    fileID,
		ParseStatus: "pending",
	}
	if s.aiSvc == nil {
		_ = s.repo.UpdateCVAIStatus(ctx, userID, "failed", fmt.Errorf("AI service unavailable"))
		result.ParseStatus = "failed"
		result.Message = "Analysis unavailable"
		return result, nil
	}
	if err := s.repo.UpdateCVAIStatus(ctx, userID, "processing", nil); err != nil {
		result.ParseStatus = "failed"
		result.Message = "Could not start analysis: " + err.Error()
		return result, nil
	}

	parsedJSON, summary, err := s.aiSvc.ExtractAndParseCV(ctx, storageKey)
	if err != nil {
		_ = s.repo.UpdateCVAIStatus(ctx, userID, "failed", err)
		result.ParseStatus = "failed"
		result.Message = "Analysis failed"
		return result, nil
	}
	var parsed interface{}
	if err := json.Unmarshal([]byte(parsedJSON), &parsed); err != nil {
		_ = s.repo.UpdateCVAIStatus(ctx, userID, "failed", err)
		result.ParseStatus = "failed"
		result.Message = "Analysis returned invalid data"
		return result, nil
	}
	saved, err := s.repo.SaveParsedCV(ctx, userID, fileID, parsedJSON, summary)
	if err != nil {
		_ = s.repo.UpdateCVAIStatus(ctx, userID, "failed", err)
		result.ParseStatus = "failed"
		result.Message = "Failed to save analysis"
		return result, nil
	}
	if saved {
		if s.matchCache != nil {
			s.matchCache.InvalidateUser(ctx, userID)
		}
		if s.matchEnqueuer != nil {
			_ = s.matchEnqueuer(userID)
		}
	}
	result.ParseStatus = "ready"
	result.ParsedData = parsed
	result.Message = "CV analyzed"
	return result, nil
}

// DeleteCV permanently removes a candidate-owned CV and purges remaining portal
// CV uploads so an older file cannot reappear after reload.
func (s *CandidatePortalService) DeleteCV(ctx context.Context, userID, fileID string) error {
	if userID == "" || fileID == "" {
		return errors.NewBadRequest("invalid CV id")
	}
	if s.fileSvc == nil {
		return errors.NewInternal("CV storage unavailable")
	}
	owned, err := s.fileSvc.GetOwnedCV(ctx, userID, fileID)
	if err != nil {
		return err
	}
	if owned == nil {
		return errors.NewNotFound("CV not found")
	}
	if err := s.repo.ClearAllUserCVs(ctx, userID); err != nil {
		return errors.NewInternal("failed to detach CV from profile")
	}
	if err := s.fileSvc.PurgeOwnedCVFiles(ctx, userID); err != nil {
		return err
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
		if err == repository.DuplicateApplication {
			return errors.NewConflict("you have already applied to this job")
		}
		return errors.NewInternal("failed to apply for job")
	}

	// Compute AI fit score in the background (best-effort). The application
	// succeeds regardless; the score fills in once the AI responds. Only runs
	// when the user has a parsed CV — otherwise fit_score stays NULL (no faking).
	go func(uid, companyID, jobID string) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		fitScore, matchJSON, merr := s.computeMatch(bgCtx, uid, job)
		if merr != nil {
			log.Printf("portal: fit score compute failed for user %s job %s: %v", uid, jobID, merr)
			return
		}
		// Resolve the candidate row id for this user+company to key the update.
		cid, cerr := s.repo.GetCandidateIDByUserCompany(bgCtx, uid, companyID)
		if cerr != nil || cid == "" {
			log.Printf("portal: cannot resolve candidate for fit score update user %s: %v", uid, cerr)
			return
		}
		if uerr := s.repo.UpdateMatch(bgCtx, jobID, cid, sql.NullFloat64{Float64: fitScore, Valid: true}, models.JSONB(matchJSON)); uerr != nil {
			log.Printf("portal: failed to save fit score user %s job %s: %v", uid, jobID, uerr)
		}
		// Seed the cache so a later preview view is instant.
		var mr MatchResult
		if json.Unmarshal(matchJSON, &mr) == nil {
			s.matchCache.Set(bgCtx, uid, jobID, &mr)
		}
	}(user.ID, job.CompanyID, job.ID)

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

func (s *CandidatePortalService) GetApplications(ctx context.Context, userID string) ([]response.CandidatePortalApplication, error) {
	apps, err := s.repo.GetApplicationsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	for i := range apps {
		if apps[i].CVStorageKey != "" && s.fileSvc != nil {
			url, _ := s.fileSvc.SignedURL(apps[i].CVStorageKey)
			if url != "" {
				apps[i].CVUrl = url
			}
		}
	}
	return apps, nil
}

func (s *CandidatePortalService) CancelApplication(ctx context.Context, userID string, applicationID string) error {
	return s.repo.CancelApplication(ctx, userID, applicationID)
}

// MatchResult mirrors the match_cv_job prompt template output.
type MatchResult struct {
	FitScore       float64  `json:"fit_score"`
	MatchedSkills  []string `json:"matched_skills"`
	MissingSkills  []string `json:"missing_skills"`
	Summary        string   `json:"summary"`
	Recommendation string   `json:"recommendation"`
	Confidence     float64  `json:"confidence"`
}

// ErrNoCV is returned by match paths when the user has no parsed CV yet.
var ErrNoCV = errors.NewBadRequest("no CV available to compute match")

// computeMatch runs AI matching of the user's parsed CV against a job. Returns
// the fit score (0-100) and the raw match JSON. Returns ErrNoCV when the user
// has no parsed CV — callers should treat that as "leave fit_score NULL".
func (s *CandidatePortalService) computeMatch(ctx context.Context, userID string, job *models.Job) (float64, []byte, error) {
	if s.orchestrator == nil {
		return 0, nil, errors.NewInternal("AI matching unavailable")
	}

	parsedJSON, err := s.repo.GetParsedCV(ctx, userID)
	if err != nil {
		return 0, nil, err
	}
	if parsedJSON == "" {
		return 0, nil, ErrNoCV
	}

	// Extract a compact skills string + summary from the parsed CV JSON.
	cvSummary := parsedJSON
	cvSkills := ""
	var parsed struct {
		Skills     []string `json:"skills"`
		Experience string   `json:"experience"`
		Education  string   `json:"education"`
	}
	if json.Unmarshal([]byte(parsedJSON), &parsed) == nil {
		cvSkills = strings.Join(parsed.Skills, ", ")
		if parsed.Experience != "" || parsed.Education != "" {
			cvSummary = fmt.Sprintf("Experience: %s. Education: %s", parsed.Experience, parsed.Education)
		}
	}

	variables := map[string]string{
		"cv_summary":       cvSummary,
		"cv_skills":        cvSkills,
		"job_title":        job.Title,
		"job_description":  job.Description,
		"job_requirements": job.Requirements.String,
	}

	data, err := s.orchestrator.CallAI(ctx, "match_cv_job", "", variables)
	if err != nil {
		return 0, nil, err
	}

	var result MatchResult
	cleanJSON := utils.CleanJSON(string(data))
	if uerr := json.Unmarshal([]byte(cleanJSON), &result); uerr != nil {
		return 0, nil, fmt.Errorf("failed to parse match result: %w", uerr)
	}

	fit := result.FitScore
	if fit < 0 {
		fit = 0
	}
	if fit > 100 {
		fit = 100
	}
	return fit, []byte(cleanJSON), nil
}

// ReviewCV runs AI CV improvement feedback for the candidate portal (not interview scoring).
func (s *CandidatePortalService) ReviewCV(ctx context.Context, userID string) (*response.CandidatePortalCVReview, error) {
	if s.aiSvc == nil {
		return nil, errors.NewInternal("AI service unavailable")
	}
	fileID, _, storageKey, err := s.repo.GetUserLatestCV(ctx, userID)
	if err != nil {
		return nil, errors.NewInternal("failed to load CV")
	}
	if fileID == "" || storageKey == "" {
		return nil, errors.NewBadRequest("please upload a CV before requesting a review")
	}
	result, err := s.aiSvc.ReviewCV(ctx, storageKey)
	if err != nil {
		return nil, errors.NewInternal("CV review failed")
	}
	return &response.CandidatePortalCVReview{
		Summary:         result.Summary,
		Issues:          result.Issues,
		Suggestions:     result.Suggestions,
		MissingSections: result.MissingSections,
		Strengths:       result.Strengths,
	}, nil
}

// GetJobMatch computes an on-demand CV↔job match for the Apply page preview.
// Returns ErrNoCV when the user has not uploaded/parsed a CV.
func (s *CandidatePortalService) GetJobMatch(ctx context.Context, userID, jobID string) (*MatchResult, error) {
	// Cache hit: skip the model entirely.
	if cached, ok := s.matchCache.Get(ctx, userID, jobID); ok {
		return cached, nil
	}

	job, err := s.jobRepo.GetByIDAnyCompany(ctx, jobID)
	if err != nil || job == nil {
		return nil, errors.NewNotFound("job not found")
	}
	_, matchJSON, err := s.computeMatch(ctx, userID, job)
	if err != nil {
		return nil, err
	}
	var result MatchResult
	if uerr := json.Unmarshal(matchJSON, &result); uerr != nil {
		return nil, errors.NewInternal("failed to parse match result")
	}
	s.matchCache.Set(ctx, userID, jobID, &result)
	return &result, nil
}

// RecomputeUserMatches recomputes and persists the fit score for every job the
// user has applied to, and reseeds the cache. Called after a CV re-upload (via
// the async queue when wired, otherwise inline). Best-effort per application:
// one failure does not abort the rest.
func (s *CandidatePortalService) RecomputeUserMatches(ctx context.Context, userID string) error {
	apps, err := s.repo.ListUserApplications(ctx, userID)
	if err != nil {
		return err
	}
	for _, app := range apps {
		job, jerr := s.jobRepo.GetByIDAnyCompany(ctx, app.JobID)
		if jerr != nil || job == nil {
			continue
		}
		fitScore, matchJSON, merr := s.computeMatch(ctx, userID, job)
		if merr != nil {
			// ErrNoCV means the CV is gone; nothing to recompute. Any other error
			// is transient — skip this application, keep going.
			if merr == ErrNoCV {
				return nil
			}
			log.Printf("portal: recompute match failed for user %s job %s: %v", userID, app.JobID, merr)
			continue
		}
		if uerr := s.repo.UpdateMatch(ctx, app.JobID, app.CandidateID, sql.NullFloat64{Float64: fitScore, Valid: true}, models.JSONB(matchJSON)); uerr != nil {
			log.Printf("portal: recompute update failed for user %s job %s: %v", userID, app.JobID, uerr)
			continue
		}
		var result MatchResult
		if json.Unmarshal(matchJSON, &result) == nil {
			s.matchCache.Set(ctx, userID, app.JobID, &result)
		}
	}
	return nil
}
