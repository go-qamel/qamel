package main

import (
	"go/build"
	"log"
	"os"
	fp "path/filepath"

	"github.com/spf13/cobra"
)

var qamelDir = fp.Join("github.com", "RadhiFadlillah", "qamel")

func init() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	qamelDir = fp.Join(gopath, "src", qamelDir)
}

func main() {
	// Create root command
	rootCmd := &cobra.Command{
		Use:   "qamel",
		Short: "qamel is tools and binding for creating GUI app in Go using Qt + QML",
		Args:  cobra.NoArgs,
	}

	// Register sub command
	rootCmd.AddCommand(cmdSetup)

	// Execute app
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
