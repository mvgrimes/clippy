package cmd

import (
	"fmt"
	"net/http"

	"github.com/clip/internal/server"
	"github.com/clip/internal/store"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the clip HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")

		s := store.NewMemory()
		handler := server.New(s)

		addr := fmt.Sprintf("%s:%d", host, port)
		fmt.Printf("clip server listening on %s\n", addr)
		return http.ListenAndServe(addr, handler)
	},
}

func init() {
	serverCmd.Flags().String("host", "0.0.0.0", "Host to bind to")
	serverCmd.Flags().Int("port", 8080, "Port to listen on")
	rootCmd.AddCommand(serverCmd)
}
