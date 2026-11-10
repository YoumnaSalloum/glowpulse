package checker

import (
	"net/http"
	"time"
)

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
// response time and TLS certificate expiry date.
func Check(url string) Result {
	client := &http.Client{Timeout: 10 * time.Second}

	start := time.Now()
	resp, err := client.Get(url)
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
