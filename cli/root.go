package cli

import (
	"github.com/mlange-42/xwrd/core"
	"github.com/spf13/cobra"
)

// RootCommand sets up the CLI
func RootCommand(config *core.Config, version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "xwrd",
		Short:         "Words tool",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	root.AddCommand(anagramCommand(config))
	root.AddCommand(matchCommand(config))

	return root
}
