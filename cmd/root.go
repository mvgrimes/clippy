package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is the current version of clip.
var Version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "clip",
	Short: "A simple pastebin service and CLI client",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of clip",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("clip v%s\n", Version)
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
