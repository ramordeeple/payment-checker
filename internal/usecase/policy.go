package usecase

import (
	"payment-checker/internal/domain"
	"time"
)

const MaxRubKopecks int64 = 15_000_00 // 15000 руб

type Policy struct {
	conv   *Converter
	maxRUB int64 // Это будет в копейках
}

type ValidateRequest struct {
	Provider string
	Amount   domain.Money
	Date     time.Time
}

type ValidateResponse struct {
	Allowed  bool
	TotalRUB domain.Money
	Reason   domain.ValidationReason
}

func NewPolicy(conv *Converter, maxRub int64) *Policy {
	if maxRub <= 0 {
		maxRub = MaxRubKopecks
	}
	return &Policy{conv: conv, maxRUB: maxRub}
}

func (p *Policy) Validate(req ValidateRequest) (ValidateResponse, error) {
	total, err := p.conv.ToRUB(req.Amount, req.Date)
	if err != nil {
		return ValidateResponse{
			Allowed: false,
			Reason:  domain.ReasonRateUnavailable}, nil
	}

	if total.Amount > p.maxRUB {
		return ValidateResponse{
			Allowed:  false,
			TotalRUB: total,
			Reason:   domain.ReasonLimitExceeded,
		}, nil
	}

	return ValidateResponse{
		Allowed:  true,
		TotalRUB: total,
		Reason:   domain.ReasonOK,
	}, nil
}
