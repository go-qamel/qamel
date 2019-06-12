package cmd

import (
	"os"

	"github.com/RadhiFadlillah/qamel/internal/config"
	"github.com/spf13/cobra"
)

func profileDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Short:   "Delete an existing profile",
		Args:    cobra.ExactArgs(1),
		Aliases: []string{"rm"},
		Run:     profileDeleteHandler,
	}
}

func profileDeleteHandler(cmd *cobra.Command, args []string) {
	// Get profile name
	profileName := args[0]

	// Load existing profiles
	profiles, err := config.LoadProfiles(configPath)
	if err != nil {
		cRedBold.Println("Failed to remove the profile:", err)
		os.Exit(1)
	}

	// Save the profiles with removed item
	delete(profiles, profileName)
	err = config.SaveProfiles(configPath, profiles)
	if err != nil {
		cRedBold.Println("Failed to remove the profile:", err)
		os.Exit(1)
	}

	cBlueBold.Printf("Profile %s has been removed.\n", profileName)
}
