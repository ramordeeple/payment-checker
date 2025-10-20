package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"payment-checker/internal/app"
	"time"
)

func main() {
	a := app.New()

	http.HandleFunc("/validate", a.Handler.ValidatePayment)
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	go func() {
		fmt.Printf("Listening http on port %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("Server closed with error: %s\n", err)
		}
	}()

	go func() {
		if err := a.StartGRPC(); err != err {
			log.Fatalf("grpc server closed with error: %s\n", err)
		}
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer stop()

	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown error: %s\n", err)
	}

	fmt.Println("Servers shutdown complete")
}
