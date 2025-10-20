package grpcapi

import (
	"context"
	"payment-checker/internal/adapter/provider"
	"payment-checker/internal/domain"
	"payment-checker/internal/usecase"
	"time"
)

type GRPCHandler struct {
	UnimplementedPaymentCheckerServer
	provider *provider.Provider
	policy   *usecase.Policy
}

func NewGRPCHandler(provider *provider.Provider, policy *usecase.Policy) *GRPCHandler {
	return &GRPCHandler{
		provider: provider,
		policy:   policy,
	}
}

func (h *GRPCHandler) ValidatePayment(
	ctx context.Context,
	req *ValidatePaymentRequest,
) (*ValidatePaymentResponse, error) {

	currency := domain.CurrencyCode(req.Currency)

	if !h.provider.HasCurrency(currency) {
		return nil, domain.ErrRateNotFound
	}

	money := domain.Money{
		Amount:   req.Amount,
		Currency: currency,
	}

	validateReq := usecase.ValidateRequest{
		Provider: req.Provider,
		Amount:   money,
		Date:     time.Now(),
	}

	resp, err := h.policy.Validate(validateReq)
	if err != nil {
		return nil, err
	}

	return &ValidatePaymentResponse{
		Allowed:           resp.Allowed,
		TotalRubInKopecks: resp.TotalRUB.Amount,
		Reason:            string(resp.Reason),
	}, nil
}
