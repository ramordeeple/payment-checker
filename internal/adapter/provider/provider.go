package provider

import (
	"payment-checker/internal/domain"
	"time"
)

type Provider struct {
	rates map[domain.CurrencyCode]domain.Rate
}

// Через map для удобного добавления новой валюты в дальнейшем
func NewProvider() *Provider {
	date := time.Now()
	return &Provider{
		rates: map[domain.CurrencyCode]domain.Rate{
			domain.CurrencyRUB: {Date: date, Currency: domain.CurrencyRUB, Nominal: 1, ValueScaled: domain.RateScale},
			domain.CurrencyUSD: {Date: date, Currency: domain.CurrencyRUB, Nominal: 1, ValueScaled: 75_1234},
			domain.CurrencyEUR: {Date: date, Currency: domain.CurrencyRUB, Nominal: 1, ValueScaled: 80_5678},
			domain.CurrencyJPY: {Date: date, Currency: domain.CurrencyRUB, Nominal: 1, ValueScaled: 50_1234},
		},
	}
}

func (p *Provider) GetRate(date time.Time, currency domain.CurrencyCode) (domain.Rate, error) {
	if rate, ok := p.rates[currency]; ok {
		return rate, nil
	}

	return domain.Rate{}, domain.ErrRateNotFound
}
