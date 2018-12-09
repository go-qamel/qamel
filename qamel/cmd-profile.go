package main

import (
	"bufio"
	"fmt"
	"os"
	fp "path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/RadhiFadlillah/qamel/qamel/config"
	"github.com/RadhiFadlillah/qamel/qamel/generator"
	"github.com/spf13/cobra"
)

var cmdProfile = &cobra.Command{
	Use:   "profile",
	Short: "Manage profiles for QML's binding",
	Args:  cobra.NoArgs,
}

var cmdProfileSetup = &cobra.Command{
	Use:   "setup [profileName]",
	Short: "Set up a new or existing profile",
	Args:  cobra.MaximumNArgs(1),
	Run:   profileSetupHandler,
}

var cmdProfileDelete = &cobra.Command{
	Use:     "delete",
	Short:   "Delete an existing profile",
	Args:    cobra.ExactArgs(1),
	Run:     profileDeleteHandler,
	Aliases: []string{"rm"},
}

var cmdProfileList = &cobra.Command{
	Use:     "list",
	Short:   "List all existing profile",
	Args:    cobra.NoArgs,
	Run:     profileListHandler,
	Aliases: []string{"ls"},
}

var cmdProfilePrint = &cobra.Command{
	Use:   "print [profileName]",
	Short: "Print data of specified profile",
	Args:  cobra.MaximumNArgs(1),
	Run:   profilePrintHandler,
}

func init() {
	cmdProfile.AddCommand(cmdProfileSetup, cmdProfileDelete, cmdProfileList, cmdProfilePrint)
}

