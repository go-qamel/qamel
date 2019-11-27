package cmd

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	fp "path/filepath"
	"strings"

	"github.com/go-qamel/qamel/internal/config"
	"github.com/go-qamel/qamel/internal/generator"
	"github.com/spf13/cobra"
)

func buildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build QML app",
		Args:  cobra.NoArgs,
		Run:   buildHandler,
	}

	cmd.Flags().StringP("output", "o", "", "location for executable file")
	cmd.Flags().StringP("profile", "p", "", "profile that used for building app")
	cmd.Flags().StringSliceP("tags", "t", []string{}, "space-separated list of build tags to satisfied during the build")
	cmd.Flags().Bool("copy-deps", false, "copy dependencies for app with dynamic linking")
	cmd.Flags().Bool("skip-vendoring", false, "if uses Go module, skip updating project's vendor")

	return cmd
}

func buildHandler(cmd *cobra.Command, args []string) {
	cBlueBold.Println("Starting build process.")
	fmt.Println()

	// Read flags
	buildTags, _ := cmd.Flags().GetStringSlice("tags")
	outputPath, _ := cmd.Flags().GetString("output")
	profileName, _ := cmd.Flags().GetString("profile")
	copyDependencies, _ := cmd.Flags().GetBool("copy-deps")
	skipVendoring, _ := cmd.Flags().GetBool("skip-vendoring")

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

	// Get project directory in workdir
	projectDir, err := os.Getwd()
	if err != nil {
		cRedBold.Println("Failed to get current working dir:", err)
		os.Exit(1)
	}

	// If this project uses Go module, make sure to use vendor.
	// This is because `qamel build` works by generating binding code
	// in project dir and Qamel dir as well. This is done every time
	// user build his app, because every user might have different
	// profile and config on each build. In old times, the Qamel dir
	// in $GOPATH/src is easily accessible and modified. However,
	// in current Go module, the Qamel dir in $GOPATH/pkg/mod is read
	// only, which make it impossible to generate binding code there.
	// Therefore, as workaround, Qamel in Go module *MUST* be used in
	// vendor by using `go mod vendor`.
	vendorDir := fp.Join(projectDir, "vendor", "github.com", "RadhiFadlillah", "qamel")
	goModFile := fp.Join(projectDir, "go.mod")
	usesGoModule := fileExists(goModFile)

	if usesGoModule && (!dirExists(vendorDir) || !skipVendoring) {
		fmt.Print("Generating vendor files...")

		cmdModVendor := exec.Command("go", "mod", "vendor")
		cmdOutput, err := cmdModVendor.CombinedOutput()
		if err != nil {
			fmt.Println()
			cRedBold.Println("Failed to vendor app:", err)
			cRedBold.Println(string(cmdOutput))
			os.Exit(1)
		}

		cGreen.Println("done")
	}

	// If vendor doesn't exist, uses Qamel dir in GOPATH
	qamelDir := vendorDir
	if !dirExists(qamelDir) {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = build.Default.GOPATH
		}
		qamelDir = fp.Join(gopath, "src", "github.com", "RadhiFadlillah", "qamel")
	}

	// Make sure the Qamel directory exists
	if !dirExists(qamelDir) {
		cRedBold.Println("Failed to access qamel: source directory doesn't exist")
		os.Exit(1)
	}

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
	os.Remove(fp.Join(qamelDir, "moc-viewer.h"))
	os.Remove(fp.Join(qamelDir, "moc-listmodel.h"))
	os.Remove(fp.Join(qamelDir, "moc-tablemodel.h"))
	os.Remove(fp.Join(qamelDir, "qamel_plugin_import.cpp"))
	os.Remove(fp.Join(qamelDir, "qamel_qml_plugin_import.cpp"))

	// Generate cgo file and moc for binding in qamel directory
	fmt.Print("Generating binding files...")
	filesToMoc := []string{"viewer.cpp", "listmodel.h", "tablemodel.h"}

	for _, fileToMoc := range filesToMoc {
		fileToMoc = fp.Join(qamelDir, fileToMoc)
		if !fileExists(fileToMoc) {
			continue
		}

		err = generator.CreateMocFile(profile.Moc, fileToMoc)
		if err != nil {
			fmt.Println()
			cRedBold.Printf("Failed to create moc file for %s: %v\n", fileToMoc, err)
			os.Exit(1)
		}
	}

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

	// Prepare default output path
	if outputPath == "" {
		outputPath = fp.Join(projectDir, fp.Base(projectDir))
		if profile.OS == "windows" {
			outputPath += ".exe"
		}
	}

	// Run go build
	fmt.Print("Building app...")
	cmdArgs := []string{"build", "-o", outputPath}

	if usesGoModule {
		cmdArgs = append(cmdArgs, "-mod", "vendor")
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
