package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"payment-checker/internal/adapter/grpcapi"
	"payment-checker/internal/adapter/httpapi"
	"payment-checker/internal/adapter/repository"
	"payment-checker/internal/port"
	"payment-checker/internal/usecase"
	"time"

	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
)

type App struct {
	Policy   *usecase.Policy
	Handler  *httpapi.Handler
	Provider *port.RateByCurrency
}

func (a *App) StartHTTP(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("/validate", a.Handler.ValidatePayment)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		fmt.Printf("Listening http on port %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("HTTP server closed with error: %s\n", err)
		}
	}()

	return httpServer
}

func (a *App) StartGRPC(addr string) *grpc.Server {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcapi.RegisterPaymentCheckerServer(grpcServer, grpcapi.NewGRPCHandler(a.Provider, a.Policy))

	go func() {
		fmt.Printf("Listening grpc on port %s\n", addr)
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalf("grpc server closed with error: %s\n", err)
		}
	}()

	return grpcServer
}

func (a *App) Shutdown(httpSrv *http.Server, grpcSrv *grpc.Server, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if httpSrv != nil {
		_ = httpSrv.Shutdown(ctx)
	}

	if grpcSrv != nil {
		grpcSrv.GracefulStop()
	}

	fmt.Println("Servers shutdown complete")
}

func InitDB(driverName, dsn string) *App {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		log.Fatalf("cannot connect to database: %w", err)
	}

	repo := repository.NewRateRepo(db)
	var fx port.RateByCurrency = repo

	converter := usecase.NewConverter(repo)
	policy := usecase.NewPolicy(converter, usecase.MaxRubKopecks)
	handler := httpapi.NewHandler(fx, policy)

	return &App{
		Handler: handler,
	}
}
