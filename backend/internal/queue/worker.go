package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"backend/internal/pkg/email"
	"backend/internal/pkg/logger"
)

type OTPMailer interface {
	SendHTML(to, subject, textBody, htmlBody string) error
}

type OTPReader interface {
	Current(ctx context.Context, purpose, email string) (code, generationID string, err error)
}

type ReportRunner interface {
	RunReportGeneration(ctx context.Context, companyID, interviewID, jobID, generatedBy, recruiterID string) error
}

// CVAnalyzer is satisfied by *service.AIService.ParseCV.
type CVAnalyzer interface {
	ParseCV(ctx context.Context, companyID, candidateID string) error
}

// MatchRecomputer is satisfied by *service.CandidatePortalService.RecomputeUserMatches.
type MatchRecomputer interface {
	RecomputeUserMatches(ctx context.Context, userID string) error
}

// Worker consumes async jobs from Redis and dispatches them to services.
type Worker struct {
	server     *asynq.Server
	reportSvc  ReportRunner
	cvAnalyzer CVAnalyzer
	matchSvc   MatchRecomputer
	mailer     OTPMailer
	otpReader  OTPReader
}

// NewWorker builds an asynq server. concurrency bounds parallel job execution.
func NewWorker(redisAddr, password string, db, concurrency int, reportSvc ReportRunner, cvAnalyzer CVAnalyzer, matchSvc MatchRecomputer, mailer OTPMailer, otpReader OTPReader) *Worker {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr, Password: password, DB: db},
		asynq.Config{
			Concurrency: concurrency,
			Logger:      &asynqLogger{},
		},
	)
	return &Worker{server: server, reportSvc: reportSvc, cvAnalyzer: cvAnalyzer, matchSvc: matchSvc, mailer: mailer, otpReader: otpReader}
}

// Run registers handlers and blocks serving jobs until the process stops.
func (w *Worker) Run() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeGenerateReport, w.handleGenerateReport)
	mux.HandleFunc(TypeAnalyzeCV, w.handleAnalyzeCV)
	mux.HandleFunc(TypeBatchTranscript, w.handleBatchTranscript)
	mux.HandleFunc(TypeRecomputeMatches, w.handleRecomputeMatches)
	mux.HandleFunc(TypeSendOTPEmail, w.handleSendOTPEmail)
	return w.server.Run(mux)
}

// Shutdown stops the worker gracefully.
func (w *Worker) Shutdown() { w.server.Shutdown() }

func (w *Worker) handleSendOTPEmail(ctx context.Context, t *asynq.Task) error {
	var p SendOTPEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil || p.To == "" || p.Purpose == "" || p.GenerationID == "" {
		logger.Error("[Worker Failure] invalid OTP payload: %v", err)
		return fmt.Errorf("bad otp email payload: %w", asynq.SkipRetry)
	}
	logger.Info("[Worker Pickup Task] email=%s purpose=%s generation=%s", p.To, p.Purpose, p.GenerationID)
	if w.mailer == nil || w.otpReader == nil {
		logger.Error("[Worker Failure] OTP dependencies not configured")
		return fmt.Errorf("otp dependencies not configured: %w", asynq.SkipRetry)
	}
	code, generationID, err := w.otpReader.Current(ctx, p.Purpose, p.To)
	if err != nil {
		logger.Error("[Read OTP Redis Failure] email=%s purpose=%s error=%v", p.To, p.Purpose, err)
		return fmt.Errorf("read current otp: %w", err)
	}
	if generationID != p.GenerationID {
		logger.Warn("[Stale OTP] email=%s purpose=%s task_generation=%s current_generation=%s", p.To, p.Purpose, p.GenerationID, generationID)
		return fmt.Errorf("stale otp email job: %w", asynq.SkipRetry)
	}
	logger.Info("[Read OTP Redis Success] email=%s purpose=%s", p.To, p.Purpose)
	textBody, htmlBody := email.RenderOTP(email.OTPTemplateData{Purpose: p.TemplateType, Code: code, ExpiryMinute: 10})
	logger.Info("[SMTP Dialing] email=%s purpose=%s", p.To, p.Purpose)
	if err := w.mailer.SendHTML(p.To, "Mã xác nhận ViệcLàm AI", textBody, htmlBody); err != nil {
		logger.Error("[SMTP Send Mail Failure] email=%s purpose=%s error=%v", p.To, p.Purpose, err)
		return fmt.Errorf("otp email delivery failed: %w", err)
	}
	logger.Info("[SMTP Send Mail Success] email=%s purpose=%s", p.To, p.Purpose)
	return nil
}

func (w *Worker) handleGenerateReport(ctx context.Context, t *asynq.Task) error {
	var p GenerateReportPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		// Non-retryable: malformed payload.
		return fmt.Errorf("bad report payload: %v: %w", err, asynq.SkipRetry)
	}
	if w.reportSvc == nil {
		return fmt.Errorf("report runner not configured: %w", asynq.SkipRetry)
	}
	// RunReportGeneration sets terminal status/notifications itself. Returning
	// an error lets asynq retry (up to MaxRetry) with backoff.
	if err := w.reportSvc.RunReportGeneration(ctx, p.CompanyID, p.InterviewID, p.JobID, p.GeneratedBy, p.RecruiterID); err != nil {
		return fmt.Errorf("report generation failed for %s: %w", p.InterviewID, err)
	}
	return nil
}

func (w *Worker) handleAnalyzeCV(ctx context.Context, t *asynq.Task) error {
	var p AnalyzeCVPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("bad cv payload: %v: %w", err, asynq.SkipRetry)
	}
	if w.cvAnalyzer == nil {
		return fmt.Errorf("cv analyzer not configured: %w", asynq.SkipRetry)
	}
	if err := w.cvAnalyzer.ParseCV(ctx, p.CompanyID, p.CandidateID); err != nil {
		return fmt.Errorf("cv analysis failed for %s: %w", p.CandidateID, err)
	}
	return nil
}

func (w *Worker) handleRecomputeMatches(ctx context.Context, t *asynq.Task) error {
	var p RecomputeMatchesPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("bad recompute-matches payload: %v: %w", err, asynq.SkipRetry)
	}
	if w.matchSvc == nil {
		return fmt.Errorf("match recomputer not configured: %w", asynq.SkipRetry)
	}
	if err := w.matchSvc.RecomputeUserMatches(ctx, p.UserID); err != nil {
		return fmt.Errorf("match recompute failed for user %s: %w", p.UserID, err)
	}
	return nil
}

func (w *Worker) handleBatchTranscript(ctx context.Context, t *asynq.Task) error {
	// Transcripts are persisted synchronously by the realtime gateway; this
	// hook exists for future batch post-processing (e.g. re-diarization).
	logger.Info("batch transcript job received (no-op)")
	return nil
}

// asynqLogger adapts the project logger to asynq's Logger interface.
type asynqLogger struct{}

func (l *asynqLogger) Debug(args ...interface{}) {}
func (l *asynqLogger) Info(args ...interface{})  { logger.Info(fmt.Sprint(args...)) }
func (l *asynqLogger) Warn(args ...interface{})  { logger.Warn(fmt.Sprint(args...)) }
func (l *asynqLogger) Error(args ...interface{}) { logger.Error(fmt.Sprint(args...)) }
func (l *asynqLogger) Fatal(args ...interface{}) { logger.Error(fmt.Sprint(args...)) }
