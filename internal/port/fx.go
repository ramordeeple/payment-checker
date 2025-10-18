package port

import (
	"payment-checker/internal/domain"
	"time"
)

type FXRateProvider interface {
	GetRate(date time.Time, currency domain.CurrencyCode) (domain.Rate, error)
}
