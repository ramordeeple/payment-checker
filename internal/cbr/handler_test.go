package cbr

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
)

var response struct {
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
