package domain

type CurrencyCode string

type Currency struct {
	Code    CurrencyCode
	NameRU  string
	NumCode string
	CBRID   string
}
