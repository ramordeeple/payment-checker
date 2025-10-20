package app

import (
	"fmt"
	"net"
	"payment-checker/internal/adapter/grpcapi"
	"payment-checker/internal/adapter/httpapi"
	"payment-checker/internal/adapter/provider"
	"payment-checker/internal/usecase"

	"google.golang.org/grpc"
)

type App struct {
	Provider *provider.Provider
	Policy   *usecase.Policy
	Handler  *httpapi.Handler
}

func New() *App {
	p := provider.NewProvider()
	converter := usecase.NewConverter(p)

	policy := usecase.NewPolicy(converter, usecase.MaxRubKopecks)
	handler := httpapi.NewHandler(p, policy)

	return &App{
		Handler: handler,
	}
}

func (a *App) StartGRPC() error {
	addr := ":8090"
	network := "tcp"

	listen, err := net.Listen(network, addr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	grpcapi.RegisterPaymentCheckerServer(s, grpcapi.NewGRPCHandler(a.Provider, a.Policy))

	fmt.Printf("Listening on port %w\n", addr)
	return s.Serve(listen)
}
