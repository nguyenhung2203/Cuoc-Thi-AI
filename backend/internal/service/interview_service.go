package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"os"
	"time"

	"github.com/google/uuid"

	"backend/internal/livekit"
	"backend/internal/models"
	"backend/internal/pkg/errors"
	"backend/internal/repository"
)

type InterviewService struct {
	repo      *repository.InterviewRepository
	reportSvc *ReportService
}

func NewInterviewService(repo *repository.InterviewRepository, reportSvc *ReportService) *InterviewService {
	return &InterviewService{
		repo:      repo,
		reportSvc: reportSvc,
	}
}

type CreateInterviewRequest struct {
	CompanyID        string    `json:"company_id"`
	JobID            string    `json:"job_id"`
	CandidateID      string    `json:"candidate_id"`
	RecruiterID      string    `json:"recruiter_id"`
	TemplateID       string    `json:"template_id"`
	RubricID         string    `json:"rubric_id"`
	Mode             string    `json:"mode"`
	Title            string    `json:"title"`
	ScheduledAt      time.Time `json:"scheduled_at"`
	ConsentRecording bool      `json:"consent_recording"`
	ConsentAI        bool      `json:"consent_ai"`
	CreatedBy        string    `json:"created_by"`
}

type CreateInterviewResponse struct {
	InterviewID string `json:"interview_id"`
	RoomCode    string `json:"room_code"`
	InviteToken string `json:"invite_token"`
}

func (s *InterviewService) CreateInterview(ctx context.Context, req CreateInterviewRequest) (*CreateInterviewResponse, error) {
	interviewID := uuid.NewString()
	roomID := uuid.NewString()
	roomCode := uuid.NewString()[:8] // simple stub
	inviteToken := uuid.NewString()

	// Hash the invite token
	hasher := sha256.New()
	hasher.Write([]byte(inviteToken))
	hash := hex.EncodeToString(hasher.Sum(nil))

	expiresAt := req.ScheduledAt.Add(24 * time.Hour) // valid 24h after scheduled time

	now := time.Now()

	interview := &models.Interview{
		ID:               interviewID,
		CompanyID:        sql.NullString{String: req.CompanyID, Valid: req.CompanyID != ""},
		JobID:            sql.NullString{String: req.JobID, Valid: req.JobID != ""},
		CandidateID:      req.CandidateID,
		RecruiterID:      sql.NullString{String: req.RecruiterID, Valid: req.RecruiterID != ""},
		TemplateID:       sql.NullString{String: req.TemplateID, Valid: req.TemplateID != ""},
		RubricID:         sql.NullString{String: req.RubricID, Valid: req.RubricID != ""},
		Mode:             req.Mode,
		Title:            req.Title,
		ScheduledAt:      sql.NullTime{Time: req.ScheduledAt, Valid: !req.ScheduledAt.IsZero()},
		Status:           "scheduled",
		RoomID:           sql.NullString{String: roomID, Valid: true},
		InviteTokenHash:  sql.NullString{String: hash, Valid: true},
		InviteExpiresAt:  sql.NullTime{Time: expiresAt, Valid: true},
		ConsentRecording: req.ConsentRecording,
		ConsentAI:        req.ConsentAI,
		CreatedBy:        sql.NullString{String: req.CreatedBy, Valid: req.CreatedBy != ""},
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	room := &models.InterviewRoom{
		ID:          roomID,
		InterviewID: interviewID,
		RoomCode:    roomCode,
		Status:           "waiting",
		Provider:         sql.NullString{String: "livekit", Valid: true},
		ConnectionConfig: []byte("{}"),
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.Create(ctx, interview, room); err != nil {
		return nil, err
	}

	return &CreateInterviewResponse{
		InterviewID: interviewID,
		RoomCode:    roomCode,
		InviteToken: inviteToken, // only returned once
	}, nil
}

func (s *InterviewService) JoinByToken(ctx context.Context, token string) (*repository.CandidateJoinInfo, error) {
	hash := token
	if len(token) != 64 {
		hasher := sha256.New()
		hasher.Write([]byte(token))
		hash = hex.EncodeToString(hasher.Sum(nil))
	}

	info, err := s.repo.GetCandidateJoinInfo(ctx, hash)
	if err != nil || info == nil {
		return nil, errors.NewNotFound("invalid or expired token")
	}

	livekitSecret := os.Getenv("LIVEKIT_API_SECRET")
	if livekitSecret == "" {
		livekitSecret = "devsecret"
	}
	livekitKey := os.Getenv("LIVEKIT_API_KEY")
	if livekitKey == "" {
		livekitKey = "devkey"
	}

	tokenString, err := livekit.GenerateToken(
		livekitKey,
		livekitSecret,
		info.RoomID,
		info.CandidateID,
		info.CandidateName,
		"candidate",
		info.InterviewID,
	)
	if err == nil {
		info.RoomAccessToken = tokenString
		info.RoomAccessTokenExpiresAt = time.Now().Add(4 * time.Hour).Format(time.RFC3339)
	}

	return info, nil
}

func (s *InterviewService) ListInterviews(ctx context.Context, companyID string, limit, offset int) ([]models.Interview, error) {
	items, err := s.repo.ListByCompany(ctx, companyID, limit, offset)
	if err != nil {
		return nil, errors.NewInternal("failed to list interviews")
	}
	return items, nil
}

func (s *InterviewService) GetInterview(ctx context.Context, interviewID, companyID string) (*models.Interview, error) {
	i, err := s.repo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return nil, errors.NewNotFound("interview not found")
	}
	return i, nil
}

func (s *InterviewService) GenerateRoomAccessToken(ctx context.Context, interviewID, companyID, userID string) (string, error) {
	// Verify recruiter access
	_, err := s.repo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return "", errors.NewNotFound("interview not found")
	}

	// Stub: return a dummy token instead of a real LiveKit token
	return "mock_livekit_token_" + userID + "_for_" + interviewID, nil
}

func (s *InterviewService) StartInterview(ctx context.Context, interviewID, companyID string) error {
	i, err := s.repo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return errors.NewNotFound("interview not found")
	}

	if i.Status != "scheduled" && i.Status != "waiting" {
		return errors.NewBadRequest("cannot start an interview in state: " + i.Status)
	}

	err = s.repo.UpdateStatus(ctx, interviewID, "active")
	if err != nil {
		return err
	}
	
	if i.RoomID.Valid {
		_ = s.repo.UpdateRoomStatus(ctx, i.RoomID.String, "active")
	}

	// Stub: Trigger webhook or Pub/Sub event here
	return nil
}

func (s *InterviewService) EndInterview(ctx context.Context, interviewID, companyID string) error {
	i, err := s.repo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return errors.NewNotFound("interview not found")
	}

	err = s.repo.UpdateStatus(ctx, interviewID, "completed")
	if err != nil {
		return err
	}
	
	if i.RoomID.Valid {
		_ = s.repo.UpdateRoomStatus(ctx, i.RoomID.String, "closed")
	}

	// Trigger report generation async
	if s.reportSvc != nil && i.JobID.Valid && i.ConsentAI {
		generatedBy := "system"
		if i.RecruiterID.Valid {
			generatedBy = i.RecruiterID.String
		}
		go func() {
			_, _ = s.reportSvc.GenerateReport(context.Background(), companyID, interviewID, i.JobID.String, generatedBy)
		}()
	}

	return nil
}
