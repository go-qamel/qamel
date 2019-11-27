package cmd

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"os/user"
	fp "path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

func dockerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docker [target]",
		Run:   dockerHandler,
		Args:  cobra.ExactArgs(1),
		Short: "Build QML app using Docker image",
		Long: "Build QML app using Docker image.\nPossible values are " +
			`"linux", "linux-static", "win32", "win32-static", "win64" and "win64-static".`,
	}

	cmd.Flags().StringP("output", "o", "", "location for executable file")
	cmd.Flags().StringSliceP("tags", "t", []string{}, "space-separated list of build tags to satisfied during the build")
	cmd.Flags().Bool("copy-deps", false, "copy dependencies for app with dynamic linking")
	cmd.Flags().Bool("skip-vendoring", false, "if uses Go module, skip updating project's vendor")

	return cmd
}

func dockerHandler(cmd *cobra.Command, args []string) {
	cBlueBold.Println("Run `qamel build` from Docker image.")

	// Read flags
	buildTags, _ := cmd.Flags().GetStringSlice("tags")
	outputPath, _ := cmd.Flags().GetString("output")
	copyDependencies, _ := cmd.Flags().GetBool("copy-deps")
	skipVendoring, _ := cmd.Flags().GetBool("skip-vendoring")

	// Get target name
	target := args[0]
	switch target {
	case "linux", "linux-static",
		"win32", "win32-static",
		"win64", "win64-static":
	default:
		cRedBold.Printf("Target %s is not supported.\n", target)
		os.Exit(1)
	}

	// Get gopath
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	// Get project directory from current working dir
	projectDir, err := os.Getwd()
	if err != nil {
		cRedBold.Println("Failed to get current working dir:", err)
		os.Exit(1)
	}

	// If this project uses Go module, vendor it first before
	// passing it to Docker
	vendorDir := fp.Join(projectDir, "vendor", "github.com", "go-qamel", "qamel")
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

	// Get docker user from current active user
	currentUser, err := user.Current()
	if err != nil {
		cRedBold.Println("Failed to get user's data:", err)
		os.Exit(1)
	}

	uid := currentUser.Uid
	gid := currentUser.Gid
	dockerUser := fmt.Sprintf("%s:%s", uid, gid)

	if runtime.GOOS == "windows" {
		uidParts := strings.Split(uid, "-")
		dockerUser = uidParts[len(uidParts)-1]
	}

	// Create directory for Go build's cache
	goCacheDir := fp.Join(projectDir, ".qamel-cache", target, "go")

	err = os.MkdirAll(goCacheDir, os.ModePerm)
	if err != nil {
		cRedBold.Println("Failed to create cache directory for Go build:", err)
		os.Exit(1)
	}

	// If target is Linux, create ccache dir as well
	ccacheDir := ""
	if strings.HasPrefix(target, "linux") {
		ccacheDir = fp.Join(projectDir, ".qamel-cache", target, "ccache")

		err = os.MkdirAll(ccacheDir, os.ModePerm)
		if err != nil {
			cRedBold.Println("Failed to create cache directory for gcc build:", err)
			os.Exit(1)
		}
	}

	// Prepare docker arguments
	dockerGopath := unixJoinPath("/", "home", "user", "go")
	dockerProjectDir := unixJoinPath(dockerGopath, "src", fp.Base(projectDir))
	dockerGoCacheDir := unixJoinPath("/", "home", "user", ".cache", "go-build")

	dockerBindProject := fmt.Sprintf(`type=bind,src=%s,dst=%s`,
		projectDir, dockerProjectDir)
	dockerBindGoSrc := fmt.Sprintf(`type=bind,src=%s,dst=%s`,
		unixJoinPath(gopath, "src"),
		unixJoinPath(dockerGopath, "src"))
	dockerBindGoCache := fmt.Sprintf(`type=bind,src=%s,dst=%s`,
		unixJoinPath(goCacheDir), dockerGoCacheDir)

	dockerArgs := []string{"run", "--rm",
		"--attach", "stdout",
		"--attach", "stderr",
		"--user", dockerUser,
		"--workdir", dockerProjectDir,
		"--mount", dockerBindProject,
		"--mount", dockerBindGoSrc,
		"--mount", dockerBindGoCache}

	if ccacheDir != "" {
		dockerCcacheDir := unixJoinPath("/", "home", "user", ".ccache")
		dockerBindCcache := fmt.Sprintf(`type=bind,src=%s,dst=%s`,
			unixJoinPath(ccacheDir), dockerCcacheDir)

		dockerArgs = append(dockerArgs, "--mount", dockerBindCcache)
	}

	// If uses go module, set env GO111MODULE to on
	if usesGoModule {
		dockerArgs = append(dockerArgs, "--env", "GO111MODULE=on")
	}

	// Finally, set image name
	dockerImageName := fmt.Sprintf("radhifadlillah/qamel:%s", target)
	dockerArgs = append(dockerArgs, dockerImageName)

	// Add qamel arguments
	dockerArgs = append(dockerArgs, "--skip-vendoring")

	if outputPath != "" {
		dockerArgs = append(dockerArgs, "-o", outputPath)
	}

	if len(buildTags) > 0 {
		dockerArgs = append(dockerArgs, "-t")
		dockerArgs = append(dockerArgs, buildTags...)
	}

	if copyDependencies {
		dockerArgs = append(dockerArgs, "--copy-deps")
	}

	// Run docker
	cmdDocker := exec.Command("docker", dockerArgs...)
	cmdDocker.Stdout = os.Stdout
	cmdDocker.Stderr = os.Stderr

	err = cmdDocker.Start()
	if err != nil {
		cRedBold.Println("Failed to start Docker:", err)
		os.Exit(1)
	}

	err = cmdDocker.Wait()
	if err != nil {
		cRedBold.Println("Failed to build app using Docker:", err)
		os.Exit(1)
	}
}
