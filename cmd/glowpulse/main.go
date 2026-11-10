package main

import (
	"fmt"
	"os"
	"time"

	"github.com/YoumnaSalloum/glowpulse/internal/checker"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: glowpulse <url>")
		os.Exit(1)
	}

	url := os.Args[1]
	result := checker.Check(url)

	if result.Err != nil {
		fmt.Printf("%s %s %s: %v\n", icon(result.Status), result.Status, result.URL, result.Err)
		os.Exit(1)
	}

	fmt.Printf("%s %s %s → %d in %v\n",
		icon(result.Status), result.Status, result.URL, result.StatusCode, result.Duration)

	if !result.CertExpires.IsZero() {
		daysLeft := int(time.Until(result.CertExpires).Hours() / 24)
		fmt.Printf("🔒 TLS certificate expires in %d days (%s)\n",
			daysLeft, result.CertExpires.Format("2006-01-02"))
	}

	if result.Status != checker.StatusUp {
		os.Exit(1)
	}

		if len(result.MissingHeaders) == 0 {
		fmt.Println("🛡️  All security headers present")
	} else {
		fmt.Printf("🛡️  Missing %d security headers:\n", len(result.MissingHeaders))
		for _, h := range result.MissingHeaders {
			fmt.Println("   -", h)
		}
	}
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
