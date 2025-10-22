package grpcapi

import (
	"payment-checker/internal/domain"
	"payment-checker/internal/port"
	"payment-checker/internal/usecase"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	UnimplementedPaymentCheckerServer
	provider      port.RateByCurrency
	policy        *usecase.Policy
	currencyCheck port.CurrencyChecker
}

func NewGRPCHandler(fx port.FXRateProvider, policy *usecase.Policy) *GRPCHandler {
	return &GRPCHandler{
		provider:      fx,
		currencyCheck: fx,
		policy:        policy,
	}
}

func (h *GRPCHandler) ValidatePayment(
	req *ValidatePaymentRequest,
) (*ValidatePaymentResponse, error) {

	currency := domain.CurrencyCode(req.Currency)

	if !h.currencyCheck.HasCurrency(currency) {
		return nil, status.Error(codes.NotFound, "currency not found")
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ValidatePaymentResponse{
		Allowed:           resp.Allowed,
		TotalRubInKopecks: resp.TotalRUB.Amount,
		Reason:            string(resp.Reason),
	}, nil
}
