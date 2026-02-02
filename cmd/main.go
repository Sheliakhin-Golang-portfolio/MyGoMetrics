package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sheliakhin-Golang-portfolio/MyGoMetrics/internal/collector"
	"github.com/Sheliakhin-Golang-portfolio/MyGoMetrics/internal/config"
	"github.com/Sheliakhin-Golang-portfolio/MyGoMetrics/internal/exporter"
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

	// Instantiate collectors
	collectors := []collector.Collector{
		collector.NewCPUCollector(),
		collector.NewMemoryCollector(),
		collector.NewDiskCollector(),
		collector.NewRuntimeCollector(),
	}

	// Create exporter registry
	registry, err := exporter.NewRegistry(collectors, cfg.HostName, cfg.Env)
	if err != nil {
		log.Fatalf("Failed to create exporter registry: %v", err)
	}

	// Start periodic collection goroutine
	go func() {
		ticker := time.NewTicker(cfg.CollectInterval)
		defer ticker.Stop()

		// Initial collection
		registry.Update(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				registry.Update(ctx)
			}
		}
	}()

	// Start HTTP server with metrics handler
	log.Printf("Starting HTTP server on %s", cfg.ListenAddr)
	log.Printf("Metrics collection interval: %v", cfg.CollectInterval)
	if err := server.Start(ctx, cfg, registry.Handler()); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server shutdown complete")
}
