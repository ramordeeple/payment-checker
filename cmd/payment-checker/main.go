package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"payment-checker/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg := app.Config{
		HTTPAddr:        ":" + envOrDefault("REST_API_PORT", "8080"),
		ShutdownTimeout: 5 * time.Second,
	}

	if err := app.Run(ctx, cfg); err != nil {
		log.Fatalf("application failed: %v", err)
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return fallback
}
