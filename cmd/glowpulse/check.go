package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/YoumnaSalloum/glowpulse/internal/checker"
	"github.com/YoumnaSalloum/glowpulse/internal/pool"
)

var (
	jsonOutput bool
	workers    int
)

var checkCmd = &cobra.Command{
	Use:   "check <url> [url...]",
	Short: "Check the health of one or more websites",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		results := pool.Run(args, workers)
		elapsed := time.Since(start)

		if jsonOutput {
			if err := printJSON(results); err != nil {
				return err
			}
		} else {
			for _, r := range results {
				printText(r)
				fmt.Println()
			}
			printSummary(results, elapsed, workers)
		}

		for _, r := range results {
			if r.Status != checker.StatusUp {
				os.Exit(1)
			}
		}
		return nil
	},
}

func init() {
	checkCmd.Flags().BoolVar(&jsonOutput, "json", false, "print the results as JSON")
	checkCmd.Flags().IntVarP(&workers, "workers", "w", 10, "how many sites to check at the same time")
	rootCmd.AddCommand(checkCmd)
}
