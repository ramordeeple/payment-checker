package domain

import (
	"strings"
)

type CurrencyCode string

func (c CurrencyCode) String() string {
	return string(c)
}

func (c CurrencyCode) Valid() bool {
	return len(c) == 3
}

func NormalizeCurrency(currency string) (CurrencyCode, error) {
	c := strings.ToUpper(strings.TrimSpace(currency))
	if !validateCurrency(c) {
		return "", ErrInvalidCurrencyCode
	}
	return CurrencyCode(c), nil
}

func validateCurrency(currency string) bool {
	if len(currency) < 3 {
		return false
	}

	for _, r := range currency {
		if r < 'A' || r > 'Z' {
			return false
		}
	}

	return true
}
