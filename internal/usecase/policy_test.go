package usecase

import (
	"payment-checker/internal/domain"
	"testing"
)

var (
	usdCurrency = domain.Currency{Code: "USD", NameRU: "Доллар США", NumCode: "840", CBRID: "R01235"}
)

func TestPolicy_Validate_OK_AtLimit(t *testing.T) {
	const rate = 75_0000
	const amountUSD = 200_00

	p, date := newPolicyWithRate(t, 1, rate, usdCurrency.Code)

	usd, err := domain.NewMoney(amountUSD, usdCurrency.Code)
	if err != nil {
		t.Fatalf("failed to create money: %v", err)
	}

	resp, err := p.Validate(ValidateRequest{
		Provider: testProvider,
		Amount:   usd,
		Date:     date,
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}

	// Проверяем reason
	if resp.Reason != domain.ReasonOK {
		t.Fatalf("expected reason %s, got %s", domain.ReasonOK, resp.Reason)
	}

	// Проверяем сумму RUB
	if resp.TotalRUB.Amount != MaxRubKopecks {
		t.Fatalf("TotalRUB.Amount = %d, want %d", resp.TotalRUB.Amount, MaxRubKopecks)
	}
}

func TestPolicy_Validate_LimitExceeded(t *testing.T) {
	p, date := newPolicyWithRate(t, 1, 75_0000, usdCurrency.Code)

	usd, err := domain.NewMoney(200_01, usdCurrency.Code)
	if err != nil {
		t.Fatalf("failed to create money: %v", err)
	}

	resp, err := p.Validate(ValidateRequest{
		Provider: testProvider,
		Amount:   usd,
		Date:     date,
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}

	if resp.Reason != domain.ReasonLimitExceeded {
		t.Fatalf("expected reason %s, got %s", domain.ReasonLimitExceeded, resp.Reason)
	}
}

func TestPolicy_Validate_RateUnavailable(t *testing.T) {
	p, date := newPolicyWithError(t, domain.ErrRateUnavailable)

	usd, err := domain.NewMoney(10_00, usdCurrency.Code)
	if err != nil {
		t.Fatalf("failed to create money: %v", err)
	}

	resp, err := p.Validate(ValidateRequest{
		Provider: testProvider,
		Amount:   usd,
		Date:     date,
	})
	if err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}

	if resp.Reason != domain.ReasonRateUnavailable {
		t.Fatalf("expected reason %s, got %s", domain.ReasonRateUnavailable, resp.Reason)
	}
}
