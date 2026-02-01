package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Sheliakhin-Golang-portfolio/MyGoMetrics/internal/config"
)

// Start starts the HTTP server with the given configuration.
// The server listens on cfg.ListenAddr and exposes a /healthcheck endpoint.
// The function blocks until the context is cancelled, at which point it performs
// a graceful shutdown.
//
// Returns an error if the server fails to start or if shutdown fails.
func Start(ctx context.Context, cfg config.Config) error {
	// Create HTTP server
	srv := &http.Server{
		Addr:         cfg.ListenAddr,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", healthcheckHandler)
	srv.Handler = mux

	// Start server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Wait for context cancellation or server error
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		// Context cancelled, perform graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}
		return nil
	}
}

// healthcheckHandler handles GET /healthcheck requests.
// Returns 200 OK with a simple status response.
func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
