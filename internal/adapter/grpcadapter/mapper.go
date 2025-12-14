package grpcadapter

import (
	"payment-checker/internal/domain"
	paymentcheckerv1 "payment-checker/internal/grpc/proto/paymentchecker/v1"
)

func mapReason(reason domain.ValidationReason) paymentcheckerv1.ValidationReason {
	switch reason {
	case domain.ReasonOK:
		return paymentcheckerv1.ValidationReason_OK

	case domain.ReasonRateUnavailable:
		return paymentcheckerv1.ValidationReason_RATE_UNAVAILABLE

	case domain.ReasonLimitExceeded:
		return paymentcheckerv1.ValidationReason_LIMIT_EXCEEDED

	default:
		return paymentcheckerv1.ValidationReason_VALIDATION_REASON_UNSPECIFIED
	}
}
