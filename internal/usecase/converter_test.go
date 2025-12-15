package usecase

import (
	"payment-checker/internal/domain"
	"testing"
)

func TestConverter_ToRUB(t *testing.T) {
	date := setupParallel(t)

	fxCurrency := domain.Currency{
		Code:    "USD",
		NameRU:  "Доллар США",
		NumCode: "840",
		CBRID:   "R01235",
	}
	currencyRUB := domain.Currency{
		Code:    "RUB",
		NameRU:  "Российский рубль",
		NumCode: "643",
		CBRID:   "R01239",
	}

	rateValue := int64(10_0000)
	amountUSD := int64(123_00)

	rate, err := domain.NewRate(date, fxCurrency.Code, 1, rateValue)
	if err != nil {
		t.Fatalf("failed to create rate: %v", err)
	}

	conv := NewConverter(fakeFX{rate: rate})

	money, err := domain.NewMoney(amountUSD, fxCurrency.Code)
	if err != nil {
		t.Fatalf("failed to create money: %v", err)
	}

	got, err := conv.ToRUB(money, date)
	if err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	expected := amountUSD * rateValue / domain.RateScale

	if got.Currency != currencyRUB.Code {
		t.Fatalf("converted currency mismatch: got %s, want %s", got.Currency, currencyRUB.Code)
	}

	if got.Amount != expected {
		t.Fatalf("converted amount mismatch: got %d, want %d", got.Amount, expected)
	}
}
