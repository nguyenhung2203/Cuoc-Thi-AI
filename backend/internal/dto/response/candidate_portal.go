package response

import "time"

type CandidatePortalDashboardStats struct {
	UpcomingInterviews int     `json:"upcoming_interviews"`
	CompletedMockTests int     `json:"completed_mock_tests"`
	AverageMockScore   float64 `json:"average_mock_score"`
	ProfileCompleteness int    `json:"profile_completeness"` // percentage 0-100
}

type CandidatePortalInterview struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	CompanyID       string     `json:"company_id"`
	CompanyName     string     `json:"company_name"`
	JobID           string     `json:"job_id"`
	JobTitle        string     `json:"job_title"`
	Mode            string     `json:"mode"`
	Status          string     `json:"status"`
	ScheduledAt     *time.Time `json:"scheduled_at"`
	InviteTokenHash string     `json:"-"`
	JoinLink        string     `json:"join_link,omitempty"`
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
