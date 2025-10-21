package grpcapi

import (
	"payment-checker/internal/domain"
	"payment-checker/internal/port"
	"payment-checker/internal/usecase"
	"time"
)

type GRPCHandler struct {
	UnimplementedPaymentCheckerServer
	provider port.FXRateProvider
	policy   *usecase.Policy
}

func NewGRPCHandler(provider port.FXRateProvider, policy *usecase.Policy) *GRPCHandler {
	return &GRPCHandler{
		provider: provider,
		policy:   policy,
	}
}

func (h *GRPCHandler) ValidatePayment(
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
