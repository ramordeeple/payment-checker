package httpapi

// From frontend
type ValidatePaymentRequest struct {
	Provider string `json:"provider"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Date     string `json:"date"`
}

// To frontend
type ValidatePaymentResponse struct {
	Allowed  bool   `json:"allowed"`
	TotalRUB int64  `json:"totalRUBInKopecks"`
	Reason   string `json:"reason"`
}
