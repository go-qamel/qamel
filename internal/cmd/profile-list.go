package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/RadhiFadlillah/qamel/internal/config"
	"github.com/spf13/cobra"
)

func profileListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List all existing profile",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Run:     profileListHandler,
	}
}

func profileListHandler(cmd *cobra.Command, args []string) {
	// Load existing profiles
	profiles, err := config.LoadProfiles(configPath)
	if err != nil {
		cRedBold.Println("Failed to get list of profile:", err)
		os.Exit(1)
	}

	// Get list of profile name
	names := []string{}
	for key := range profiles {
		names = append(names, key)
	}

	sort.Strings(names)
	for i := 0; i < len(names); i++ {
		fmt.Println(names[i])
	}
}
