package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

// Dispatcher enqueues async jobs onto Redis via asynq.
type Dispatcher struct {
	client *asynq.Client
}

// NewDispatcher creates a dispatcher connected to Redis. A quick PING verifies
// connectivity so callers can fall back to inline processing when Redis is down.
func NewDispatcher(redisAddr, password string, db int) (*Dispatcher, error) {
	opt := asynq.RedisClientOpt{Addr: redisAddr, Password: password, DB: db}
	client := asynq.NewClient(opt)

	// Verify Redis is reachable; asynq.Client is lazy otherwise.
	inspector := asynq.NewInspector(opt)
	defer inspector.Close()
	if _, err := inspector.Queues(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis unreachable at %s: %w", redisAddr, err)
	}

	return &Dispatcher{client: client}, nil
}

func (d *Dispatcher) Close() error {
	if d.client == nil {
		return nil
	}
	return d.client.Close()
}

// defaultOpts: retry up to 3 times (per spec), with asynq's exponential
// backoff, and a per-task timeout.
func defaultOpts(timeout time.Duration) []asynq.Option {
	return []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Timeout(timeout),
	}
}

func (d *Dispatcher) enqueue(taskType string, payload interface{}, opts ...asynq.Option) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	task := asynq.NewTask(taskType, b)
	_, err = d.client.Enqueue(task, opts...)
	return err
}

// EnqueueSendOTPEmail schedules OTP delivery outside the request path.
func (d *Dispatcher) EnqueueSendOTPEmail(payload SendOTPEmailPayload) error {
	err := d.enqueue(TypeSendOTPEmail, payload, asynq.MaxRetry(3), asynq.Timeout(30*time.Second), asynq.Retention(10*time.Minute))
	if err == nil {
		log.Printf("[Enqueue Success] type=%s email=%s purpose=%s generation=%s", TypeSendOTPEmail, payload.To, payload.Purpose, payload.GenerationID)
	}
	return err
}

func (d *Dispatcher) EnqueueGenerateReport(companyID, interviewID, jobID, generatedBy, recruiterID string) error {
	return d.enqueue(TypeGenerateReport, GenerateReportPayload{
		CompanyID:   companyID,
		InterviewID: interviewID,
		JobID:       jobID,
		GeneratedBy: generatedBy,
		RecruiterID: recruiterID,
	}, defaultOpts(10*time.Minute)...)
}

// EnqueueAnalyzeCV schedules async CV analysis.
func (d *Dispatcher) EnqueueAnalyzeCV(companyID, candidateID string) error {
	return d.enqueue(TypeAnalyzeCV, AnalyzeCVPayload{
		CompanyID:   companyID,
		CandidateID: candidateID,
	}, defaultOpts(5*time.Minute)...)
}

// EnqueueBatchTranscript schedules async transcript batch processing.
func (d *Dispatcher) EnqueueBatchTranscript(companyID, interviewID string) error {
	return d.enqueue(TypeBatchTranscript, BatchTranscriptPayload{
		CompanyID:   companyID,
		InterviewID: interviewID,
	}, defaultOpts(5*time.Minute)...)
}

// EnqueueRecomputeMatches schedules an async recompute of all of a user's
// application fit scores (after a CV change).
func (d *Dispatcher) EnqueueRecomputeMatches(userID string) error {
	return d.enqueue(TypeRecomputeMatches, RecomputeMatchesPayload{
		UserID: userID,
	}, defaultOpts(10*time.Minute)...)
}
