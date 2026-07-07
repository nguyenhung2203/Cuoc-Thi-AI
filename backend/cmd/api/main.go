package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"backend/internal/ai"
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/service"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// 2. Connect to PostgreSQL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
		)
	}
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()
	log.Println("database: connected")

	// 3. Repositories
	userRepo := repository.NewUserRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	jobRepo := repository.NewJobRepository(db)
	candidateRepo := repository.NewCandidateRepository(db)
	fileRepo := repository.NewFileRepository(db)
	interviewRepo := repository.NewInterviewRepository(db)
	transcriptRepo := repository.NewTranscriptRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	reportRepo := repository.NewReportRepository(db)
	mockRepo := repository.NewMockInterviewRepository(db)
	candidatePortalRepo := repository.NewCandidatePortalRepository(db)


	// 4. Services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	companySvc := service.NewCompanyService(companyRepo)
	jobSvc := service.NewJobService(jobRepo)
	candidateSvc := service.NewCandidateService(candidateRepo, jobRepo)
	fileSvc := service.NewFileService(fileRepo)
	interviewSvc := service.NewInterviewService(interviewRepo)
	transcriptSvc := service.NewTranscriptService(transcriptRepo, interviewRepo)
	aiSvc := service.NewAIService(fileRepo, candidateRepo, cfg.GeminiAPIKey)
	notificationSvc := service.NewNotificationService(notificationRepo)
	reportSvc := service.NewReportService(reportRepo, transcriptSvc, aiSvc)
	mockSvc := service.NewMockService(mockRepo, aiSvc)
	candidatePortalSvc := service.NewCandidatePortalService(candidatePortalRepo, userRepo, candidateRepo, jobRepo)

	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	rubricRepo := repository.NewRubricRepository(db.DB)
	questionRepo := repository.NewQuestionRepository(db.DB)


	// AI Setup
	aiPromptRepo := repository.NewAIPromptRepository(db)
	aiLogRepo := repository.NewAILogRepository(db)
	scoreRepo := repository.NewScoreRepository(db.DB)
	reportRepo := repository.NewReportRepository(db)

	// 3. AI Providers & Internal Services
	promptSvc := service.NewPromptService(aiPromptRepo)
	aiLogSvc := service.NewAILogService(aiLogRepo)
	aiOrchestrator := service.NewAIOrchestratorService(promptSvc, aiLogSvc, cfg.AIServiceURL)

	jdAnalyzer := ai.NewJDAnalyzer(aiOrchestrator)
	cvAnalyzer := ai.NewCVAnalyzer(aiOrchestrator)
	qGenerator := ai.NewQuestionGenerator(aiOrchestrator)

	// 4. Services
	authSvc := service.NewAuthService(userRepo, refreshTokenRepo, cfg.JWTSecret)
	companySvc := service.NewCompanyService(companyRepo)
	jobSvc := service.NewJobService(jobRepo, jdAnalyzer, qGenerator, rubricRepo, questionRepo)
	candidateSvc := service.NewCandidateService(candidateRepo, jobRepo, cvAnalyzer)
	fileSvc := service.NewFileService(fileRepo)
	transcriptSvc := service.NewTranscriptService(transcriptRepo, interviewRepo)
	rubricSvc := service.NewRubricService(rubricRepo)
	scoreSvc := service.NewScoreService(scoreRepo, transcriptRepo, rubricRepo, interviewRepo, aiOrchestrator)
	notifRepo := repository.NewNotificationRepository(db)
	mockRepo := repository.NewMockRepository(db)
	reportSvc := service.NewReportService(reportRepo, transcriptRepo, scoreRepo, jobRepo, interviewRepo, notifRepo, aiOrchestrator)
	interviewSvc := service.NewInterviewService(interviewRepo, reportSvc)
	mockSvc := service.NewMockService(mockRepo, aiOrchestrator, promptSvc)

	// 5. Handlers
	authHandler := handler.NewAuthHandler(authSvc, cfg.JWTSecret)
	companyHandler := handler.NewCompanyHandler(companySvc)
	jobHandler := handler.NewJobHandler(jobSvc)
	candidateHandler := handler.NewCandidateHandler(candidateSvc, aiSvc)
	fileHandler := handler.NewFileHandler(fileSvc)
	interviewHandler := handler.NewInterviewHandler(interviewSvc)
	mockHandler := handler.NewMockHandler(mockSvc)
	transcriptHandler := handler.NewTranscriptHandler(transcriptSvc)
	reportHandler := handler.NewReportHandler(reportSvc)
	aiAdminHandler := handler.NewAIAdminHandler(promptSvc)
	notificationHandler := handler.NewNotificationHandler(notificationSvc)
	mockHandler := handler.NewMockHandler(mockSvc)
	candidatePortalHandler := handler.NewCandidatePortalHandler(candidatePortalSvc, fileSvc)
	rubricHandler := handler.NewRubricHandler(rubricSvc)
	aiHandler := handler.NewAiHandler(scoreSvc, reportSvc)
	reportHandler := handler.NewReportHandler(reportSvc)
	notifHandler := handler.NewNotificationHandler(notifRepo)


	// In a real app, aiOrchestrator would be injected into handlers that need it (e.g. JobHandler for AnalyzeJD)

	// 6. Router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.RecoveryMiddleware)

	// Health check (no auth)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"healthy"}`))
	})

	// Serve static uploads
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "uploads"))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(filesDir)))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Use(middleware.AuthRateLimitMiddleware)
			authHandler.Routes(r)
		})

		// Public career site routes
		r.Route("/public", func(r chi.Router) {
			publicJobHandler := handler.NewPublicJobHandler(jobSvc)
			publicJobHandler.Routes(r)
		})

		// Public interview routes
		r.Route("/interviews", func(r chi.Router) {
			interviewHandler.Routes(r)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))

			// Admin routes

				// Mock interview routes
				r.Route("/mock-interviews", func(r chi.Router) {
					r.Post("/", mockHandler.Create)
					r.Get("/me", mockHandler.ListMine)
					r.Get("/{id}", mockHandler.GetByID)
					r.Post("/{id}/start", mockHandler.Start)
					r.Post("/{id}/end", mockHandler.End)
					r.Post("/{id}/messages", mockHandler.SendMessage)
					r.Get("/{id}/report", mockHandler.GetReport)
				})

		// Notification routes — user's own notifications
			r.Get("/notifications", notifHandler.ListNotifications)
			r.Put("/notifications/{notification_id}/read", notifHandler.MarkNotificationRead)

			r.Route("/admin", func(r chi.Router) {
				// Require admin role middleware would normally go here
				r.Post("/ai-prompts", aiAdminHandler.CreatePromptTemplate)
			})

			// Candidate Portal routes
			r.Route("/portal", func(r chi.Router) {
				candidatePortalHandler.Routes(r)
			})

			// File routes — user accesses own files directly (not company-scoped)
			fileHandler.Routes(r)

			// Notifications (user scoped)
			r.Route("/notifications", func(r chi.Router) {
				notificationHandler.ProtectedRoutes(r)
			})

			// Mock Interviews (candidate scoped)
			r.Route("/mock-interviews", func(r chi.Router) {
				mockHandler.ProtectedRoutes(r)
			})

			// Company management routes
			r.Route("/companies", func(r chi.Router) {
				companyHandler.Routes(r)
			})

			// Company-scoped routes
			r.Route("/companies/{company_id}", func(r chi.Router) {
				r.Use(middleware.CompanyScopeMiddleware(db))

				jobHandler.Routes(r)
				candidateHandler.Routes(r)

				r.Route("/rubrics", func(r chi.Router) {
					r.Post("/", rubricHandler.CreateRubric)
					r.Get("/", rubricHandler.ListCompanyRubrics)
					r.Get("/{rubric_id}", rubricHandler.GetRubric)
					r.Delete("/{rubric_id}", rubricHandler.DeleteRubric)
				})

				r.Route("/interviews", func(r chi.Router) {
					interviewHandler.ProtectedRoutes(r)

					r.Route("/{interview_id}/transcripts", func(r chi.Router) {
						transcriptHandler.Routes(r)
					})
					r.Route("/{interview_id}/report", func(r chi.Router) {
						r.Get("/", reportHandler.GetReport)
						r.Put("/decision", reportHandler.OverrideDecision)
						r.Post("/retry", reportHandler.RetryReport)
						reportHandler.Routes(r)
					})
					r.Route("/{interview_id}/ai", func(r chi.Router) {
						r.Post("/score-answer", aiHandler.ScoreAnswer)
						r.Post("/generate-report", aiHandler.GenerateReport)
					})
				})
			})
		})
	})

	// 7. Start server
	addr := ":" + cfg.AppPort
	log.Printf("server: listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
