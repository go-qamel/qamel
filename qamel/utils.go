package main

import (
	"os"

	"github.com/fatih/color"
)

var (
	cCyan      = color.New(color.FgCyan)
	cCyanBold  = color.New(color.FgCyan, color.Bold)
	cRed       = color.New(color.FgRed)
	cRedBold   = color.New(color.FgRed, color.Bold)
	cGreen     = color.New(color.FgGreen)
	cGreenBold = color.New(color.FgGreen, color.Bold)
	cBlueBold  = color.New(color.FgBlue, color.Bold)
	cYellow    = color.New(color.FgYellow)
)

// fileExists checks if the file in specified path is exists
func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	return !os.IsNotExist(err) && !info.IsDir()
}

// dirExists checks if the directory in specified path is exists
func dirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	return !os.IsNotExist(err) && info.IsDir()
}
