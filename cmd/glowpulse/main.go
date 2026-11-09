package main

import (
	"fmt"
	"os"

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
		fmt.Printf("❌ %s is DOWN: %v\n", result.URL, result.Err)
		os.Exit(1)
	}

	fmt.Printf("✅ %s → %d in %v\n", result.URL, result.StatusCode, result.Duration)
}
