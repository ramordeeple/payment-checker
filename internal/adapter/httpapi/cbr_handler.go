package httpapi

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"payment-checker/internal/port"
	"time"
)

type CBRDailyHandler struct {
	provider port.RatesByDateProvider
}

type ValCurs struct {
	XMLName xml.Name    `xml:"ValCurs"`
	Date    string      `xml:"Date,attr"`
	Name    string      `xml:"name,attr"`
	Valutes []ValuteXML `xml:"Valute"`
}

type ValuteXML struct {
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int32  `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

func (h *CBRDailyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dateReq := r.URL.Query().Get("date_req")
	if dateReq == "" {
		http.Error(w, "date_req is required, format DD/MM/YYYY", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("02/01/2006", dateReq)
	if err != nil {
		http.Error(w, "invalid date_req format, expected DD/MM/YYYY", http.StatusBadRequest)
		return
	}

	// Получаем курсы по дате
	rates, err := h.provider.GetRatesByDate(date)
	if err != nil {
		fmt.Println("GetRatesByDate error:", err)
		http.Error(w, "rate not found", http.StatusNotFound)
		return
	}

	fmt.Println("Rates returned from GetRatesByDate:", rates)

	response := ValCurs{
		Date: date.Format("02.01.2006"),
		Name: "Foreign Currency Market",
	}

	for _, rate := range rates {
		fmt.Printf("Processing rate: Currency=%s Nominal=%d ValueScaled=%d\n",
			rate.Currency, rate.Nominal, rate.ValueScaled)

		meta, err := h.provider.GetCurrencyMeta(rate.Currency)
		if err != nil {
			fmt.Println("Currency meta not found for", rate.Currency, "error:", err)
			continue
		}

		fmt.Println("Currency meta found:", meta)

		response.Valutes = append(response.Valutes, ValuteXML{
			ID:       meta.CBRID,
			NumCode:  meta.NumCode,
			CharCode: string(rate.Currency),
			Nominal:  rate.Nominal,
			Name:     meta.NameRU,
			Value:    formatCBRValue(rate.ValueScaled),
		})
	}

	if len(response.Valutes) == 0 {
		fmt.Println("No Valutes added to response; check rates and currency meta in DB")
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	xml.NewEncoder(w).Encode(response)
}

func NewCBRDailyHandler(provider port.RatesByDateProvider) *CBRDailyHandler {
	return &CBRDailyHandler{provider: provider}
}

func formatCBRValue(valueScaled int64) string {
	intPart := valueScaled / 10000
	fracPart := valueScaled % 10000
	return fmt.Sprintf("%d,%04d", intPart, fracPart)
}
