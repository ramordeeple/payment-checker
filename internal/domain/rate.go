package domain

import (
	"fmt"
	"time"
)

const RateScale = int64(10_000) // Чтобы можно было 75.1234 хранить как целое число(751234) во избежания ошибок округления

type Rate struct {
	Date        time.Time
	Currency    CurrencyCode
	Nominal     int32
	ValueScaled int64
}

func NewRate(date time.Time, currency string, nominal int32, valueScaled int64) (Rate, error) {
	cc, err := NormalizeCurrency(currency)
	if err != nil {
		return Rate{}, ErrInvalidCurrencyCode
	}

	if nominal <= 0 {
		return Rate{}, ErrInvalidNominal
	}

	if valueScaled <= 0 {
		return Rate{}, ErrInvalidRateValue
	}

	return Rate{
		Date:        date,
		Currency:    cc,
		Nominal:     nominal,
		ValueScaled: valueScaled,
	}, nil
}

// Для удобной совместимости с ЦБ xml
func (r Rate) FormatValueScaled() string {
	return fmt.Sprintf("%d,%04d", r.ValueScaled/RateScale, r.ValueScaled%RateScale)
}
