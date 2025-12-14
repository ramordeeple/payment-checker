package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"payment-checker/internal/adapter/grpcadapter"
	"payment-checker/internal/adapter/httpapi"
	"payment-checker/internal/adapter/repository"
	paymentcheckerv1 "payment-checker/internal/grpc/proto/paymentchecker/v1"

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
	Provider port.FXRateProvider
}

func (a *App) StartHTTP(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("/validate", a.Handler.ValidatePayment)

	cbrHandler := httpapi.NewCBRDailyHandler(a.Provider)
	mux.Handle("/scripts/XML_daily.asp", cbrHandler)

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

func (a *App) StartGRPC(addr string) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	grpcServer := grpc.NewServer()

	paymentcheckerv1.RegisterPaymentCheckerServer(
		grpcServer,
		grpcadapter.NewHandler(a.Provider, a.Policy),
	)

	go func() {
		log.Printf("gRPC server listening on %s", addr)
		if err := grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf("gRPC server stopped with error: %v", err)
		}
	}()

	return grpcServer, lis, nil
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

func InitDB(driverName, dsn, migrationsPath string) *App {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("cannot ping database: %v", err)
	}

	if err := repository.RunMigrations(db, migrationsPath); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	repo := repository.NewRateRepo(db)
	converter := usecase.NewConverter(repo)
	policy := usecase.NewPolicy(converter, usecase.MaxRubKopecks)
	handler := httpapi.NewHandler(repo, policy)

	return &App{
		Handler:  handler,
		Policy:   policy,
		Provider: repo,
	}
}
