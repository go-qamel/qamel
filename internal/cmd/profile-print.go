package cmd

import (
	"fmt"
	"os"

	"github.com/go-qamel/qamel/internal/config"
	"github.com/spf13/cobra"
)

func profilePrintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "print [profileName]",
		Short: "Print data of specified profile",
		Args:  cobra.MaximumNArgs(1),
		Run:   profilePrintHandler,
	}
}

func profilePrintHandler(cmd *cobra.Command, args []string) {
	// Get profile name
	profileName := "default"
	if len(args) == 1 {
		profileName = args[0]
	}

	// Load existing profiles
	profiles, err := config.LoadProfiles(configPath)
	if err != nil {
		cRedBold.Println("Failed to print the profile:", err)
		os.Exit(1)
	}

	// Save the profiles with removed item
	profile, ok := profiles[profileName]
	if !ok {
		cRedBold.Printf("Failed to print the profile: %s not exists\n", profileName)
		os.Exit(1)
	}

	cBlueBold.Printf("Details of profile %s\n", profileName)
	cCyanBold.Print("OS      : ")
	fmt.Println(profile.OS)
	cCyanBold.Print("Arch    : ")
	fmt.Println(profile.Arch)
	cCyanBold.Print("Static  : ")
	fmt.Println(profile.Static)
	cCyanBold.Print("Qmake   : ")
	fmt.Println(profile.Qmake)
	cCyanBold.Print("Moc     : ")
	fmt.Println(profile.Moc)
	cCyanBold.Print("Rcc     : ")
	fmt.Println(profile.Rcc)
	cCyanBold.Print("Gcc     : ")
	fmt.Println(profile.Gcc)
	cCyanBold.Print("G++     : ")
	fmt.Println(profile.Gxx)

	if profile.OS == "windows" {
		if !profile.Static {
			cCyanBold.Print("Objdump : ")
			fmt.Println(profile.Objdump)
		}

		cCyanBold.Print("Windres : ")
		fmt.Println(profile.Windres)
	}
}
