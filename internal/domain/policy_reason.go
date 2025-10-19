package domain

type ValidationReason string

const (
	ReasonOK              ValidationReason = "OK"
	ReasonRateUnavailable ValidationReason = "RATE_UNAVAILABLE"
	ReasonLimitExceeded   ValidationReason = "LIMIT_EXCEEDED"
)
