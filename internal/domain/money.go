package domain

import (
	"fmt"
)

type Money struct {
	Amount   int64        // В копейках, центах
	Currency CurrencyCode // "RUB", "USD"
}

func NewMoney(amount int64, currency CurrencyCode) (Money, error) {
	return Money{Amount: amount, Currency: currency}, nil
}

func (m Money) String() string {
	sign := ""
	amount := m.Amount

	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	intPart := amount / 100      // Целые рубли / доллары
	fractionPart := amount % 100 // Копейки / центы

	return fmt.Sprintf("%s%d.%02d %s",
		sign, intPart, fractionPart, m.Currency)
}
