package domain

type CurrencyCode string

func validateCurrency(currency CurrencyCode) bool {

	// Код валюты должен иметь строго 3 символа
	if len(currency) != 3 {
		return false
	}

	// Код валюты должен быть строго на латинице
	for _, r := range currency {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}
