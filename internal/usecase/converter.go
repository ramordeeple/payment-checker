package usecase

import (
	"payment-checker/internal/domain"
	"payment-checker/internal/port"
	"time"
)

type Converter struct {
	fx port.FXRateProvider
}

func NewConverter(fx port.FXRateProvider) *Converter {
	return &Converter{fx: fx}
}

func (c *Converter) ToRUB(m domain.Money, date time.Time) (domain.Money, error) {
	if m.Currency == domain.CurrencyRUB {
		return m, nil
	}

	r, err := c.fx.GetRate(date, m.Currency)
	if err != nil {
		return domain.Money{}, domain.ErrRateUnavailable
	}

	rub := (m.Amount * r.ValueScaled) / (int64(r.Nominal) * domain.RateScale)

	return domain.Money{Amount: rub, Currency: domain.CurrencyRUB}, nil
}
