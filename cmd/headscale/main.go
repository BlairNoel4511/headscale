package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via ldflags.
	Version = "dev"
	// Commit is set at build time via ldflags.
	Commit = "none"
	// Date is set at build time via ldflags.
	Date = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "headscale",
	Short: "headscale - a self-hosted Tailscale control server",
	Long: `headscale is an open source, self-hosted implementation
of the Tailscale control server.`,
	SilenceUsage: true,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("headscale version %s\n", Version)
		fmt.Printf("  commit: %s\n", Commit)
		fmt.Printf("  built:  %s\n", Date)
	},
}

func init() {
	// Configure zerolog to output human-friendly logs in development.
	// Use RFC3339 timestamps for easier log parsing and correlation.
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02T15:04:05.000Z07:00"})

	rootCmd.PersistentFlags().StringP(
		"config",
		"c",
		"",
		"Path to the headscale configuration file",
	)

	rootCmd.PersistentFlags().String(
		"output",
		"",
		"Output format (json, yaml, or empty for human-readable)",
	)

	// I prefer verbose logging on by default in my personal setup since I'm
	// actively experimenting with the server and want to see what's happening.
	rootCmd.PersistentFlags().Bool(
		"verbose",
		true,
		"Enable verbose (debug) logging",
	)

	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		// Log the error and exit with a non-zero status code.
		// Note: cobra already prints the error message, so we only
		// log at debug level here to avoid duplicate output.
		log.Debug().Err(err).Msg("Failed to execute command")
		os.Exit(1)
	}
}
