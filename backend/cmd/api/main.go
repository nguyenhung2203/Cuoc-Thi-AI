package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"backend/internal/ai"
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/pkg/email"
	"backend/internal/queue"
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

	// 2b. Connect to Redis
	redisDB := 0
	if v, convErr := strconv.Atoi(cfg.RedisDB); convErr == nil {
		redisDB = v
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.RedisPassword,
		DB:       redisDB,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Printf("redis: unavailable (%v) — notifications will not persist until redis is ready", err)
	} else {
		log.Println("redis: connected")
	}

	// 3. Repositories
	userRepo := repository.NewUserRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	jobRepo := repository.NewJobRepository(db)
	candidateRepo := repository.NewCandidateRepository(db)
	fileRepo := repository.NewFileRepository(db)
	interviewRepo := repository.NewInterviewRepository(db)
	transcriptRepo := repository.NewTranscriptRepository(db)
	systemSettingsRepo := repository.NewSystemSettingsRepository(db)
	notificationRepo := repository.NewNotificationRepository(redisClient, systemSettingsRepo)
	reportRepo := repository.NewReportRepository(db)
	mockRepo := repository.NewMockRepository(db)
	candidatePortalRepo := repository.NewCandidatePortalRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	rubricRepo := repository.NewRubricRepository(db.DB)
	questionRepo := repository.NewQuestionRepository(db.DB)
	interviewTemplateRepo := repository.NewInterviewTemplateRepository(db.DB)
	aiPromptRepo := repository.NewAIPromptRepository(db)
	aiLogRepo := repository.NewAILogRepository(db)
	scoreRepo := repository.NewScoreRepository(db.DB)

	// AI Setup
	promptSvc := service.NewPromptService(aiPromptRepo)
	aiLogSvc := service.NewAILogService(aiLogRepo)
	aiOrchestrator := service.NewAIOrchestratorService(promptSvc, aiLogSvc, cfg.AIServiceURL)

	jdAnalyzer := ai.NewJDAnalyzer(aiOrchestrator)
	cvAnalyzer := ai.NewCVAnalyzer(aiOrchestrator)
	qGenerator := ai.NewQuestionGenerator(aiOrchestrator)

	// 4. Services
	auditRepo := repository.NewAuditRepository(db)
	auditSvc := service.NewAuditService(auditRepo)
	authSvc := service.NewAuthService(userRepo, refreshTokenRepo, cfg.JWTSecret)
	companySvc := service.NewCompanyService(companyRepo)
	jobSvc := service.NewJobService(jobRepo, jdAnalyzer, qGenerator, rubricRepo, questionRepo)
	candidateSvc := service.NewCandidateService(candidateRepo, jobRepo, cvAnalyzer)
	fileSvc := service.NewFileService(fileRepo, cfg.PublicBaseURL)
	transcriptSvc := service.NewTranscriptService(transcriptRepo, interviewRepo)
	rubricSvc := service.NewRubricService(rubricRepo)
	scoreSvc := service.NewScoreService(scoreRepo, transcriptRepo, rubricRepo, interviewRepo, aiOrchestrator)
	reportSvc := service.NewReportService(reportRepo, transcriptRepo, scoreRepo, jobRepo, interviewRepo, notificationRepo, aiOrchestrator)
	mailer := email.NewSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)
	otpSvc := service.NewOTPService(redisClient)
	authSvc.WithOTP(otpSvc, mailer)
	authSvc.WithGoogle(cfg.GoogleClientID)
	interviewSvc := service.NewInterviewService(interviewRepo, reportSvc, notificationRepo).
		WithMailer(candidateRepo, mailer, cfg.FrontendURL)
	suggestionSvc := service.NewSuggestionService(aiOrchestrator, transcriptRepo, interviewRepo, jobRepo)
	mockSvc := service.NewMockService(mockRepo, aiOrchestrator, promptSvc)
	aiSvc := service.NewAIService(fileRepo, candidateRepo, cfg.GeminiAPIKey)
	matchCache := service.NewMatchCacheService(redisClient)

	candidatePortalSvc := service.NewCandidatePortalService(candidatePortalRepo, userRepo, candidateRepo, jobRepo, aiSvc, aiOrchestrator, matchCache, notificationRepo, fileSvc)
	userSvc := service.NewUserService(userRepo, refreshTokenRepo)

	// 4b. Async queue (Redis/asynq). Best-effort: if Redis is unavailable the
	// platform still runs, with heavy jobs processed inline instead.
	if dispatcher, derr := queue.NewDispatcher(cfg.RedisAddr(), cfg.RedisPassword, redisDB); derr != nil {
		log.Printf("queue: Redis unavailable (%v) — reports will run inline", derr)
	} else {
		reportSvc.SetEnqueuer(dispatcher.EnqueueGenerateReport)
		candidatePortalSvc.SetMatchEnqueuer(dispatcher.EnqueueRecomputeMatches)
		worker := queue.NewWorker(cfg.RedisAddr(), cfg.RedisPassword, redisDB, 10, reportSvc, aiSvc, candidatePortalSvc)
		go func() {
			log.Println("queue: async worker started")
			if werr := worker.Run(); werr != nil {
				log.Printf("queue: worker stopped: %v", werr)
			}
		}()
	}

	// 5. Handlers
	authHandler := handler.NewAuthHandler(authSvc, cfg.JWTSecret, auditSvc)
	userHandler := handler.NewUserHandler(userSvc, companySvc, auditSvc, cfg.JWTSecret)
	companyHandler := handler.NewCompanyHandler(companySvc)
	jobHandler := handler.NewJobHandler(jobSvc)
	candidateHandler := handler.NewCandidateHandler(candidateSvc, aiSvc, fileSvc)
	fileHandler := handler.NewFileHandler(fileSvc)
	interviewHandler := handler.NewInterviewHandler(interviewSvc)
	mockHandler := handler.NewMockHandler(mockSvc)
	transcriptHandler := handler.NewTranscriptHandler(transcriptSvc)
	reportHandler := handler.NewReportHandler(reportSvc)
	aiAdminHandler := handler.NewAIAdminHandler(promptSvc, aiLogSvc)
	notificationHandler := handler.NewNotificationHandler(notificationRepo)
	candidatePortalHandler := handler.NewCandidatePortalHandler(candidatePortalSvc, fileSvc)
	rubricHandler := handler.NewRubricHandler(rubricSvc)

	questionBankSvc := service.NewQuestionBankService(questionRepo)
	questionBankHandler := handler.NewQuestionBankHandler(questionBankSvc)

	interviewTemplateSvc := service.NewInterviewTemplateService(interviewTemplateRepo)
	interviewTemplateHandler := handler.NewInterviewTemplateHandler(interviewTemplateSvc)

	aiHandler := handler.NewAiHandler(scoreSvc, reportSvc, suggestionSvc)
	auditHandler := handler.NewAuditHandler(auditSvc)

	systemSettingsSvc := service.NewSystemSettingsService(systemSettingsRepo)
	systemSettingsHandler := handler.NewSystemSettingsHandler(systemSettingsSvc, auditSvc)

	// 6. Router
	r := chi.NewRouter()

	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.RecoveryMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"healthy"}`))
	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "uploads"))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(filesDir)))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Use(middleware.AuthRateLimitMiddleware)
			authHandler.Routes(r)
		})

		r.Route("/public", func(r chi.Router) {
			publicJobHandler := handler.NewPublicJobHandler(jobSvc)
			publicJobHandler.Routes(r)
			r.Get("/companies/{company_id}", companyHandler.Get)
		})

		r.Route("/interviews", func(r chi.Router) {
			interviewHandler.Routes(r)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			r.Use(middleware.AuditMiddleware(auditSvc))

			r.Route("/mock-interviews", func(r chi.Router) {
				r.Post("/", mockHandler.Create)
				r.Get("/me", mockHandler.ListMine)
				r.Get("/{id}", mockHandler.GetByID)
				r.Post("/{id}/start", mockHandler.Start)
				r.Post("/{id}/end", mockHandler.End)
				r.Post("/{id}/messages", mockHandler.SendMessage)
				r.Get("/{id}/messages", mockHandler.GetMessages)
				r.Post("/{id}/save-live-transcript", mockHandler.SaveLiveTranscript)
				r.Get("/{id}/report", mockHandler.GetReport)
			})

			r.Route("/notifications", func(r chi.Router) {
				notificationHandler.ProtectedRoutes(r)
			})

			r.Route("/admin", func(r chi.Router) {
				r.Use(middleware.RoleMiddleware("admin"))
				r.Get("/ai-prompts", aiAdminHandler.ListTemplates)
				r.Post("/ai-prompts", aiAdminHandler.CreatePromptTemplate)
				r.Get("/ai-logs", aiAdminHandler.ListAILogs)
				r.Get("/logs", auditHandler.ListAllGlobalLogs)
				r.Get("/settings", systemSettingsHandler.GetSettings)
				r.Put("/settings", systemSettingsHandler.UpdateSettings)
			})
			userHandler.Routes(r)

			r.Route("/portal", func(r chi.Router) {
				candidatePortalHandler.Routes(r)
			})

			fileHandler.Routes(r)

			r.Route("/companies", func(r chi.Router) {
				companyHandler.Routes(r)
			})

			r.Route("/companies/{company_id}", func(r chi.Router) {
				r.Use(middleware.CompanyScopeMiddleware(db))

				jobHandler.Routes(r)
				candidateHandler.Routes(r)

				r.Get("/audit-logs", auditHandler.ListAuditLogs)

				r.Route("/rubrics", func(r chi.Router) {
					r.Post("/", rubricHandler.CreateRubric)
					r.Get("/", rubricHandler.ListCompanyRubrics)
					r.Get("/{rubric_id}", rubricHandler.GetRubric)
					r.Put("/{rubric_id}", rubricHandler.UpdateRubric)
					r.Delete("/{rubric_id}", rubricHandler.DeleteRubric)
				})
				questionBankHandler.Routes(r)
				interviewTemplateHandler.Routes(r)

				r.Route("/interviews", func(r chi.Router) {
					interviewHandler.ProtectedRoutes(r)

					r.Route("/{interview_id}/transcripts", func(r chi.Router) {
						transcriptHandler.Routes(r)
					})
					r.Route("/{interview_id}/report", func(r chi.Router) {
						r.Get("/", reportHandler.GetReport)
						r.Put("/decision", reportHandler.OverrideDecision)
						r.Post("/retry", reportHandler.RetryReport)
					})
					r.Route("/{interview_id}/ai", func(r chi.Router) {
						r.Post("/score-answer", aiHandler.ScoreAnswer)
						r.Post("/generate-report", aiHandler.GenerateReport)
						r.Post("/suggest-follow-up", aiHandler.SuggestFollowUp)
					})
				})
			})
		})
	})

	addr := ":" + cfg.AppPort
	log.Printf("server: listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
