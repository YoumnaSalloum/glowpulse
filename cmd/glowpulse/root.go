package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "glowpulse",
	Short: "GlowPulse checks if websites are up, fast and secure",
}
