package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"time"

	"github.com/google/uuid"

	"backend/internal/livekit"
	"backend/internal/models"
	"backend/internal/pkg/email"
	"backend/internal/pkg/errors"
	"backend/internal/repository"
)

type InterviewService struct {
	repo          *repository.InterviewRepository
	reportSvc     *ReportService
	candidateRepo *repository.CandidateRepository
	notifRepo     *repository.NotificationRepository
	mailer        *email.Sender
	frontendURL   string
}

func NewInterviewService(repo *repository.InterviewRepository, reportSvc *ReportService, notifRepo *repository.NotificationRepository) *InterviewService {
	return &InterviewService{
		repo:      repo,
		reportSvc: reportSvc,
		notifRepo: notifRepo,
	}
}

// WithMailer attaches email/candidate dependencies for invite & reminder emails.
// Optional so existing construction keeps working.
func (s *InterviewService) WithMailer(candidateRepo *repository.CandidateRepository, mailer *email.Sender, frontendURL string) *InterviewService {
	s.candidateRepo = candidateRepo
	s.mailer = mailer
	s.frontendURL = frontendURL
	return s
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
		ID:               roomID,
		InterviewID:      interviewID,
		RoomCode:         roomCode,
		Status:           "waiting",
		Provider:         sql.NullString{String: "livekit", Valid: true},
		ConnectionConfig: []byte("{}"),
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.Create(ctx, interview, room); err != nil {
		return nil, err
	}

	// Notify candidate if possible
	if s.candidateRepo != nil && s.notifRepo != nil {
		if cand, candErr := s.candidateRepo.GetByID(ctx, req.CompanyID, req.CandidateID); candErr == nil && cand != nil {
			if cand.UserID.Valid {
				_ = s.notifRepo.Create(ctx, &models.Notification{
					UserID:  cand.UserID.String,
					Title:   "Lịch phỏng vấn mới",
					Message: "Bạn vừa được mời tham gia phỏng vấn: " + req.Title,
					Type:    "interview_invite",
					Link:    "/my-interviews",
				})
			}
		}
	}

	return &CreateInterviewResponse{
		InterviewID: interviewID,
		RoomCode:    roomCode,
		InviteToken: inviteToken, // only returned once
	}, nil
}

type CandidateJoinInfo struct {
	InterviewID              string    `json:"interview_id"`
	RoomID                   string    `json:"room_id"`
	CandidateName            string    `json:"candidate_name"`
	CompanyName              string    `json:"company_name"`
	JobTitle                 string    `json:"job_title"`
	ScheduledAt              time.Time `json:"scheduled_at"`
	RoomAccessToken          string    `json:"room_access_token"`
	RoomAccessTokenExpiresAt string    `json:"room_access_token_expires_at"`
}

