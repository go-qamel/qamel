package main

import (
	"bufio"
	"fmt"
	"os"
	fp "path/filepath"
	"strings"

	"github.com/RadhiFadlillah/qamel/qamel/generator"
	"github.com/spf13/cobra"
)

var cmdSetup = &cobra.Command{
	Use:   "setup",
	Short: "Set up the Qt's binding",
	Args:  cobra.NoArgs,
	Run:   setupHandler,
}

func setupHandler(cmd *cobra.Command, args []string) {
	// Define input reader
	reader := bufio.NewReader(os.Stdin)

	// Fetch path to Qt's directory and tools
	cBlueBold.Println("Thanks for using qamel, QML's binding for Go.")
	fmt.Println()
	fmt.Println("Please specify the *full path* to your Qt's tools directory.")
	fmt.Println("This might be different depending on your platform. " +
		"For example, in Linux with Qt 5.11.1, " +
		"the tools are located in $HOME/Qt5.11.1/5.11.1/gcc_64/bin/")
	fmt.Println()

	cCyanBold.Print("Qt tools dir : ")
	qtDir, _ := reader.ReadString('\n')
	qtDir = strings.TrimSpace(qtDir)
	if !dirExists(qtDir) {
		cRedBold.Println("The specified directory does not exist")
		return
	}

	// Make sure qmake, moc, and rcc is exists
	qmakePath := fp.Join(qtDir, "qmake")
	mocPath := fp.Join(qtDir, "moc")
	rccPath := fp.Join(qtDir, "rcc")

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
				return
			}
		}

		if !mocExists {
			cCyanBold.Print("Path to moc   : ")
			mocPath, _ = reader.ReadString('\n')
			mocPath = strings.TrimSpace(mocPath)
			if !fileExists(mocPath) {
				cRedBold.Println("The specified path does not exist")
				return
			}
		}

		if !rccExists {
			cCyanBold.Print("Path to rcc   : ")
			rccPath, _ = reader.ReadString('\n')
			rccPath = strings.TrimSpace(rccPath)
			if !fileExists(rccPath) {
				cRedBold.Println("The specified path does not exist")
				return
			}
		}
	}

	// Generating cgo code for binding
	fmt.Println()
	fmt.Print("Generating some code for binding...")

	cgoFlags, err := generator.CreateCgoFlags(qmakePath)
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to create cgo file:", err)
		return
	}

	err = generator.CreateCgoFile(qamelDir, cgoFlags)
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to create cgo file:", err)
		return
	}

	// Generating moc file for viewer
	err = generator.CreateMocFile(mocPath, fp.Join(qamelDir, "viewer.cpp"))
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to create moc file for viewer:", err)
		return
	}

	cGreen.Println("done")

	// Save config file
	fmt.Print("Saving config file...")

	err = saveConfig(config{
		Qmake: qmakePath,
		Moc:   mocPath,
		Rcc:   rccPath})
	if err != nil {
		fmt.Println()
		cRedBold.Println("Failed to save the config file:", err)
		return
	}

	cGreen.Println("done")

	// Setup finished
	fmt.Println()
	cBlueBold.Println("Setup finished.")
	cBlueBold.Println("Now you can get started on your own QML app.")
}
