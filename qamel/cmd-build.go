package main

import (
	"fmt"
	"os"
	fp "path/filepath"

	"github.com/RadhiFadlillah/qamel/qamel/generator"
	"github.com/spf13/cobra"
)

var cmdBuild = &cobra.Command{
	Use:   "build",
	Short: "Build QML app",
	Args:  cobra.MaximumNArgs(1),
	Run:   buildHandler,
}

func buildHandler(cmd *cobra.Command, args []string) {
	cBlueBold.Println("Starting build process.")
	fmt.Println()

	// Get destination directory
	dstDir := ""
	if len(args) == 1 {
		dstDir = args[0]
	}

	// If destination directory is empty, use current working directory
	if dstDir == "" {
		var err error
		dstDir, err = os.Getwd()
		if err != nil {
			cRedBold.Println("Failed to get current working dir:", err)
			return
		}
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

	// Create rcc file
	fmt.Print("Generating Qt resource file...")

	err = generator.CreateRccFile(cfg.Rcc, dstDir)
	if err != nil {
		cYellow.Println(err)
	} else {
		cGreen.Println("done")
	}
}
