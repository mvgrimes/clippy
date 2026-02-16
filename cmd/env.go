package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// bindEnv sets a flag's default from an environment variable if the flag
// was not explicitly set on the command line.
func bindEnv(cmd *cobra.Command, flagName, envVar string) {
	f := cmd.Flags().Lookup(flagName)
	if f == nil {
		return
	}
	if !f.Changed {
		if v, ok := os.LookupEnv(envVar); ok {
			f.Value.Set(v)
		}
	}
}
