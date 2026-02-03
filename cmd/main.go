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
	"github.com/Sheliakhin-Golang-portfolio/MyGoMetrics/internal/logger"
	"github.com/Sheliakhin-Golang-portfolio/MyGoMetrics/internal/server"
	"go.uber.org/zap"
)

func main() {
	// Load configuration (parses flags and env vars)
	cfg, err := config.Load()
	if err != nil {
		// Using standard library log instead of zap before zap is initialized
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	if err := logger.Init(cfg.LogLevel); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Create root context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for SIGINT and SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Handle signals in a goroutine
	go func() {
		sig := <-sigCh
		logger.Logger.Info("Received signal, initiating graceful shutdown", zap.Any("signal", sig))
		cancel()
	}()

	// Instantiate all collectors
	allCollectors := []collector.Collector{
		collector.NewCPUCollector(),
		collector.NewMemoryCollector(),
		collector.NewDiskCollector(),
		collector.NewRuntimeCollector(),
	}

	// Filter collectors by EnabledCollectors if specified
	collectors := filterCollectors(allCollectors, cfg.EnabledCollectors)

	// Log enabled collectors
	if len(cfg.EnabledCollectors) > 0 {
		collectorNames := make([]string, len(collectors))
		for i, c := range collectors {
			collectorNames[i] = c.Name()
		}
		logger.Logger.Info("Enabled collectors", zap.Strings("collectors", collectorNames))
	} else {
		logger.Logger.Info("All collectors enabled")
	}

	// Create exporter registry
	registry, err := exporter.NewRegistry(collectors, cfg.HostName, cfg.Env, logger.Logger)
	if err != nil {
		logger.Logger.Fatal("Failed to create exporter registry", zap.Error(err))
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
	logger.Logger.Info("Starting HTTP server", zap.String("listen_addr", cfg.ListenAddr))
	logger.Logger.Info("Metrics collection interval", zap.Duration("collect_interval", cfg.CollectInterval))
	if err := server.Start(ctx, cfg, registry.Handler()); err != nil {
		logger.Logger.Fatal("Server error", zap.Error(err))
	}

	logger.Logger.Info("Server shutdown complete")
}

// filterCollectors filters collectors based on the enabled list.
// If enabled is nil or empty, returns all collectors.
// Unknown collector names are ignored (no-op).
func filterCollectors(all []collector.Collector, enabled []string) []collector.Collector {
	if len(enabled) == 0 {
		return all
	}

	// Build a map for O(1) lookup
	enabledMap := make(map[string]struct{}, len(enabled))
	for _, name := range enabled {
		enabledMap[name] = struct{}{}
	}

	// Filter collectors
	filtered := make([]collector.Collector, 0, len(all))
	for _, c := range all {
		if _, ok := enabledMap[c.Name()]; ok {
			filtered = append(filtered, c)
		}
	}

	return filtered
}
