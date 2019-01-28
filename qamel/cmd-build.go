package main

import (
	"fmt"
	"os"
	"os/exec"
	fp "path/filepath"
	"strings"

	"github.com/RadhiFadlillah/qamel/qamel/config"
	"github.com/RadhiFadlillah/qamel/qamel/generator"
	"github.com/spf13/cobra"
)

var cmdBuild = &cobra.Command{
	Use:   "build",
	Short: "Build QML app",
	Args:  cobra.MaximumNArgs(1),
	Run:   buildHandler,
}

func init() {
	cmdBuild.Flags().StringP("output", "o", "", "location for executable file")
	cmdBuild.Flags().StringP("profile", "p", "", "profile that used for building app")
	cmdBuild.Flags().StringSliceP("tags", "t", []string{}, "space-separated list of build tags to satisfied during the build")
	cmdBuild.Flags().Bool("copy-deps", false, "copy dependencies for app with dynamic linking")
}

func buildHandler(cmd *cobra.Command, args []string) {
	cBlueBold.Println("Starting build process.")
	fmt.Println()

	// Read flags
	buildTags, _ := cmd.Flags().GetStringSlice("tags")
	outputPath, _ := cmd.Flags().GetString("output")
	profileName, _ := cmd.Flags().GetString("profile")
	copyDependencies, _ := cmd.Flags().GetBool("copy-deps")

	// Get project directory
	projectDir := ""
	if len(args) == 1 {
		projectDir = args[0]
	}

	// If project directory is empty, use current working directory
	// Else, make sure project dir is exists
	if projectDir == "" {
		var err error
		projectDir, err = os.Getwd()
		if err != nil {
			cRedBold.Println("Failed to get current working dir:", err)
			os.Exit(1)
		}
	} else if !dirExists(projectDir) {
		cRedBold.Println("Destination directory doesn't exist")
		os.Exit(1)
	}

	// Make sure project dir is absolute
	projectDir, err := fp.Abs(projectDir)
	if err != nil {
		cRedBold.Println("Failed to get project dir:", err)
		os.Exit(1)
	}

	// Load config file
	fmt.Print("Load config file...")

	profileName = strings.TrimSpace(profileName)
	if profileName == "" {
		profileName = "default"
	}

	profile, err := config.LoadProfile(configPath, profileName)
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to load profile file:", err)
		cRedBold.Println("You might need to run `qamel profile setup` again.")
		os.Exit(1)
	}
	cGreen.Println("done")

	// Remove old qamel files
	fmt.Print("Removing old build files...")
	err = removeQamelFiles(projectDir)
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to remove old build files:", err)
		os.Exit(1)
	}
	cGreen.Println("done")

	os.Remove(fp.Join(projectDir, "qamel-icon.syso"))
	os.Remove(fp.Join(qamelDir, "qamel_plugin_import.cpp"))
	os.Remove(fp.Join(qamelDir, "qamel_qml_plugin_import.cpp"))

	// Generate cgo file for binding in qamel directory
	fmt.Print("Generating cgo file...")
	err = generator.CreateCgoFile(profile, qamelDir, "")
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to create cgo file:", err)
		os.Exit(1)
	}
	cGreen.Println("done")

	// Create rcc file
	fmt.Print("Generating Qt resource file...")
	err = generator.CreateRccFile(profile, projectDir)
	if err != nil {
		if err == generator.ErrNoResourceDir {
			cYellow.Println(err)
		} else {
			fmt.Println()
			cRedBold.Println("Failed to create Qt resource file:", err)
			os.Exit(1)
		}
	} else {
		cGreen.Println("done")
	}

	// Create syso file
	if profile.OS == "windows" && profile.Windres != "" {
		fmt.Print("Generating syso file...")
		err = generator.CreateSysoFile(profile, projectDir)
		if err != nil {
			if err == generator.ErrNoIcon {
				cYellow.Println(err)
			} else {
				fmt.Println()
				cRedBold.Println("Failed to create syso file:", err)
				os.Exit(1)
			}
		} else {
			cGreen.Println("done")
		}
	}

	// Generate code for QmlObject structs
	fmt.Print("Generating code for QML objects...")
	errs := generator.CreateQmlObjectCode(profile, projectDir, buildTags...)
	if errs != nil && len(errs) > 0 {
		fmt.Println()
		for _, err := range errs {
			cRedBold.Println("Failed:", err)
		}
		os.Exit(1)
	}
	cGreen.Println("done")

	// Run go build
	fmt.Print("Building app...")
	cmdArgs := []string{"build"}

	if outputPath != "" {
		cmdArgs = append(cmdArgs, "-o", outputPath)
	}

	if len(buildTags) > 0 {
		cmdArgs = append(cmdArgs, "-tags")
		cmdArgs = append(cmdArgs, buildTags...)
	}

	ldFlags := "all=-s -w"
	if profile.OS == "windows" {
		ldFlags += " -H=windowsgui"
	}

	cmdArgs = append(cmdArgs, "-ldflags")
	cmdArgs = append(cmdArgs, ldFlags)

	cmdGo := exec.Command("go", cmdArgs...)
	cmdGo.Dir = projectDir
	cmdGo.Env = append(os.Environ(),
		`CGO_ENABLED=1`,
		`CGO_CFLAGS_ALLOW=.*`,
		`CGO_CXXFLAGS_ALLOW=.*`,
		`CGO_LDFLAGS_ALLOW=.*`,
		"GOOS="+profile.OS,
		"GOARCH="+profile.Arch,
		"CC="+profile.Gcc,
		"CXX="+profile.Gxx)

	cmdOutput, err := cmdGo.CombinedOutput()
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to build app:", err)
		cRedBold.Println(string(cmdOutput))
		os.Exit(1)
	}
	cGreen.Println("done")

	// If it's shared mode, copy dependencies
	if !profile.Static && copyDependencies {
		fmt.Print("Copying dependencies...")
		err = generator.CopyDependencies(profile, projectDir, outputPath)
		if err != nil {
			fmt.Println()
			cRedBold.Println("Failed to copy dependencies:", err)
			os.Exit(1)
		}
		cGreen.Println("done")
	}

	// Build finished
	fmt.Println()
	cBlueBold.Println("Build finished succesfully.")
}

// removeQamelFiles remove old generated qamel files in the specified dir
func removeQamelFiles(rootDir string) error {
	prefixes := []string{"moc-qamel-", "qamel-"}

	err := fp.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return fp.SkipDir
			}
			return nil
		}

		switch fileExt := fp.Ext(info.Name()); fileExt {
		case ".h", ".go", ".cpp":
		default:
			return nil
		}

		if strings.HasSuffix(info.Name(), "_plugin_import.cpp") {
			return os.Remove(path)
		}

		for i := range prefixes {
			if strings.HasPrefix(info.Name(), prefixes[i]) {
				return os.Remove(path)
			}
		}

		return nil
	})

	return err
}
