package main

import (
	"context"
	"os"
	"os/signal"
	_ "payment-checker/docs"
	"payment-checker/internal/app"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer stop()

	a := app.InitDB(
		"postgres",
		"postgres://postgres:123@localhost:5433/payment_db?sslmode=disable")

	httpSrv := a.StartHTTP(":8080")
	grpcSrv := a.StartGRPC(":8090")
	<-ctx.Done()

	a.Shutdown(httpSrv, grpcSrv, 5*time.Second)
}
