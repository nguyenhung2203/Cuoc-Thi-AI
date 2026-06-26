package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend/internal/realtime"
)

func main() {
	fmt.Println("Starting AI Interview Platform Realtime Gateway on :8080...")

	srv := realtime.NewServer(":8080")
	
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
