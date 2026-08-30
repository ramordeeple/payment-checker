package cbr

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"time"
)

//go:embed fixtures/*.xml
var fixtures embed.FS

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dateReq := r.URL.Query().Get("date_req")
	if dateReq == "" {
		http.Error(w, "date_req is required", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("02/01/2006", dateReq)
	if err != nil {
		http.Error(
			w,
			"invalid date_req: expected DD/MM/YYYY",
			http.StatusBadRequest,
		)
		return
	}

	fixturePath := fmt.Sprintf(
		"fixtures/%s.xml",
		date.Format("2006-01-02"),
	)

	body, err := fixtures.ReadFile(fixturePath)
	if errors.Is(err, fs.ErrNotExist) {
		http.Error(w, "rates not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to read fixture", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
