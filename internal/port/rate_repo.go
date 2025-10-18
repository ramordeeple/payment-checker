package port

import (
	"payment-checker/internal/domain"
	"time"
)

type RateRepository interface {
	// Обновление или вставка курса
	Upsert(rate domain.Rate) error

	// Получение курса
	GetByDateAndCurrency(date time.Time, currency domain.CurrencyCode) (domain.Rate, error)

	// Все курсы на дату
	ListByDate(date time.Time, currency domain.CurrencyCode) (domain.Rate, error)
}
