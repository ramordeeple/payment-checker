package cbr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type cbrResponse struct {
	Date string `xml:"Date,attr"`
	Name string `xml:"name,attr"`

	Valutes []struct {
		ID       string `xml:"ID,attr"`
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Nominal  int    `xml:"Nominal"`
		Name     string `xml:"Name"`
		Value    string `xml:"Value"`
	} `xml:"Valute"`
}

func TestHandler_KnownDateReturnsFixture(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(
		http.MethodGet,
		"/scripts/XML_daily.asp?date_req=02/03/2002",
		nil,
	)
	rec := httptest.NewRecorder()

	NewHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	wantContentType := "application/xml; charset=utf-8"
	if got := rec.Header().Get("Content-Type"); got != wantContentType {
		t.Errorf("Content-Type: got %q, want %q", got, wantContentType)
	}

	wantBody, err := fixtures.ReadFile("fixtures/2002-03-02.xml")
	if err != nil {
		t.Fatalf("read expected fixture: %v", err)
	}

	if !bytes.Equal(rec.Body.Bytes(), wantBody) {
		t.Errorf(
			"response body differs from fixture\ngot:\n%s\nwant:\n%s",
			rec.Body.Bytes(),
			wantBody,
		)
	}
}

func TestHandler_ReturnsValidCBRXML(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(
		http.MethodGet,
		"/scripts/XML_daily.asp?date_req=02/03/2002",
		nil,
	)
	rec := httptest.NewRecorder()

	NewHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"status: got %d, want %d; body=%q",
			rec.Code,
			http.StatusOK,
			rec.Body.String(),
		)
	}

	var response cbrResponse

	if err := xml.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("response is not valid XML: %v", err)
	}

	if len(response.Valutes) != 1 {
		t.Fatalf(
			"Valutes count: got %d, want 1",
			len(response.Valutes),
		)
	}

	valute := response.Valutes[0]

	if valute.CharCode != "USD" {
		t.Errorf("CharCode: got %q, want %q", valute.CharCode, "USD")
	}

	if valute.Nominal != 1 {
		t.Errorf("Nominal: got %d, want %d", valute.Nominal, 1)
	}

	if valute.Value != "75,0000" {
		t.Errorf("Value: got %q, want %q", valute.Value, "75,0000")
	}
}

func TestHandler_ErrorResponses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		target     string
		wantStatus int
		wantAllow  string
	}{
		{
			name:       "missing date",
			method:     http.MethodGet,
			target:     "/scripts/XML_daily.asp",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid date format",
			method:     http.MethodGet,
			target:     "/scripts/XML_daily.asp?date_req=2002-03-02",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "fixture not found",
			method:     http.MethodGet,
			target:     "/scripts/XML_daily.asp?date_req=01/01/2030",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "method not allowed",
			method:     http.MethodPost,
			target:     "/scripts/XML_daily.asp?date_req=02/03/2002",
			wantStatus: http.StatusMethodNotAllowed,
			wantAllow:  http.MethodGet,
		},
	}

	handler := NewHandler()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.method, tt.target, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf(
					"status: got %d, want %d; body=%q",
					rec.Code,
					tt.wantStatus,
					rec.Body.String(),
				)
			}

			if got := rec.Header().Get("Allow"); got != tt.wantAllow {
				t.Errorf(
					"Allow header: got %q, want %q",
					got,
					tt.wantAllow,
				)
			}
		})
	}
}

func TestHandler_ConcurrentRequests(t *testing.T) {
	t.Parallel()

	const requestCount = 100

	wantBody, err := fixtures.ReadFile("fixtures/2002-03-02.xml")
	if err != nil {
		t.Fatalf("read expected fixture: %v", err)
	}

	handler := NewHandler()
	errs := make(chan error, requestCount)

	var wg sync.WaitGroup
	wg.Add(requestCount)

	for requestNumber := range requestCount {
		go func() {
			defer wg.Done()

			req := httptest.NewRequest(
				http.MethodGet,
				"/scripts/XML_daily.asp?date_req=02/03/2002",
				nil,
			)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				errs <- fmt.Errorf(
					"request %d: status got %d, want %d",
					requestNumber,
					rec.Code,
					http.StatusOK,
				)
				return
			}

			if !bytes.Equal(rec.Body.Bytes(), wantBody) {
				errs <- fmt.Errorf(
					"request %d: response differs from fixture",
					requestNumber,
				)
			}
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}
