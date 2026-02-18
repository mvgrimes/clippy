package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var pasteCmd = &cobra.Command{
	Use:   "paste [key]",
	Short: "Store a paste from stdin",
	Args:  cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		bindEnv(cmd, "server", "CLIPR_SERVER")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, _ := cmd.Flags().GetString("server")
		path := "/@"
		if len(args) == 1 {
			path = "/@/" + args[0]
		}

		resp, err := http.Post(srv+path, "application/octet-stream", os.Stdin)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("server returned %d: %s", resp.StatusCode, body)
		}
		fmt.Print(string(body))
		return nil
	},
}

func init() {
	pasteCmd.Flags().String("server", "http://localhost:8080", "Server URL")
	rootCmd.AddCommand(pasteCmd)
}