func profileSetupHandler(cmd *cobra.Command, args []string) {
	// Get profile name
	profileName := "default"
	if len(args) > 0 {
		profileName = args[0]
	}

	// Define input reader
	reader := bufio.NewReader(os.Stdin)

	// Fetch target OS
	cBlueBold.Println("Thanks for using qamel, QML's binding for Go.")
	fmt.Println()
	fmt.Println("Please specify the target OS for this profile. " +
		`Possible values are "windows", "linux" and "darwin". ` + "\n" +
		"Keep it empty to use your system OS.")
	fmt.Println()

	cCyanBold.Printf("Target OS (default %s) : ", runtime.GOOS)
	targetOS, _ := reader.ReadString('\n')
	targetOS = strings.TrimSpace(targetOS)
	if targetOS == "" {
		targetOS = runtime.GOOS
	}

	switch targetOS {
	case "linux", "windows", "darwin":
	default:
		cRedBold.Printf("OS %s is not supported\n", targetOS)
		os.Exit(1)
	}

	// Fetch target architecture
	fmt.Println()
	fmt.Println("Please specify the target architecture for this profile. " +
		`Possible values are "386" and "amd64".` + "\n" +
		"Keep it empty to use your system architecture.")
	fmt.Println()

	cCyanBold.Printf("Target arch (default %s) : ", runtime.GOARCH)
	targetArch, _ := reader.ReadString('\n')
	targetArch = strings.TrimSpace(targetArch)
	if targetArch == "" {
		targetArch = runtime.GOARCH
	}

	switch targetArch {
	case "386", "amd64":
	default:
		cRedBold.Printf("Architecture %s is not supported\n", targetArch)
		os.Exit(1)
	}

	// Fetch build mode
	fmt.Println()
	fmt.Println("Please specify whether this profile used to build static or shared app.")
	fmt.Println()

	cCyanBold.Print("Use static mode (y/n, default n) : ")
	staticMode, _ := reader.ReadString('\n')
	staticMode = strings.TrimSpace(staticMode)
	staticMode = strings.ToLower(staticMode)
	if staticMode == "" {
		staticMode = "n"
	}

	if staticMode != "y" && staticMode != "n" {
		cRedBold.Println("Input value is not valid")
		os.Exit(1)
	}

	// Fetch path to Qt's directory and tools
	fmt.Println()
	fmt.Println("Please specify the *full path* to your Qt's tools directory.")
	fmt.Println("This might be different depending on your platform or your target. " +
		"For example, in Linux with Qt 5.11.1, " +
		"the tools are located in $HOME/Qt5.11.1/5.11.1/gcc_64/bin/")
	fmt.Println()

	cCyanBold.Print("Qt tools dir : ")
	qtDir, _ := reader.ReadString('\n')
	qtDir = strings.TrimSpace(qtDir)
	if !dirExists(qtDir) {
		cRedBold.Println("The specified directory does not exist")
		os.Exit(1)
	}

	// Make sure qmake, moc, and rcc is exists
	qmakePath := fp.Join(qtDir, "qmake")
	mocPath := fp.Join(qtDir, "moc")
	rccPath := fp.Join(qtDir, "rcc")
	if runtime.GOOS == "windows" {
		qmakePath += ".exe"
		mocPath += ".exe"
		rccPath += ".exe"
	}

	qmakeExists := fileExists(qmakePath)
	mocExists := fileExists(mocPath)
	rccExists := fileExists(rccPath)

	cCyanBold.Print("qmake        : ")
	if qmakeExists {
		cGreen.Println("found")
	} else {
		cRed.Println("not found")
	}

	cCyanBold.Print("moc          : ")
	if mocExists {
		cGreen.Println("found")
	} else {
		cRed.Println("not found")
	}

	cCyanBold.Print("rcc          : ")
	if rccExists {
		cGreen.Println("found")
	} else {
		cRed.Println("not found")
	}

	if !qmakeExists || !mocExists || !rccExists {
		fmt.Println()
		fmt.Println("Unable to find some of the tools. Please specify the *full path* to it manually.")
		fmt.Println()

		if !qmakeExists {
			cCyanBold.Print("Path to qmake : ")
			qmakePath, _ = reader.ReadString('\n')
			qmakePath = strings.TrimSpace(qmakePath)
			if !fileExists(qmakePath) {
				cRedBold.Println("The specified path does not exist")
				os.Exit(1)
			}
		}

		if !mocExists {
			cCyanBold.Print("Path to moc   : ")
			mocPath, _ = reader.ReadString('\n')
			mocPath = strings.TrimSpace(mocPath)
			if !fileExists(mocPath) {
				cRedBold.Println("The specified path does not exist")
				os.Exit(1)
			}
		}

		if !rccExists {
			cCyanBold.Print("Path to rcc   : ")
			rccPath, _ = reader.ReadString('\n')
			rccPath = strings.TrimSpace(rccPath)
			if !fileExists(rccPath) {
				cRedBold.Println("The specified path does not exist")
				os.Exit(1)
			}
		}
	}

	// Fetch custom C and C++ compiler
	defaultGcc := "gcc"
	defaultGxx := "g++"
	if runtime.GOOS == "windows" {
		defaultGcc += ".exe"
		defaultGxx += ".exe"
	}

	fmt.Println()
	fmt.Println("Please specify the *full path* to your C and C++ compiler.")
	fmt.Println("Keep it empty to use the default compiler on your system.")
	fmt.Println()

	cCyanBold.Printf("C compiler (default %s)   : ", defaultGcc)
	gccPath, _ := reader.ReadString('\n')
	gccPath = strings.TrimSpace(gccPath)
	if gccPath == "" {
		gccPath = defaultGcc
	}

	cCyanBold.Printf("C++ compiler (default %s) : ", defaultGxx)
	gxxPath, _ := reader.ReadString('\n')
	gxxPath = strings.TrimSpace(gxxPath)
	if gxxPath == "" {
		gxxPath = defaultGxx
	}

	// Generating moc file for viewer
	fmt.Println()
	fmt.Print("Generating some code for binding...")

	err := generator.CreateMocFile(mocPath, fp.Join(qamelDir, "viewer.cpp"))
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to create moc file for viewer:", err)
		os.Exit(1)
	}

	cGreen.Println("done")

	// Save config file
	fmt.Printf("Saving profile %s...", profileName)

	profiles, err := config.LoadProfiles(configPath)
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to save the profile:", err)
		os.Exit(1)
	}

	profiles[profileName] = config.Profile{
		OS:     targetOS,
		Arch:   targetArch,
		Static: staticMode == "y",
		Qmake:  qmakePath,
		Moc:    mocPath,
		Rcc:    rccPath,
		Gcc:    gccPath,
		Gxx:    gxxPath,
	}

	err = config.SaveProfiles(configPath, profiles)
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to save the profile:", err)
		os.Exit(1)
	}

	cGreen.Println("done")

	// Setup finished
	fmt.Println()
	cBlueBold.Println("Setup finished.")
	cBlueBold.Println("Now you can get started on your own QML app.")
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
	cCyanBold.Print("OS     : ")
	fmt.Println(profile.OS)
	cCyanBold.Print("Arch   : ")
	fmt.Println(profile.Arch)
	cCyanBold.Print("Static : ")
	fmt.Println(profile.Static)
	cCyanBold.Print("Qmake  : ")
	fmt.Println(profile.Qmake)
	cCyanBold.Print("Moc    : ")
	fmt.Println(profile.Moc)
	cCyanBold.Print("Rcc    : ")
	fmt.Println(profile.Rcc)
	cCyanBold.Print("Gcc    : ")
	fmt.Println(profile.Gcc)
	cCyanBold.Print("G++    : ")
	fmt.Println(profile.Gxx)
}
