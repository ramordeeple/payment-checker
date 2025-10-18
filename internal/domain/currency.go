package domain

import (
	"fmt"
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
	if err := validateCurrency(c); err != nil {
		return "", err
	}
	return CurrencyCode(c), nil
}

func validateCurrency(currency string) error {
	if len(currency) < 3 {
		return fmt.Errorf("invalid currency code length: %q", currency)
	}

	for _, r := range currency {
		if r < 'A' || r > 'Z' {
			return fmt.Errorf("currency should contain latin characters: %q", currency)
		}
	}

	return nil
}
