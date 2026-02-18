package cmd

import (
	"fmt"
	"net/http"

	"github.com/mvgrimes/clipr/internal/server"
	"github.com/mvgrimes/clipr/internal/store"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the clipr HTTP server",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindEnv(cmd, "host", "CLIPR_HOST")
		bindEnv(cmd, "port", "CLIPR_PORT")
		bindEnv(cmd, "max-per-clip", "CLIPR_MAX_PER_CLIP")
		bindEnv(cmd, "max-clips", "CLIPR_MAX_CLIPS")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		maxPerClipStr, _ := cmd.Flags().GetString("max-per-clip")
		maxClipsStr, _ := cmd.Flags().GetString("max-clips")

		var opts []store.MemoryOption
		if maxPerClipStr != "" {
			n, err := parseSize(maxPerClipStr)
			if err != nil {
				return fmt.Errorf("invalid --max-per-clip: %w", err)
			}
			opts = append(opts, store.WithMaxPerClip(n))
		}
		if maxClipsStr != "" {
			n, err := parseSize(maxClipsStr)
			if err != nil {
				return fmt.Errorf("invalid --max-clips: %w", err)
			}
			opts = append(opts, store.WithMaxTotal(n))
		}

		s := store.NewMemory(opts...)
		handler := server.New(s)

		addr := fmt.Sprintf("%s:%d", host, port)
		fmt.Printf("clipr server listening on %s\n", addr)
		return http.ListenAndServe(addr, handler)
	},
}

func init() {
	serverCmd.Flags().String("host", "0.0.0.0", "Host to bind to")
	serverCmd.Flags().Int("port", 8080, "Port to listen on")
	serverCmd.Flags().String("max-per-clip", "", "Max size per clip (e.g. 1M, 512K)")
	serverCmd.Flags().String("max-clips", "", "Max total size of all clips (e.g. 100M, 1G)")
	rootCmd.AddCommand(serverCmd)
}
