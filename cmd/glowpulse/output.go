package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/YoumnaSalloum/glowpulse/internal/checker"
)

// printText prints a human-friendly report.
func printText(r checker.Result) {
	if r.Err != nil {
		fmt.Printf("%s %s %s: %v\n", icon(r.Status), r.Status, r.URL, r.Err)
		return
	}

	fmt.Printf("%s %s %s → %d in %v\n",
		icon(r.Status), r.Status, r.URL, r.StatusCode, r.Duration)

	if !r.CertExpires.IsZero() {
		daysLeft := int(time.Until(r.CertExpires).Hours() / 24)
		fmt.Printf("🔒 TLS certificate expires in %d days (%s)\n",
			daysLeft, r.CertExpires.Format("2006-01-02"))
	}

	if len(r.MissingHeaders) == 0 {
		fmt.Println("🛡️  All security headers present")
	} else {
		fmt.Printf("🛡️  Missing %d security headers:\n", len(r.MissingHeaders))
		for _, h := range r.MissingHeaders {
			fmt.Println("   -", h)
		}
	}
}

// jsonResult is the shape of the --json output.
type jsonResult struct {
	URL            string   `json:"url"`
	Status         string   `json:"status"`
	StatusCode     int      `json:"status_code,omitempty"`
	DurationMs     int64    `json:"duration_ms"`
	CertExpires    string   `json:"cert_expires,omitempty"`
	MissingHeaders []string `json:"missing_headers,omitempty"`
	Error          string   `json:"error,omitempty"`
}

// printJSON prints the result as indented JSON.
// toJSON converts one result into its JSON shape.
func toJSON(r checker.Result) jsonResult {
	out := jsonResult{
		URL:            r.URL,
		Status:         string(r.Status),
		StatusCode:     r.StatusCode,
		DurationMs:     r.Duration.Milliseconds(),
		MissingHeaders: r.MissingHeaders,
	}
	if !r.CertExpires.IsZero() {
		out.CertExpires = r.CertExpires.Format(time.RFC3339)
	}
	if r.Err != nil {
		out.Error = r.Err.Error()
	}
	return out
}

// printJSON prints all results as an indented JSON array.
func printJSON(results []checker.Result) error {
	out := make([]jsonResult, 0, len(results))
	for _, r := range results {
		out = append(out, toJSON(r))
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// icon returns an emoji for each health status.
func icon(s checker.Status) string {
	switch s {
	case checker.StatusUp:
		return "✅"
	case checker.StatusDegraded:
		return "⚠️"
	default:
		return "❌"
	}
}

// printSummary prints one line with totals and timing.
func printSummary(results []checker.Result, elapsed time.Duration, workers int) {
	up := 0
	for _, r := range results {
		if r.Status == checker.StatusUp {
			up++
		}
	}
	fmt.Printf("📊 Checked %d sites in %v using %d workers: %d up, %d not up\n",
		len(results), elapsed.Round(time.Millisecond), workers, up, len(results)-up)
}
