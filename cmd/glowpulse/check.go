package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/YoumnaSalloum/glowpulse/internal/checker"
)

var jsonOutput bool

var checkCmd = &cobra.Command{
	Use:   "check <url>",
	Short: "Check the health of a website",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result := checker.Check(args[0])

		if jsonOutput {
			if err := printJSON(result); err != nil {
				return err
			}
		} else {
			printText(result)
		}

		if result.Status != checker.StatusUp {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	checkCmd.Flags().BoolVar(&jsonOutput, "json", false, "print the result as JSON")
	rootCmd.AddCommand(checkCmd)
}
