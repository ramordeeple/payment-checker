package httpapi

import (
	"encoding/json"
	"net/http"
	"payment-checker/internal/domain"
	"payment-checker/internal/port"
	"payment-checker/internal/usecase"
	"time"
)

type Handler struct {
	provider      port.RateByCurrency
	policy        *usecase.Policy
	currencyCheck port.CurrencyChecker
}

func NewHandler(fx port.RateByCurrency, policy *usecase.Policy) *Handler {
	return &Handler{
		provider: fx,
		policy:   policy,
	}
}

// ValidatePayment godoc
// @Summary Validate foreign payment
// @Description Validates a payment request with amount, currency, and date
// @Tags payments
// @Accept json
// @Produce json
// @Param request body ValidatePaymentRequest true "Payment request"
// @Success 200 {object} ValidatePaymentResponse
// @Failure 400 {string} string "Bad request"
// @Router /validate [post]
func (h *Handler) ValidatePayment(w http.ResponseWriter, r *http.Request) {
	var reqDTO ValidatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currency := domain.CurrencyCode(reqDTO.Currency)

	if !h.currencyCheck.HasCurrency(currency) {
		http.Error(w, domain.ErrCurrencyNotFound.Error(), http.StatusNotFound)
		return
	}

	money := domain.Money{
		Amount:   reqDTO.Amount,
		Currency: currency,
	}

	date, err := time.Parse("2006-01-02", reqDTO.Date)
	if err != nil {
		http.Error(w, domain.ErrInvalidDate.Error(), http.StatusBadRequest)
		return
	}

	validateReq := usecase.ValidateRequest{
		Provider: reqDTO.Provider,
		Amount:   money,
		Date:     date,
	}

	resp, err := h.policy.Validate(validateReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	httpResponse := ValidatePaymentResponse{
		Allowed:  resp.Allowed,
		TotalRUB: resp.TotalRUB.Amount,
		Reason:   string(resp.Reason),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(httpResponse)

}
