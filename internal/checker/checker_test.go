package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStatusFromCode(t *testing.T) {
	tests := []struct {
		code int
		want Status
	}{
		{200, StatusUp},
		{301, StatusUp},
		{404, StatusDegraded},
		{500, StatusDown},
	}

	for _, tt := range tests {
		got := statusFromCode(tt.code)
		if got != tt.want {
			t.Errorf("statusFromCode(%d) = %s, want %s", tt.code, got, tt.want)
		}
	}
}

func TestCheckReturnsDegradedOn404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	result := Check(context.Background(), server.URL, time.Second)

	if result.Status != StatusDegraded {
		t.Errorf("Status = %s, want %s", result.Status, StatusDegraded)
	}
	if result.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d, want %d", result.StatusCode, http.StatusNotFound)
	}
}

func TestCheckTimesOutOnSlowSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
	}))
	defer server.Close()

	result := Check(context.Background(), server.URL, 50*time.Millisecond)

	if result.Status != StatusDown {
		t.Errorf("Status = %s, want %s", result.Status, StatusDown)
	}
	if result.Err == nil {
		t.Error("expected a timeout error, got nil")
	}
}
