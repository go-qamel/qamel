package main

import (
	"fmt"
	"os"
	"os/exec"
	fp "path/filepath"
	"strings"

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
	cmdBuild.Flags().StringSliceP("tags", "t", []string{}, "space-separated list of build tags to satisfied during the build")
}

func buildHandler(cmd *cobra.Command, args []string) {
	cBlueBold.Println("Starting build process.")
	fmt.Println()

	// Read flags
	buildTags, _ := cmd.Flags().GetStringSlice("tags")
	outputPath, _ := cmd.Flags().GetString("output")

	// Get destination directory
	dstDir := ""
	if len(args) == 1 {
		dstDir = args[0]
	}

	// If destination directory is empty, use current working directory
	// Else, make sure destination dir is exists
	if dstDir == "" {
		var err error
		dstDir, err = os.Getwd()
		if err != nil {
			cRedBold.Println("Failed to get current working dir:", err)
			return
		}
	} else if !dirExists(dstDir) {
		cRedBold.Println("Destination directory doesn't exist")
		return
	}

	// Make sure destination dir is absolute
	dstDir, err := fp.Abs(dstDir)
	if err != nil {
		cRedBold.Println("Failed to get destination dir:", err)
		return
	}

	// Load config file
	fmt.Print("Load config file...")
	cfg, err := loadConfig()
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to load config file:", err)
		cRedBold.Println("You might need to run `qamel setup` again.")
		return
	}
	cGreen.Println("done")

	// Remove old qamel files
	fmt.Print("Removing old build files...")
	err = removeQamelFiles(dstDir)
	if err != nil {
		cRedBold.Println("Failed to remove old build files:", err)
		return
	}
	cGreen.Println("done")

	// Create cgo flags
	fmt.Print("Generating cgo flags...")
	cgoFlags, err := generator.CreateCgoFlags(cfg.Qmake)
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to create cgo flags:", err)
		return
	}
	cGreen.Println("done")

	// Create rcc file
	fmt.Print("Generating Qt resource file...")
	err = generator.CreateRccFile(cfg.Rcc, dstDir, cgoFlags)
	if err != nil {
		cYellow.Println(err)
	} else {
		cGreen.Println("done")
	}

	// Generate code for QmlObject structs
	fmt.Print("Generating code for QML objects...")
	errs := generator.CreateQmlObjectCode(cfg.Moc, dstDir, cgoFlags, buildTags...)
	if errs != nil && len(errs) > 0 {
		fmt.Println()
		for _, err := range errs {
			cRedBold.Println("Failed:", err)
		}
		return
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

	cmdGo := exec.Command("go", cmdArgs...)
	cmdGo.Dir = dstDir
	cmdGo.Env = append(os.Environ(),
		`CGO_CFLAGS_ALLOW=.*`,
		`CGO_CXXFLAGS_ALLOW=.*`,
		`CGO_LDFLAGS_ALLOW=.*`)

	cmdOutput, err := cmdGo.CombinedOutput()
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to build app:", err)
		cRedBold.Println(string(cmdOutput))
		return
	}
	cGreen.Println("done")

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

		for i := range prefixes {
			if strings.HasPrefix(info.Name(), prefixes[i]) {
				if err = os.Remove(path); err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}
