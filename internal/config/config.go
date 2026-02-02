package config

import (
	"flag"
	"fmt"
	"os"
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
}

// Load loads configuration from flags, environment variables, and defaults.
// Precedence: flags > environment variables > defaults.
//
// Flags:
//   - -listen-addr: HTTP server listen address (default: ":9000")
//   - -collect-interval: Interval between collector runs (default: "15s")
//   - -host: Hostname label value for Prometheus metrics (default: os.Hostname() or "unknown")
//   - -env: Environment label value for Prometheus metrics (default: "")
//
// Environment variables:
//   - LISTEN_ADDR: HTTP server listen address
//   - COLLECT_INTERVAL: Interval between collector runs (e.g., "15s", "30s")
//   - HOST: Hostname label value for Prometheus metrics
//   - ENV: Environment label value for Prometheus metrics
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

	// Parse flags
	flag.Parse()

	// Get hostname for default
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Start with defaults
	cfg := Config{
		ListenAddr:      ":9000",
		CollectInterval: 15 * time.Second,
		HostName:        hostname,
		Env:             "",
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

	// Validate configuration
	if cfg.ListenAddr == "" {
		return Config{}, fmt.Errorf("listen address cannot be empty")
	}
	if cfg.CollectInterval <= 0 {
		return Config{}, fmt.Errorf("collect interval must be positive")
	}

	return cfg, nil
}
