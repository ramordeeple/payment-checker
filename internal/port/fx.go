package port

import (
	"payment-checker/internal/domain"
	"time"
)

type RateByCurrency interface {
	GetRate(date time.Time, currency domain.CurrencyCode) (domain.Rate, error)
}

type CurrencyChecker interface {
	HasCurrency(currency domain.CurrencyCode) bool
}

type RatesByDateProvider interface {
	GetRatesByDate(date time.Time) ([]domain.Rate, error)
	GetCurrencyMeta(code domain.CurrencyCode) (domain.Currency, error)
}

type FXRateProvider interface {
	RateByCurrency
	CurrencyChecker
	RatesByDateProvider
}