func (s *InterviewService) JoinByToken(ctx context.Context, token string) (*CandidateJoinInfo, error) {
	var hash string
	if len(token) == 64 {
		hash = token
	} else {
		hasher := sha256.New()
		hasher.Write([]byte(token))
		hash = hex.EncodeToString(hasher.Sum(nil))
	}

	info, err := s.repo.GetCandidateJoinInfoByInviteTokenHash(ctx, hash)
	if err != nil || info == nil {
		return nil, errors.NewNotFound("invalid or expired token")
	}

	if info.InviteExpiresAt.Before(time.Now()) {
		return nil, errors.NewForbidden("token expired")
	}

	// Generate LiveKit token for candidate — credentials come from the shared
	// config (production refuses dev fallbacks at startup).
	lkCfg := livekit.Current()

	identity := "candidate-" + info.CandidateID
	if info.UserID != nil {
		identity = *info.UserID
	}

	roomID := "room-" + info.InterviewID
	if room, rerr := s.repo.GetRoomByInterviewID(ctx, info.InterviewID); rerr == nil && room != nil && room.RoomCode != "" {
		roomID = room.RoomCode
	}

	tokenString, err := livekit.GenerateToken(
		lkCfg.APIKey,
		lkCfg.APISecret,
		roomID,
		identity,
		info.CandidateName,
		"candidate",
		info.InterviewID,
	)
	if err != nil {
		return nil, errors.NewInternal("failed to generate candidate token")
	}

	expiresAt := time.Now().Add(4 * time.Hour)

	return &CandidateJoinInfo{
		InterviewID:              info.InterviewID,
		RoomID:                   roomID,
		CandidateName:            info.CandidateName,
		CompanyName:              info.CompanyName,
		JobTitle:                 info.JobTitle,
		ScheduledAt:              info.ScheduledAt,
		RoomAccessToken:          tokenString,
		RoomAccessTokenExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
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

// SendReminder emails the candidate a reminder for an upcoming interview.
// Returns the recipient email on success.
func (s *InterviewService) SendReminder(ctx context.Context, interviewID, companyID string) (string, error) {
	iv, err := s.repo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return "", errors.NewNotFound("interview not found")
	}
	if s.candidateRepo == nil || s.mailer == nil {
		return "", errors.NewInternal("email service not configured")
	}
	cand, err := s.candidateRepo.GetByID(ctx, companyID, iv.CandidateID)
	if err != nil || cand == nil {
		return "", errors.NewNotFound("candidate not found")
	}
	if cand.Email == "" {
		return "", errors.NewBadRequest("candidate has no email")
	}

	when := "thời gian đã hẹn"
	if iv.ScheduledAt.Valid {
		when = iv.ScheduledAt.Time.Format("15:04 02/01/2006")
	}
	subject := "Nhắc lịch phỏng vấn: " + iv.Title
	body := "Xin chào " + cand.FullName + ",\n\n" +
		"Đây là lời nhắc cho buổi phỏng vấn \"" + iv.Title + "\" vào " + when + ".\n" +
		"Vui lòng đăng nhập hệ thống để tham gia đúng giờ.\n\n" +
		"Trân trọng,\nAI Interview Platform"

	if err := s.mailer.Send(cand.Email, subject, body); err != nil {
		return "", errors.NewInternal("failed to send reminder email")
	}
	return cand.Email, nil
}

// UpdateNotes saves recruiter internal notes after verifying the interview exists in the company.
func (s *InterviewService) UpdateNotes(ctx context.Context, interviewID, companyID, notes string) error {
	if _, err := s.repo.GetByIDAndCompany(ctx, interviewID, companyID); err != nil {
		return errors.NewNotFound("interview not found")
	}
	return s.repo.UpdateNotes(ctx, interviewID, companyID, notes)
}

// RescheduleInterview changes an interview's scheduled time. Only allowed while
// the interview is still in the "scheduled" state — once it is waiting/active/
// completed/cancelled, rescheduling is rejected (cancel + recreate instead).
func (s *InterviewService) RescheduleInterview(ctx context.Context, interviewID, companyID string, scheduledAt time.Time) error {
	i, err := s.repo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return errors.NewNotFound("interview not found")
	}
	if i.Status != "scheduled" {
		return errors.NewBadRequest("can only reschedule an interview in 'scheduled' state, current: " + i.Status)
	}
	if err := s.repo.UpdateScheduledAt(ctx, interviewID, companyID, scheduledAt); err != nil {
		return err
	}

	// Notify candidate about the new time if possible.
	if s.candidateRepo != nil && s.notifRepo != nil {
		if cand, candErr := s.candidateRepo.GetByID(ctx, companyID, i.CandidateID); candErr == nil && cand != nil {
			if cand.UserID.Valid {
				_ = s.notifRepo.Create(ctx, &models.Notification{
					UserID:  cand.UserID.String,
					Title:   "Lịch phỏng vấn được dời",
					Message: "Lịch phỏng vấn: " + i.Title + " đã được dời sang thời gian mới.",
					Type:    "interview_rescheduled",
					Link:    "/my-interviews",
				})
			}
		}
	}

	return nil
}

// GetRoom returns room metadata for an interview (after verifying company scope).
func (s *InterviewService) GetRoom(ctx context.Context, interviewID, companyID string) (*models.InterviewRoom, error) {
	if _, err := s.repo.GetByIDAndCompany(ctx, interviewID, companyID); err != nil {
		return nil, errors.NewNotFound("interview not found")
	}
	room, err := s.repo.GetRoomByInterviewID(ctx, interviewID)
	if err != nil || room == nil {
		return nil, errors.NewNotFound("room not found")
	}
	return room, nil
}

