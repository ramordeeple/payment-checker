package http

// От фронтенда
type ValidatePaymentRequest struct {
	Provider string `json:"provider"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Date     string `json:"date"`
}

// На фронтенд
type ValidatePaymentResponse struct {
	Allowed  bool   `json:"allowed"`
	TotalRUB int64  `json:"totalRUB"`
	Reason   string `json:"reason"`
}
