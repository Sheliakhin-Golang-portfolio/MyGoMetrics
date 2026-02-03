package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	// ListenAddr is the address the HTTP server will listen on (e.g., ":9000").
	ListenAddr string

	// CollectInterval is the interval between collector runs.
	CollectInterval time.Duration

	// HostName is the hostname label value for Prometheus metrics.
	HostName string

	// Env is the environment label value for Prometheus metrics.
	Env string

	// EnabledCollectors is a list of collector names to enable.
	// If empty, all collectors are enabled.
	// Valid names: cpu, memory, disk, runtime
	EnabledCollectors []string

	// LogLevel is the logging level (debug, info, warn, error).
	// Default: "info"
	LogLevel string
}

// Load loads configuration from flags, environment variables, and defaults.
// Precedence: flags > environment variables > defaults.
//
// Flags:
//   - -listen-addr: HTTP server listen address (default: ":9000")
//   - -collect-interval: Interval between collector runs (default: "15s")
//   - -host: Hostname label value for Prometheus metrics (default: os.Hostname() or "unknown")
//   - -env: Environment label value for Prometheus metrics (default: "")
//   - -enabled-collectors: Comma-separated list of collector names to enable (default: empty = all enabled)
//   - -log-level: Logging level: debug, info, warn, error (default: "info")
//
// Environment variables:
//   - LISTEN_ADDR: HTTP server listen address
//   - COLLECT_INTERVAL: Interval between collector runs (e.g., "15s", "30s")
//   - HOST: Hostname label value for Prometheus metrics
//   - ENV: Environment label value for Prometheus metrics
//   - ENABLED_COLLECTORS: Comma-separated list of collector names to enable (empty = all enabled)
//   - LOG_LEVEL: Logging level: debug, info, warn, error (default: "info")
//
// Returns an error if flag parsing fails or if required configuration is missing.
func Load() (Config, error) {
	// Load .env file if present (ignore "file not found" errors)
	_ = godotenv.Load()

	// Define flags
	listenAddrFlag := flag.String("listen-addr", "", "HTTP server listen address (e.g., :9000)")
	collectIntervalFlag := flag.String("collect-interval", "", "Interval between collector runs (e.g., 15s, 30s)")
	hostFlag := flag.String("host", "", "Hostname label value for Prometheus metrics")
	envFlag := flag.String("env", "", "Environment label value for Prometheus metrics")
	enabledCollectorsFlag := flag.String("enabled-collectors", "", "Comma-separated list of collector names to enable (empty = all enabled)")
	logLevelFlag := flag.String("log-level", "", "Logging level: debug, info, warn, error (default: info)")

	// Parse flags
	flag.Parse()

	// Get hostname for default
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Start with defaults
	cfg := Config{
		ListenAddr:        ":9000",
		CollectInterval:   15 * time.Second,
		HostName:          hostname,
		Env:               "",
		EnabledCollectors: nil, // nil/empty = all enabled
		LogLevel:          "info",
	}

	// Override with environment variables if set
	if envAddr := os.Getenv("LISTEN_ADDR"); envAddr != "" {
		cfg.ListenAddr = envAddr
	}
	if envInterval := os.Getenv("COLLECT_INTERVAL"); envInterval != "" {
		parsedInterval, err := time.ParseDuration(envInterval)
		if err != nil {
			return Config{}, fmt.Errorf("invalid COLLECT_INTERVAL: %w", err)
		}
		cfg.CollectInterval = parsedInterval
	}
	if envHost := os.Getenv("HOST"); envHost != "" {
		cfg.HostName = envHost
	}
	if envEnv := os.Getenv("ENV"); envEnv != "" {
		cfg.Env = envEnv
	}
	if envCollectors := os.Getenv("ENABLED_COLLECTORS"); envCollectors != "" {
		cfg.EnabledCollectors = parseCollectors(envCollectors)
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	}

	// Override with flags if set (highest precedence)
	if *listenAddrFlag != "" {
		cfg.ListenAddr = *listenAddrFlag
	}
	if *collectIntervalFlag != "" {
		parsedInterval, err := time.ParseDuration(*collectIntervalFlag)
		if err != nil {
			return Config{}, fmt.Errorf("invalid -collect-interval: %w", err)
		}
		cfg.CollectInterval = parsedInterval
	}
	if *hostFlag != "" {
		cfg.HostName = *hostFlag
	}
	if *envFlag != "" {
		cfg.Env = *envFlag
	}
	if *enabledCollectorsFlag != "" {
		cfg.EnabledCollectors = parseCollectors(*enabledCollectorsFlag)
	}
	if *logLevelFlag != "" {
		cfg.LogLevel = *logLevelFlag
	}

	// Validate configuration
	if cfg.ListenAddr == "" {
		return Config{}, fmt.Errorf("listen address cannot be empty")
	}
	if cfg.CollectInterval <= 0 {
		return Config{}, fmt.Errorf("collect interval must be positive")
	}

	return cfg, nil
}

// parseCollectors parses a comma-separated list of collector names,
// trims whitespace, and returns a slice of names.
// Empty string returns nil (all collectors enabled).
func parseCollectors(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
