package domain

// Выноска в отдельный файл, а не в currency.go, т.к. при большом кол-ве валют будет целеобразнее хранить их здесь
const (
	CurrencyRUB = CurrencyCode("RUB")
	CurrencyUSD = CurrencyCode("USD")
	CurrencyEUR = CurrencyCode("EUR")
	CurrencyJPY = CurrencyCode("JPY")
)
