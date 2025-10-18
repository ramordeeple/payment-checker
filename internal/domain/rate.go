package domain

import (
	"errors"
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
		return Rate{}, err
	}

	if nominal <= 0 {
		return Rate{}, errors.New("nominal must be greater than 0")
	}

	if valueScaled <= 0 {
		return Rate{}, errors.New("valueScaled must be greater or equal 0")
	}

	return Rate{
		Date:        date,
		Currency:    cc,
		Nominal:     nominal,
		ValueScaled: valueScaled,
	}, nil
}
