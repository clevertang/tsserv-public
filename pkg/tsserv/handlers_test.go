package tsserv

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSayHello(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sayHello)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Hello\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetRawDataPoints(t *testing.T) {
	begin := "2023-01-01T00:00:00Z"
	end := "2023-01-01T01:00:00Z"
	url := "/data?begin=" + begin + "&end=" + end

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getRawDataPoints)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check for a valid timestamp in the response
	expectedPrefix := "2023-01-01T00:00:00Z"
	if !startsWith(rr.Body.String(), expectedPrefix) {
		t.Errorf("handler returned unexpected body: got %v want prefix %v",
			rr.Body.String(), expectedPrefix)
	}
}

func startsWith(body, prefix string) bool {
	return len(body) >= len(prefix) && body[:len(prefix)] == prefix
}

func TestGetRawDataPointsInvalidParams(t *testing.T) {
	url := "/data?begin=invalid&end=invalid"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getRawDataPoints)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := "Param 'begin' must be in RFC3339 format"
	if !contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func contains(body, substr string) bool {
	return len(body) >= len(substr) && body[:len(substr)] == substr
}
