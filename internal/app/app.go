package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"payment-checker/internal/cbr"
)

type Config struct {
	HTTPAddr        string
	ShutdownTimeout time.Duration
}

func Run(ctx context.Context, cfg Config) error {
	server := newHTTPServer(cfg.HTTPAddr)
	serverErr := startHTTPServer(server)

	select {
	case err := <-serverErr:
		return checkHTTPServerError(err)

	case <-ctx.Done():
		log.Printf("shutdown signal received")
		return shutdownHTTPServer(
			server,
			serverErr,
			cfg.ShutdownTimeout,
		)
	}
}

func newHTTPServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle(
		"GET /scripts/XML_daily.asp",
		cbr.NewHandler(),
	)

	return &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       time.Minute,
	}
}

func startHTTPServer(server *http.Server) <-chan error {
	serverErr := make(chan error, 1)

	go func() {
		log.Printf("CBR mock is listening on %s", server.Addr)
		serverErr <- server.ListenAndServe()
	}()

	return serverErr
}

func shutdownHTTPServer(
	server *http.Server,
	serverErr <-chan error,
	timeout time.Duration,
) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		_ = server.Close()
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	if err := checkHTTPServerError(<-serverErr); err != nil {
		return err
	}

	log.Printf("CBR mock stopped")

	return nil
}

func checkHTTPServerError(err error) error {
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return fmt.Errorf("HTTP server stopped: %w", err)
}
