package usecase

import (
	"payment-checker/internal/domain"
	"payment-checker/internal/port"
	"time"
)

type Converter struct {
	fx port.RateByCurrency
}

func NewConverter(fx port.RateByCurrency) *Converter {
	return &Converter{fx: fx}
}

func (c *Converter) ToRUB(m domain.Money, date time.Time) (domain.Money, error) {
	rub := domain.CurrencyCode(("RUB"))

	if m.Currency == rub {
		return m, nil
	}

	r, err := c.fx.GetRate(date, m.Currency)
	if err != nil {
		return domain.Money{}, domain.ErrRateUnavailable
	}

	rubAmount := (m.Amount * r.ValueScaled) / (int64(r.Nominal) * domain.RateScale)

	return domain.Money{Amount: rubAmount, Currency: "RUB"}, nil
}
