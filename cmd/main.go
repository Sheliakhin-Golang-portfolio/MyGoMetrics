package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sheliakhin-Golang-portfolio/MyGoMetrics/internal/config"
	"github.com/Sheliakhin-Golang-portfolio/MyGoMetrics/internal/server"
)

func main() {
	// Load configuration (parses flags and env vars)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create root context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for SIGINT and SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Handle signals in a goroutine
	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v, initiating graceful shutdown...", sig)
		cancel()
	}()

	// Start HTTP server
	log.Printf("Starting HTTP server on %s", cfg.ListenAddr)
	if err := server.Start(ctx, cfg); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server shutdown complete")
}
