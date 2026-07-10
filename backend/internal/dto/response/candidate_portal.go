package response

import "time"

type CandidatePortalDashboardStats struct {
	UpcomingInterviews int     `json:"upcoming_interviews"`
	CompletedMockTests int     `json:"completed_mock_tests"`
	AverageMockScore   float64 `json:"average_mock_score"`
	ProfileCompleteness int    `json:"profile_completeness"` // percentage 0-100
}

type CandidatePortalInterview struct {
	ID              string     `json:"id" db:"id"`
	Title           string     `json:"title" db:"title"`
	CompanyID       string     `json:"company_id" db:"company_id"`
	CompanyName     string     `json:"company_name" db:"company_name"`
	JobID           string     `json:"job_id" db:"job_id"`
	JobTitle        string     `json:"job_title" db:"job_title"`
	Mode            string     `json:"mode" db:"mode"`
	Status          string     `json:"status" db:"status"`
	ScheduledAt     *time.Time `json:"scheduled_at" db:"scheduled_at"`
	InviteTokenHash string     `json:"-" db:"invite_token_hash"`
	JoinLink        string     `json:"join_link,omitempty" db:"-"`
}

type CandidatePortalProfile struct {
	UserID    string `json:"user_id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
	CVFileID  string `json:"cv_file_id"`
	CVUrl     string `json:"cv_url"`
	CVName    string `json:"cv_name"`
}

type CandidateApplication struct {
	JobID          string    `json:"job_id" db:"job_id"`
	JobTitle       string    `json:"job_title" db:"job_title"`
	CompanyID      string    `json:"company_id" db:"company_id"`
	CompanyName    string    `json:"company_name" db:"company_name"`
	PipelineStatus string    `json:"pipeline_status" db:"pipeline_status"`
	AppliedAt      time.Time `json:"applied_at" db:"applied_at"`
}
