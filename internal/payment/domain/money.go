package domain

import (
	"errors"
	"fmt"
)

type Money struct {
	Amount   int64  // В копейках, центах
	Currency string // RUB, USD
	Scale    uint8  // Кол-во знаков после запятой
}

func NewMoney(amount int64, currency string, scale uint8) (Money, error) {
	if scale > 8 {
		return Money{}, errors.New("Scale is too large")
	}
	if currency == "" {
		return Money{}, errors.New("Currency is required")
	}

	return Money{amount, currency, scale}, nil
}

func (m Money) String() string {
	sign := ""
	amount := m.Amount

	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	scaleFactor := pow10(int64(m.Scale))

	intPart := amount / scaleFactor
	fractionPart := amount % 10

	return fmt.Sprintf("%s%d.%0*d %s",
		sign, intPart, m.Scale, fractionPart, m.Currency)
}

func pow10(n int64) int64 {
	res := int64(1)
	for i := int64(0); i < n; i++ {
		res *= 10
	}
	return res
}
