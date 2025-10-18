package domain

import (
	"errors"
	"fmt"
)

type Money struct {
	Amount   int64        // В копейках, центах
	Currency CurrencyCode // RUB, USD
}

func NewMoney(amount int64, currency string) (Money, error) {
	cc, err := NormalizeCurrency(currency)
	if err != nil {
		return Money{}, err
	}
	return Money{Amount: amount, Currency: cc}, nil
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

func (m Money) GreaterThan(other Money) bool {
	if m.Currency != other.Currency {
		return false
	}
	return m.Amount > other.Amount
}

func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot add different currencies")
	}

	return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}, nil
}

func (m Money) Equal(other Money) bool {
	return m.Currency == other.Currency &&
		m.Amount == other.Amount
}