// CancelInterview marks an interview cancelled (only if not already ended).
func (s *InterviewService) CancelInterview(ctx context.Context, interviewID, companyID string) error {
	i, err := s.repo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return errors.NewNotFound("interview not found")
	}
	if i.Status == "completed" || i.Status == "cancelled" {
		return errors.NewBadRequest("cannot cancel an interview in state: " + i.Status)
	}
	if err := s.repo.UpdateStatus(ctx, interviewID, companyID, "cancelled"); err != nil {
		return err
	}
	if i.RoomID.Valid {
		_ = s.repo.UpdateRoomStatus(ctx, i.RoomID.String, "closed")
	}

	// Notify candidate if possible
	if s.candidateRepo != nil && s.notifRepo != nil {
		if cand, candErr := s.candidateRepo.GetByID(ctx, companyID, i.CandidateID); candErr == nil && cand != nil {
			if cand.UserID.Valid {
				_ = s.notifRepo.Create(ctx, &models.Notification{
					UserID:  cand.UserID.String,
					Title:   "Lịch phỏng vấn bị hủy",
					Message: "Lịch phỏng vấn: " + i.Title + " đã bị hủy.",
					Type:    "interview_cancelled",
					Link:    "/my-interviews",
				})
			}
		}
	}

	return nil
}

func (s *InterviewService) GenerateRoomAccessToken(ctx context.Context, interviewID, companyID, userID string) (string, string, error) {
	// Verify access + load interview for role/display resolution.
	interview, err := s.repo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return "", "", errors.NewNotFound("interview not found")
	}

	// Resolve the LiveKit room name from the room record; fall back to the
	// deterministic "room-<interviewID>" convention used elsewhere.
	roomName := "room-" + interviewID
	room, rerr := s.repo.GetRoomByInterviewID(ctx, interviewID)
	if rerr == nil && room != nil && room.RoomCode != "" {
		roomName = room.RoomCode
		log.Printf("Resolved roomName for interview %s to %s", interviewID, roomName)
	} else {
		log.Printf("Failed to resolve roomName for interview %s, rerr: %v, room: %v", interviewID, rerr, room)
	}

	// The recruiter who owns the interview gets the recruiter role; everyone
	// else joining via this endpoint is treated as a participant.
	role := "participant"
	displayName := "Participant"
	if interview.RecruiterID.Valid && interview.RecruiterID.String == userID {
		role = "recruiter"
		displayName = "Recruiter"
	}

	lkCfg := livekit.Current()
	token, err := livekit.GenerateToken(lkCfg.APIKey, lkCfg.APISecret, roomName, userID, displayName, role, interviewID)
	if err != nil {
		return "", "", errors.NewInternal("failed to generate room access token")
	}
	return token, roomName, nil
}

// canStartInterview reports whether an interview in status s may be started.
func canStartInterview(status string) bool {
	return status == "scheduled" || status == "waiting"
}

// canEndInterview reports whether an interview in status s may be ended.
func canEndInterview(status string) bool {
	return status != "completed" && status != "cancelled"
}

func (s *InterviewService) StartInterview(ctx context.Context, interviewID, companyID string) error {
	i, err := s.repo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return errors.NewNotFound("interview not found")
	}

	if !canStartInterview(i.Status) {
		return errors.NewConflict("cannot start an interview in state: " + i.Status)
	}

	// Atomic: the WHERE clause re-checks the state, so a concurrent start
	// loses cleanly instead of double-stamping.
	ok, err := s.repo.MarkStarted(ctx, interviewID, time.Now().UTC())
	if err != nil {
		return errors.NewInternal("failed to start interview")
	}
	if !ok {
		return errors.NewConflict("interview already started or not in a startable state")
	}

	_ = s.repo.MarkRoomOpened(ctx, interviewID, time.Now().UTC())
	return nil
}

func (s *InterviewService) EndInterview(ctx context.Context, interviewID, companyID string) error {
	i, err := s.repo.GetByIDAndCompany(ctx, interviewID, companyID)
	if err != nil {
		return errors.NewNotFound("interview not found")
	}

	if !canEndInterview(i.Status) {
		return errors.NewConflict("cannot end an interview in state: " + i.Status)
	}

	ok, err := s.repo.MarkEnded(ctx, interviewID, time.Now().UTC())
	if err != nil {
		return errors.NewInternal("failed to end interview")
	}
	if !ok {
		// Lost a race with another end call — report generation already
		// triggered by the winner, do not spawn a second one.
		return errors.NewConflict("interview already ended")
	}

	_ = s.repo.MarkRoomClosed(ctx, interviewID, time.Now().UTC())

	// Trigger report generation async — only on the winning transition.
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
