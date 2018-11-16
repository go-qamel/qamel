package main

import (
	"go/build"
	"log"
	"os"
	fp "path/filepath"

	ap "github.com/muesli/go-app-paths"
	"github.com/spf13/cobra"
)

var (
	qamelDir   = fp.Join("github.com", "RadhiFadlillah", "qamel")
	configPath = "config.json"
)

func init() {
	// Get qamel directory
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	qamelDir = fp.Join(gopath, "src", qamelDir)

	// Get config path in ${XDG_CONFIG_HOME}/qamel/config.json
	var err error
	userScope := ap.NewScope(ap.User, "qamel", "qamel")
	configPath, err = userScope.ConfigPath("config.json")
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	// Create root command
	rootCmd := &cobra.Command{
		Use:   "qamel",
		Short: "qamel is tools and binding for creating GUI app in Go using Qt + QML",
		Args:  cobra.NoArgs,
	}

	// Register sub command
	rootCmd.AddCommand(cmdProfile, cmdBuild)

	// Execute app
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
