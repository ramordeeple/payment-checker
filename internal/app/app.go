package app

import (
	"payment-checker/internal/adapter/httpapi"
	"payment-checker/internal/adapter/provider"
	"payment-checker/internal/usecase"
)

type App struct {
	Handler *httpapi.Handler
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
