package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"backend/internal/pkg/logger"
)

// ReportRunner is satisfied by *service.ReportService.RunReportGeneration.
type ReportRunner interface {
	RunReportGeneration(ctx context.Context, companyID, interviewID, jobID, generatedBy, recruiterID string) error
}

// CVAnalyzer is satisfied by *service.AIService.ParseCV.
type CVAnalyzer interface {
	ParseCV(ctx context.Context, companyID, candidateID string) error
}

// Worker consumes async jobs from Redis and dispatches them to services.
type Worker struct {
	server     *asynq.Server
	reportSvc  ReportRunner
	cvAnalyzer CVAnalyzer
}

// NewWorker builds an asynq server. concurrency bounds parallel job execution.
func NewWorker(redisAddr, password string, db, concurrency int, reportSvc ReportRunner, cvAnalyzer CVAnalyzer) *Worker {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr, Password: password, DB: db},
		asynq.Config{
			Concurrency: concurrency,
			Logger:      &asynqLogger{},
		},
	)
	return &Worker{server: server, reportSvc: reportSvc, cvAnalyzer: cvAnalyzer}
}

// Run registers handlers and blocks serving jobs until the process stops.
func (w *Worker) Run() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeGenerateReport, w.handleGenerateReport)
	mux.HandleFunc(TypeAnalyzeCV, w.handleAnalyzeCV)
	mux.HandleFunc(TypeBatchTranscript, w.handleBatchTranscript)
	return w.server.Run(mux)
}

// Shutdown stops the worker gracefully.
func (w *Worker) Shutdown() { w.server.Shutdown() }

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
