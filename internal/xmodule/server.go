package xmodule

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// Serve runs the API and the embedded React app on one listener until ctx is
// cancelled, then shuts down gracefully.
func Serve(ctx context.Context, cfg Config, deps Dependencies, log *slog.Logger) error {
	top := http.NewServeMux()
	top.Handle("/api/", deps.API.Handler())
	top.Handle("/api", deps.API.Handler())
	top.Handle("/", deps.WebUI.Handler())

	server := &http.Server{
		Addr:              cfg.Listen,
		Handler:           top,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	listener, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	errCh := make(chan error, 1)
	go func() { errCh <- server.Serve(listener) }()
	log.Info("server listening", "addr", cfg.Listen, "version", Version)

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		log.Info("server stopped")
		return nil
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
