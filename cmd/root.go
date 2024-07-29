package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:                "backend",
	Short:              "backend",
	PersistentPreRunE:  setup,
	PersistentPostRunE: cleanup,
}

func init() {
	rootCmd.AddCommand(
		restCommand,
		kafkaCommand,
	)
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
