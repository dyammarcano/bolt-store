package cmd

import (
	"context"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bolt-store",
	Short: "A simple CLI to interact with a bolt store",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Hello, World!")
		return nil
	},
}

func Execute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cobra.CheckErr(rootCmd.ExecuteContext(ctx))
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
