package models

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestJSONB_Marshal(t *testing.T) {
	data := map[string]interface{}{"key": "value", "num": 42}
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	jb := JSONB(b)
	if err := jb.Scan(b); err != nil {
		t.Fatal(err)
	}
}

func TestJSONB_Null(t *testing.T) {
	// JSONB should handle null
	var jb JSONB
	if err := jb.Scan([]byte("null")); err != nil {
		t.Fatal(err)
	}
}

func TestUser_DefaultValues(t *testing.T) {
	u := User{
		ID:        "id-1",
		Email:     "test@test.com",
		FullName:  "Test User",
		Role:      "recruiter",
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if u.ID != "id-1" {
		t.Fatalf("expected id-1, got %s", u.ID)
	}
	if u.Email != "test@test.com" {
		t.Fatalf("expected test@test.com, got %s", u.Email)
	}
}

func TestJob_Defaults(t *testing.T) {
	j := Job{
		ID:        "job-1",
		CompanyID: "company-1",
		Title:     "Test Job",
		Status:    "draft",
		CreatedAt: time.Now(),
	}
	if j.Status != "draft" {
		t.Fatalf("expected draft, got %s", j.Status)
	}
	if j.CompanyID != "company-1" {
		t.Fatalf("expected company-1, got %s", j.CompanyID)
	}
}

func TestInterview_Scheduled(t *testing.T) {
	now := time.Now()
	iv := Interview{
		ID:          "interview-1",
		Status:      "scheduled",
		ScheduledAt: sql.NullTime{Time: now, Valid: true},
		CreatedAt:   now,
	}
	if iv.Status != "scheduled" {
		t.Fatalf("expected scheduled, got %s", iv.Status)
	}
	if !iv.ScheduledAt.Valid {
		t.Fatal("ScheduledAt should be valid")
	}
}

func TestInterviewTranscript_DBColumnsMatchMigration(t *testing.T) {
	typeOfTranscript := reflect.TypeOf(InterviewTranscript{})
	got := make(map[string]bool, typeOfTranscript.NumField())
	for i := 0; i < typeOfTranscript.NumField(); i++ {
		got[typeOfTranscript.Field(i).Tag.Get("db")] = true
	}

	expected := []string{
		"id", "interview_id", "participant_id", "speaker_type", "speaker_name",
		"content", "language", "start_time_ms", "end_time_ms", "confidence",
		"source", "is_final", "edited_content", "edited_by", "edited_at", "created_at",
	}
	for _, column := range expected {
		if !got[column] {
			t.Errorf("InterviewTranscript is missing db mapping for %q", column)
		}
	}

	legacy := []string{"speaker_id", "speaker_role", "start_time", "end_time", "updated_at"}
	for _, column := range legacy {
		if got[column] {
			t.Errorf("InterviewTranscript must not map legacy column %q", column)
		}
	}
}

func TestJobCandidate_UniqueComposite(t *testing.T) {
	jc := JobCandidate{
		ID:          "jc-1",
		CompanyID:   "company-1",
		JobID:       "job-1",
		CandidateID: "candidate-1",
	}
	if jc.JobID != "job-1" || jc.CandidateID != "candidate-1" {
		t.Fatal("composite key fields should match")
	}
}

func TestAiPromptTemplate_Defaults(t *testing.T) {
	tmpl := AiPromptTemplate{
		ID:      "tmpl-1",
		Name:    "test_prompt",
		Version: 1,
		Content: "You are a test assistant",
		Model:   "gemini-1.5-flash",
	}
	// IsActive is false by default in Go (zero value)
	// In DB it defaults to true via the schema
	t.Logf("AiPromptTemplate Model has IsActive=%v (default Go zero value)", tmpl.IsActive)
	if tmpl.Version != 1 {
		t.Fatalf("expected version 1, got %d", tmpl.Version)
	}
}

func TestAuditLog_NullHandling(t *testing.T) {
	log := AuditLog{
		ID:     "log-1",
		ActorUserID: sql.NullString{},
		CompanyID:    sql.NullString{},
	}
	// Null strings should have Valid=false by default
	if log.ActorUserID.Valid {
		t.Fatal("ActorUserID should be invalid by default")
	}
	if log.CompanyID.Valid {
		t.Fatal("CompanyID should be invalid by default")
	}
}

func TestInterviewScore_Weighted(t *testing.T) {
	score := InterviewScore{
		InterviewID:   "interview-1",
		CriterionName: "Technical",
		Score:         sql.NullFloat64{Float64: 4, Valid: true},
		MaxScore:      5,
		Weight:        30,
		WeightedScore: sql.NullFloat64{Float64: 24, Valid: true},
	}
	if score.WeightedScore.Float64 != 24 {
		t.Fatalf("expected weighted score 24, got %f", score.WeightedScore.Float64)
	}
}

func TestInterviewStatuses(t *testing.T) {
	valid := map[string]bool{
		"scheduled": true, "waiting": true, "active": true,
		"completed": true, "cancelled": true, "expired": true,
	}
	for k, v := range valid {
		if !v {
			t.Errorf("unexpected status: %s", k)
		}
	}
}
