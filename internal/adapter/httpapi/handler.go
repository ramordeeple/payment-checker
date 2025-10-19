package httpapi

import (
	"encoding/json"
	"net/http"
	"payment-checker/internal/adapter/provider"
	"payment-checker/internal/domain"
	"payment-checker/internal/usecase"
	"time"
)

type Handler struct {
	provider *provider.Provider
	policy   *usecase.Policy
}

func NewHandler(p *provider.Provider, policy *usecase.Policy) *Handler {
	return &Handler{
		provider: p,
		policy:   policy,
	}
}

func (h *Handler) ValidatePayment(w http.ResponseWriter, r *http.Request) {
	var reqDTO ValidatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&reqDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currency := domain.CurrencyCode(reqDTO.Currency)

	if !h.provider.HasCurrency(currency) {
		http.Error(w, domain.ErrCurrencyNotFound.Error(), http.StatusBadRequest)
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
