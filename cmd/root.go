package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is the current version of clippy.
var Version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "clippy",
	Short: "A simple pastebin service and CLI client",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of clippy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("clippy v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
