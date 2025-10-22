package domain

import "errors"

var (
	ErrInvalidNominal   = errors.New("invalid rate nominal")
	ErrInvalidRateValue = errors.New("invalid rate value")
	ErrRateNotFound     = errors.New("rate not found")
	ErrRateUnavailable  = errors.New("foreign exchange rate is unavailable")
	ErrCurrencyNotFound = errors.New("currency not found")
	ErrInvalidDate      = errors.New("invalid date")
)
