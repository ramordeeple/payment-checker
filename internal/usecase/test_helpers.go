package usecase

import (
	"payment-checker/internal/domain"
	"testing"
	"time"
)

const testProvider = "test"

type fakeFX struct {
	rate domain.Rate
	err  error
}

func (f fakeFX) GetRate() (domain.Rate, error) {
	return f.rate, nil
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func setupParallel(t *testing.T) time.Time {
	t.Helper()
	t.Parallel()

	return time.Date(
		2025,
		3,
		2,
		0,
		0,
		0,
		0,
		time.UTC)
}

func newPolicyWithRate(t *testing.T, nominal int32, valueScaled int64, currency domain.CurrencyCode) (*Policy, time.Time) {
	t.Helper()
	date := setupParallel(t)

	r, err := domain.NewRate(date, currency, nominal, valueScaled)
	must(t, err)

	conv := NewConverter(fakeFX{rate: r})
	p := NewPolicy(conv, MaxRubKopecks)

	return p, date
}

func newPolicyWithError(t *testing.T, fxErr error) (*Policy, time.Time) {
	t.Helper()
	date := setupParallel(t)

	conv := NewConverter(fakeFX{err: fxErr})
	p := NewPolicy(conv, MaxRubKopecks)

	return p, date
}

func assertValidation(t *testing.T, got ValidateResponse, wantAllowed bool, wantReason domain.ValidationReason) {
	t.Helper()

	if got.Allowed != wantAllowed {
		t.Fatalf("got allowed=%t, expected allowed=%t", got.Allowed, wantAllowed)
	}

	if got.Reason != wantReason {
		t.Fatalf("got reason=%s, expected reason=%s", got.Reason, wantReason)
	}
}
