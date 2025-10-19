package usecase

import (
	"payment-checker/internal/domain"
	"testing"
)

const (
	rateValue   = 10_0000 // Условно 10.0000 RUB за 1 USD
	amountUSD   = 123_00  // 123.00 USD
	fxCurrency  = domain.CurrencyUSD
	currencyRUB = domain.CurrencyRUB
)

func TestConverter_ToRUB(t *testing.T) {
	date := setupParallel(t)

	rate, err := domain.NewRate(date, fxCurrency, 1, rateValue) // 10.0000 RUB за USD
	must(t, err)

	conv := NewConverter(fakeFX{rate: rate})
	money, err := domain.NewMoney(amountUSD, fxCurrency) // $123.00
	must(t, err)

	got, err := conv.ToRUB(money, date)
	must(t, err)

	expected := int64(amountUSD) * rateValue / domain.RateScale

	if got.Currency != currencyRUB {
		t.Fatalf("converted currency: %s (from %s), expected %s", got.Currency, fxCurrency, currencyRUB)
	}

	if got.Amount != expected {
		t.Fatalf("got amount %d, expected %d", got.Amount, expected)
	}
}
