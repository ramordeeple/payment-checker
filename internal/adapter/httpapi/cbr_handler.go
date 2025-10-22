package httpapi

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"payment-checker/internal/domain"
	"payment-checker/internal/port"
	"time"
)

type CBRHandler struct {
	rateProvider port.RateByCurrency
}

type CBRRateResponse struct {
	XMLName     xml.Name `xml:"Valute" json:"-"`
	CharCode    string   `xml:"CharCode" json:"char_code"`
	Nominal     int32    `xml:"Nominal" json:"nominal"`
	ValueScaled int64    `xml:"Value" json:"value_scaled"`
	Date        string   `xml:"-" json:"date"`
}

func NewCBRHandler(rateProvider port.RateByCurrency) *CBRHandler {
	return &CBRHandler{rateProvider: rateProvider}
}

// GetCBRRate godoc
// @Summary Get currency rate
// @Description Returns the rate for a given currency and date from the mock CBR service
// @Tags cbr
// @Accept json
// @Produce json
// @Accept xml
// @Produce xml
// @Param currency query string true "Currency code"
// @Param date query string false "Date in YYYY-MM-DD format"
// @Success 200 {object} CBRRateResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Rate not found"
// @Router /cbr [get]
func (h *CBRHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	currency, date, err := parseAndValidateParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rate, err := h.rateProvider.GetRate(date, currency)
	if handleRateError(w, err) {
		return
	}

	response := CBRRateResponse{
		CharCode:    string(rate.Currency),
		Nominal:     rate.Nominal,
		ValueScaled: rate.ValueScaled,
		Date:        date.Format("2006-01-02"),
	}

	// Определяем формат ответа
	accept := r.Header.Get("Accept")
	if accept == "application/xml" {
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		xml.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
	}
}

func parseAndValidateParams(r *http.Request) (domain.CurrencyCode, time.Time, error) {
	currency := domain.CurrencyCode(r.URL.Query().Get("currency"))
	if currency == "" {
		return "", time.Time{}, errors.New("currency query parameter is required")
	}

	dateStr := r.URL.Query().Get("date")
	date, err := parseDate(dateStr)
	if err != nil {
		return "", time.Time{}, err
	}

	return currency, date, nil
}

func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Now(), nil
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}
	return date, nil
}

func handleRateError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, domain.ErrRateNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return true
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
	return true
}
