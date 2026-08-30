package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"payment-checker/internal/cbr"
)

func main() {
	port := os.Getenv("REST_API_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.Handle(
		"GET /scripts/XML_daily.asp",
		cbr.NewHandler(),
	)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       time.Minute,
	}

	log.Printf("CBR mock is listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
