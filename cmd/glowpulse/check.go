package main

import (
	"errors"
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
	filePath   string
	timeout    time.Duration
)

var checkCmd = &cobra.Command{
	Use:   "check [url...]",
	Short: "Check the health of one or more websites",
	Example: `  glowpulse check https://google.com
  glowpulse check -w 20 https://google.com https://github.com
  glowpulse check -f examples/urls.txt --json`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		urls := args
		if filePath != "" {
			fromFile, err := readURLs(filePath)
			if err != nil {
				return err
			}
			urls = append(urls, fromFile...)
		}
		if len(urls) == 0 {
			return errors.New("no URLs given: pass them as arguments or use --file")
		}

		start := time.Now()
		results := pool.Run(cmd.Context(), urls, workers, timeout)
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
			if cmd.Context().Err() != nil {
				fmt.Println("⏹️  Stopped early: some sites were not checked")
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
	checkCmd.Flags().StringVarP(&filePath, "file", "f", "", "read URLs from a file, one per line")
	checkCmd.Flags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "maximum time to wait for each site")
	rootCmd.AddCommand(checkCmd)
}
