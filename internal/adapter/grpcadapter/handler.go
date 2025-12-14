package grpcadapter

import (
	"context"
	paymentcheckerv1 "payment-checker/internal/grpc/proto/paymentchecker/v1"
	"time"

	"payment-checker/internal/domain"
	"payment-checker/internal/port"
	"payment-checker/internal/usecase"
)

type Handler struct {
	paymentcheckerv1.UnimplementedPaymentCheckerServer

	provider port.FXRateProvider
	policy   *usecase.Policy
}

func NewHandler(provider port.FXRateProvider, policy *usecase.Policy) *Handler {
	return &Handler{
		provider: provider,
		policy:   policy,
	}
}

func (h *Handler) Validate(
	ctx context.Context,
	req *paymentcheckerv1.ValidatePaymentRequest,
) (*paymentcheckerv1.ValidatePaymentResponse, error) {

	currency := domain.CurrencyCode(req.Currency)
	if !h.provider.HasCurrency(currency) {
		return &paymentcheckerv1.ValidatePaymentResponse{
			Allowed: false,
			Reason:  paymentcheckerv1.ValidationReason_RATE_UNAVAILABLE,
		}, nil
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return &paymentcheckerv1.ValidatePaymentResponse{
			Allowed: false,
			Reason:  paymentcheckerv1.ValidationReason_RATE_UNAVAILABLE,
		}, nil
	}

	resp, err := h.policy.Validate(usecase.ValidateRequest{
		Provider: req.Provider,
		Amount: domain.Money{
			Amount:   req.Amount,
			Currency: currency,
		},
		Date: date,
	})
	if err != nil {
		return &paymentcheckerv1.ValidatePaymentResponse{
			Allowed: false,
			Reason:  paymentcheckerv1.ValidationReason_RATE_UNAVAILABLE,
		}, nil
	}

	return &paymentcheckerv1.ValidatePaymentResponse{
		Allowed:           resp.Allowed,
		TotalRubInKopecks: resp.TotalRUB.Amount,
		Reason:            mapReason(resp.Reason),
	}, nil
}
