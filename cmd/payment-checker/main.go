package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	_ "payment-checker/docs"
	"payment-checker/internal/app"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	a := app.InitDB("postgres", dsn, "file:///app/migrations")

	httpSrv := a.StartHTTP(":8080")
	grpcSrv, _, _ := a.StartGRPC(":8090")

	<-ctx.Done()

	a.Shutdown(httpSrv, grpcSrv, 5*time.Second)
}
