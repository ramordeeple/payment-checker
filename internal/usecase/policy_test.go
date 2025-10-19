package usecase

import (
	"payment-checker/internal/domain"
	"testing"
)

func TestPolicy_Validate_OK_AtLimit(t *testing.T) {
	p, date := newPolicyWithRate(t, 1, 75_0000, domain.CurrencyUSD) // Условно 75.0000 за 1 доллар
	usd, err := domain.NewMoney(200_00, domain.CurrencyUSD)         // 200 * 75 == 15_000.00 руб
	must(t, err)

	resp, err := p.Validate(ValidateRequest{
		Provider: testProvider,
		Amount:   usd,
		Date:     date,
	})
	must(t, err)

	assertValidation(t, resp, true, domain.ReasonOK)

	if resp.TotalRUB.Amount != MaxRubKopecks {
		t.Fatalf("Total rub=%s, expected=%d", resp.TotalRUB, MaxRubKopecks)
	}
}

func TestPolicy_Validate_LimitExceeded(t *testing.T) {
	p, date := newPolicyWithRate(t, 1, 75_0000, domain.CurrencyUSD)
	usd, err := domain.NewMoney(200_01, domain.CurrencyUSD)

	resp, err := p.Validate(ValidateRequest{
		Provider: testProvider,
		Amount:   usd,
		Date:     date,
	})
	must(t, err)

	assertValidation(t, resp, false, domain.ReasonLimitExceeded)
}

func TestPolicy_Validate_RateUnavailable(t *testing.T) {
	p, date := newPolicyWithError(t, domain.ErrRateUnavailable)
	usd, err := domain.NewMoney(10_00, domain.CurrencyUSD)
	must(t, err)

	resp, err := p.Validate(ValidateRequest{
		Provider: testProvider,
		Amount:   usd,
		Date:     date,
	})

	assertValidation(t, resp, false, domain.ReasonRateUnavailable)
}
