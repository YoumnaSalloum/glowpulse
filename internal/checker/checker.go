package checker

import (
	"context"
	"net/http"
	"time"
)

// client is shared by every check so connections can be reused.
var client = &http.Client{}

// Status represents the status of a URL check.
type Status string

const (
	StatusUp       Status = "UP"
	StatusDegraded Status = "DEGRADED"
	StatusDown     Status = "DOWN"
)

// Result represents the result of a URL check.
type Result struct {
	URL            string
	StatusCode     int
	Duration       time.Duration
	Status         Status
	CertExpires    time.Time
	MissingHeaders []string
	Err            error
}

// Check visits the given URL and returns its health, status code,
// response time, TLS certificate expiry and missing security headers.
// It stops early if ctx is cancelled or the timeout runs out.
func Check(ctx context.Context, url string, timeout time.Duration) Result {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{URL: url, Status: StatusDown, Err: err}
	}

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return Result{URL: url, Status: StatusDown, Duration: duration, Err: err}
	}
	defer resp.Body.Close()

	result := Result{
		URL:            url,
		Status:         statusFromCode(resp.StatusCode),
		StatusCode:     resp.StatusCode,
		Duration:       duration,
		MissingHeaders: missingHeaders(resp.Header),
	}

	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		result.CertExpires = resp.TLS.PeerCertificates[0].NotAfter
	}

	return result

}

// statusFromCode turns an HTTP status code into a health status.
func statusFromCode(code int) Status {
	switch {
	case code >= 200 && code < 400:
		return StatusUp
	case code >= 400 && code < 500:
		return StatusDegraded
	default:
		return StatusDown
	}
}

// securityHeaders lists the headers a well-protected website should send.
var securityHeaders = []string{
	"Strict-Transport-Security",
	"Content-Security-Policy",
	"X-Content-Type-Options",
	"X-Frame-Options",
	"Referrer-Policy",
}

// missingHeaders returns the security headers that the response did not include.
func missingHeaders(h http.Header) []string {
	var missing []string
	for _, name := range securityHeaders {
		if h.Get(name) == "" {
			missing = append(missing, name)
		}
	}
	return missing
}
