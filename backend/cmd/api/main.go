package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

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

	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	// 4. Services
	authSvc := service.NewAuthService(userRepo, refreshTokenRepo, cfg.JWTSecret)
	companySvc := service.NewCompanyService(companyRepo)
	jobSvc := service.NewJobService(jobRepo)
	candidateSvc := service.NewCandidateService(candidateRepo, jobRepo)
	fileSvc := service.NewFileService(fileRepo)
	interviewSvc := service.NewInterviewService(interviewRepo)
	transcriptSvc := service.NewTranscriptService(transcriptRepo, interviewRepo)

	aiPromptRepo := repository.NewAIPromptRepository(db)
	aiLogRepo := repository.NewAILogRepository(db)
	promptSvc := service.NewPromptService(aiPromptRepo)
	aiLogSvc := service.NewAILogService(aiLogRepo)
	_ = service.NewAIOrchestratorService(promptSvc, aiLogSvc, cfg.AIServiceURL)

	// 5. Handlers
	authHandler := handler.NewAuthHandler(authSvc, cfg.JWTSecret)
	companyHandler := handler.NewCompanyHandler(companySvc)
	jobHandler := handler.NewJobHandler(jobSvc)
	candidateHandler := handler.NewCandidateHandler(candidateSvc)
	fileHandler := handler.NewFileHandler(fileSvc)
	interviewHandler := handler.NewInterviewHandler(interviewSvc)
	transcriptHandler := handler.NewTranscriptHandler(transcriptSvc)
	aiAdminHandler := handler.NewAIAdminHandler(promptSvc)

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

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Use(middleware.AuthRateLimitMiddleware)
			authHandler.Routes(r)
		})

		// Public interview routes
		r.Route("/interviews", func(r chi.Router) {
			interviewHandler.Routes(r)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))

			// Admin routes
			r.Route("/admin", func(r chi.Router) {
				// Require admin role middleware would normally go here
				r.Post("/ai-prompts", aiAdminHandler.CreatePromptTemplate)
			})

			// File routes — user accesses own files directly (not company-scoped)
			fileHandler.Routes(r)

			// Company management routes
			r.Route("/companies", func(r chi.Router) {
				companyHandler.Routes(r)
			})

			// Company-scoped routes
			r.Route("/companies/{company_id}", func(r chi.Router) {
				r.Use(middleware.CompanyScopeMiddleware(db))

				jobHandler.Routes(r)
				candidateHandler.Routes(r)
				
				r.Route("/interviews", func(r chi.Router) {
					interviewHandler.ProtectedRoutes(r)
					
					r.Route("/{interview_id}/transcripts", func(r chi.Router) {
						transcriptHandler.Routes(r)
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
