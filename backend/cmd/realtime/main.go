package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"backend/internal/config"
	"backend/internal/realtime"
	"backend/internal/repository"
	"backend/internal/service"
)

func main() {
	// Load .env so DATABASE_URL / GEMINI_API_KEY / LIVEKIT_* are available.
	_ = godotenv.Load()
	_ = godotenv.Load("../../.env")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// LiveKit config: fail fast in production if dev fallbacks would be used.
	livekitCfg, err := realtime.LoadLiveKitConfig()
	if err != nil {
		log.Fatalf("livekit config: %v", err)
	}

	// Connect to PostgreSQL. Fail fast: the gateway now persists transcripts,
	// chat, audit logs, and resolves room/invite tokens against the DB.
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

	// Build the real service dependencies the router/handlers call.
	interviewRepo := repository.NewInterviewRepository(db)
	transcriptRepo := repository.NewTranscriptRepository(db)
	jobRepo := repository.NewJobRepository(db)
	rubricRepo := repository.NewRubricRepository(db.DB)
	scoreRepo := repository.NewScoreRepository(db.DB)
	aiPromptRepo := repository.NewAIPromptRepository(db)
	aiLogRepo := repository.NewAILogRepository(db)

	promptSvc := service.NewPromptService(aiPromptRepo)
	aiLogSvc := service.NewAILogService(aiLogRepo)
	aiOrchestrator := service.NewAIOrchestratorService(promptSvc, aiLogSvc, cfg.AIServiceURL)

	suggestionSvc := service.NewSuggestionService(aiOrchestrator, transcriptRepo, interviewRepo, jobRepo)
	scoreSvc := service.NewScoreService(scoreRepo, transcriptRepo, rubricRepo, interviewRepo, aiOrchestrator)

	deps := realtime.RouterDeps{
		InterviewRepo:  interviewRepo,
		TranscriptRepo: transcriptRepo,
		SuggestionSvc:  suggestionSvc,
		ScoreSvc:       scoreSvc,
	}

	port := os.Getenv("REALTIME_PORT")
	if port == "" {
		port = "8081"
	}
	addr := ":" + port
	fmt.Printf("Starting AI Interview Platform Realtime Gateway on %s...\n", addr)

	srv := realtime.NewServer(addr, db, deps, livekitCfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Received termination signal, shutting down...")
		cancel()
	}()

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
