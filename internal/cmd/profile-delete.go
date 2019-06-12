package cmd

import (
	"os"
	"strings"

	"github.com/RadhiFadlillah/qamel/internal/config"
	"github.com/spf13/cobra"
)

func profileDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Short:   "Delete an existing profile",
		Args:    cobra.MinimumNArgs(1),
		Aliases: []string{"rm"},
		Run:     profileDeleteHandler,
	}
}

func profileDeleteHandler(cmd *cobra.Command, args []string) {
	// Load existing profiles
	profiles, err := config.LoadProfiles(configPath)
	if err != nil {
		cRedBold.Println("Failed to remove the profile:", err)
		os.Exit(1)
	}

	// Delete each profile
	for _, profileName := range args {
		delete(profiles, profileName)
	}

	// Save the new profiles after removal
	err = config.SaveProfiles(configPath, profiles)
	if err != nil {
		cRedBold.Println("Failed to remove the profile:", err)
		os.Exit(1)
	}

	logFormat := "Profile %s has been removed\n"
	if len(args) > 1 {
		logFormat = "Profiles [%s] have been removed\n"
	}

	cBlueBold.Printf(logFormat, strings.Join(args, ","))
}
