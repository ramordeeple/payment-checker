package domain

import (
	"time"
)

type Rate struct {
	Date        time.Time
	Currency    CurrencyCode
	Nominal     int32
	ValueScaled int64
	CBRID       string
	NumCode     string
	NameRU      string
}

const RateScale = int64(10_000) // Чтобы можно было 75.1234 хранить как целое число(751234) во избежания ошибок округления

func NewRate(date time.Time, currency CurrencyCode, nominal int32, valueScaled int64) (Rate, error) {
	if nominal <= 0 {
		return Rate{}, ErrInvalidNominal
	}

	if valueScaled <= 0 {
		return Rate{}, ErrInvalidRateValue
	}

	return Rate{
		Date:        date,
		Currency:    currency,
		Nominal:     nominal,
		ValueScaled: valueScaled,
	}, nil
}
