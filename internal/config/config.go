package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	// ListenAddr is the address the HTTP server will listen on (e.g., ":9000").
	ListenAddr string
}

// Load loads configuration from flags, environment variables, and defaults.
// Precedence: flags > environment variables > defaults.
//
// Flags:
//   - -listen-addr: HTTP server listen address (default: ":9000")
//
// Environment variables:
//   - LISTEN_ADDR: HTTP server listen address
//
// Returns an error if flag parsing fails or if required configuration is missing.
func Load() (Config, error) {
	// Load .env file if present (ignore "file not found" errors)
	_ = godotenv.Load()

	// Define flags
	listenAddrFlag := flag.String("listen-addr", "", "HTTP server listen address (e.g., :9000)")

	// Parse flags
	flag.Parse()

	// Start with defaults
	cfg := Config{
		ListenAddr: ":9000",
	}

	// Override with environment variable if set
	if envAddr := os.Getenv("LISTEN_ADDR"); envAddr != "" {
		cfg.ListenAddr = envAddr
	}

	// Override with flag if set (highest precedence)
	if *listenAddrFlag != "" {
		cfg.ListenAddr = *listenAddrFlag
	}

	// Validate configuration
	if cfg.ListenAddr == "" {
		return Config{}, fmt.Errorf("listen address cannot be empty")
	}

	return cfg, nil
}
