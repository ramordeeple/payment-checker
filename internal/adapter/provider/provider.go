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
	return &Provider{
		rates: map[domain.CurrencyCode]domain.Rate{
			domain.CurrencyRUB: {Currency: domain.CurrencyRUB, Nominal: 1, ValueScaled: domain.RateScale},
			domain.CurrencyUSD: {Currency: domain.CurrencyRUB, Nominal: 1, ValueScaled: 75_1234},
			domain.CurrencyEUR: {Currency: domain.CurrencyRUB, Nominal: 1, ValueScaled: 80_5678},
			domain.CurrencyJPY: {Currency: domain.CurrencyRUB, Nominal: 1, ValueScaled: 50_1234},
		},
	}
}

func (p *Provider) GetRate(date time.Time, currency domain.CurrencyCode) (domain.Rate, error) {
	rate, ok := p.rates[currency]
	if !ok {
		return domain.Rate{}, domain.ErrRateNotFound
	}
	rate.Date = date

	return rate, nil
}
