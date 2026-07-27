package response

import "time"

// InterviewResponse is the API-facing shape of an interview.
type InterviewResponse struct {
	ID          string     `json:"id"`
	CompanyID   string     `json:"company_id"`
	JobID       string     `json:"job_id,omitempty"`
	CandidateID string     `json:"candidate_id,omitempty"`
	RecruiterID string     `json:"recruiter_id,omitempty"`
	Title       string     `json:"title"`
	Mode        string     `json:"mode"`
	Status      string     `json:"status"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	EndedAt     *time.Time `json:"ended_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// RoomAccessTokenResponse is returned by the room access-token endpoints.
type RoomAccessTokenResponse struct {
	Token    string `json:"token"`
	RoomCode string `json:"room_code"`
	Role     string `json:"role"`
}
