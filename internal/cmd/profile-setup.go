package cmd

import (
	"bufio"
	"fmt"
	"os"
	fp "path/filepath"
	"runtime"
	"strings"

	"github.com/RadhiFadlillah/qamel/internal/config"
	"github.com/spf13/cobra"
)

func profileSetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup [profileName]",
		Short: "Set up a new or existing profile",
		Args:  cobra.MaximumNArgs(1),
		Run:   profileSetupHandler,
	}
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
	defaultObjdump := "objdump"
	defaultRes := "windres"
	if runtime.GOOS == "windows" {
		defaultGcc += ".exe"
		defaultGxx += ".exe"
		defaultObjdump += ".exe"
		defaultRes += ".exe"
	}

	fmt.Println()
	fmt.Println("Please specify the *full path* to your compiler.")
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

	// If Windows shared, specify path to objdump, which used to fetch dependencies
	objdumpPath := ""
	if targetOS == "windows" && staticMode != "y" {
		cCyanBold.Printf("Objdump (default %s)  : ", defaultObjdump)
		objdumpPath, _ = reader.ReadString('\n')
		objdumpPath = strings.TrimSpace(objdumpPath)
		if objdumpPath == "" {
			objdumpPath = defaultObjdump
		}
	}

	// If Windows, specify path to windres
	// which will be used to create icon for the executable
	windresPath := ""
	if targetOS == "windows" {
		fmt.Println()
		fmt.Println("Since you are targeting Windows, you might want to set icon for your executable file. " +
			"To do so, please specify the *full path* to windres on your system. " +
			"It's usually located in the directory where MinGW is installed.")
		fmt.Println()

		cCyanBold.Print("Path to windres : ")
		windresPath, _ = reader.ReadString('\n')
		windresPath = strings.TrimSpace(windresPath)
		if windresPath == "" {
			windresPath = defaultRes
		} else {
			if !strings.HasSuffix(windresPath, defaultRes) {
				windresPath = fp.Join(windresPath, defaultRes)
			}
			windresPath = strings.ReplaceAll(windresPath, "Program Files", "Progra~1")
			windresPath = strings.ReplaceAll(windresPath, "Program Files (x86)", "Progra~2")
			if !fileExists(windresPath) {
				cRedBold.Println("The specified path does not exist")
				os.Exit(1)
			}
		}

	}

	// Save config file
	fmt.Printf("Saving profile %s...", profileName)

	profiles, err := config.LoadProfiles(configPath)
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to save the profile:", err)
		os.Exit(1)
	}

	profiles[profileName] = config.Profile{
		OS:      targetOS,
		Arch:    targetArch,
		Static:  staticMode == "y",
		Qmake:   qmakePath,
		Moc:     mocPath,
		Rcc:     rccPath,
		Gcc:     gccPath,
		Gxx:     gxxPath,
		Windres: windresPath,
		Objdump: objdumpPath,
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
