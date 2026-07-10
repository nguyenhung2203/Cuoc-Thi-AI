package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"backend/internal/realtime"
)

func main() {
	// Load .env so GEMINI_API_KEY / LIVEKIT_* are available to the gateway.
	_ = godotenv.Load()
	_ = godotenv.Load("../../.env")

	port := os.Getenv("REALTIME_PORT")
	if port == "" {
		port = "8081"
	}
	addr := ":" + port
	fmt.Printf("Starting AI Interview Platform Realtime Gateway on %s...\n", addr)

	srv := realtime.NewServer(addr)
	
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
