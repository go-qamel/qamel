package cmd

import (
	"github.com/spf13/cobra"
)

func profileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage profiles for QML's binding",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		profileSetupCmd(),
		profileDeleteCmd(),
		profileListCmd(),
		profilePrintCmd(),
	)

	return cmd
}
