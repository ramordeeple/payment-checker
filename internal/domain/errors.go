package domain

import "errors"

var (
	ErrInvalidCurrencyCode = errors.New("invalid currency code")
	ErrInvalidNominal      = errors.New("invalid rate nominal")
	ErrInvalidRateValue    = errors.New("invalid rate value")
	ErrRateNotFound        = errors.New("rate not found")
	ErrDifferentCurrencies = errors.New("cannot provide an operation with different currencies")
	ErrRateUnavailable     = errors.New("foreign exchange rate is unavailable")
)
