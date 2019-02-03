package main

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

var cmdDocker = &cobra.Command{
	Use:   "docker [target]",
	Run:   dockerHandler,
	Args:  cobra.ExactArgs(1),
	Short: "Build QML app using Docker image",
	Long: "Build QML app using Docker image.\nPossible values are " +
		`"linux", "linux-static", "win32", "win32-static", "win64" and "win64-static".`,
}

func init() {
	cmdDocker.Flags().StringP("output", "o", "", "location for executable file")
	cmdDocker.Flags().StringSliceP("tags", "t", []string{}, "space-separated list of build tags to satisfied during the build")
	cmdDocker.Flags().Bool("copy-deps", false, "copy dependencies for app with dynamic linking")
}

func dockerHandler(cmd *cobra.Command, args []string) {
	cBlueBold.Println("Run `qamel build` from Docker image.")

	// Read flags
	buildTags, _ := cmd.Flags().GetStringSlice("tags")
	outputPath, _ := cmd.Flags().GetString("output")
	copyDependencies, _ := cmd.Flags().GetBool("copy-deps")

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

	// Get relative project dir from $GOPATH
	projectDir, err := os.Getwd()
	if err != nil {
		cRedBold.Println("Failed to get current working dir:", err)
		os.Exit(1)
	}

	hostGopath := build.Default.GOPATH
	projectDir, err = fp.Rel(hostGopath, projectDir)
	if err != nil {
		cRedBold.Println("Failed to get relative path of project dir:", err)
		os.Exit(1)
	}

	// Create docker user
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

	// Prepare docker arguments
	dockerGopath := fp.Join("/", "home", "user", "go")
	dockerWorkdir := fp.Join(dockerGopath, projectDir)
	dockerBindSrc := fp.Join(hostGopath, "src")
	dockerBindDst := fp.Join(dockerGopath, "src")
	dockerBindFs := fmt.Sprintf("type=bind,src=%s,dst=%s", dockerBindSrc, dockerBindDst)
	dockerImageName := fmt.Sprintf("radhifadlillah/qamel:%s", target)

	dockerArgs := []string{"run", "--rm",
		"--attach", "stdout",
		"--attach", "stderr",
		"--user", dockerUser,
		"--workdir", dockerWorkdir,
		"--mount", dockerBindFs,
		dockerImageName}

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
