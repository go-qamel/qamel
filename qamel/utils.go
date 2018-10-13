package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	fp "path/filepath"

	"github.com/fatih/color"
)

var (
	cBold      = color.New(color.FgCyan, color.Bold)
	cRed       = color.New(color.FgRed)
	cRedBold   = color.New(color.FgRed, color.Bold)
	cGreen     = color.New(color.FgGreen)
	cGreenBold = color.New(color.FgGreen, color.Bold)
	cBlueBold  = color.New(color.FgBlue, color.Bold)
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

// getPackageName gets the package name from specified file
func getPackageName(filePath string) (string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", fmt.Errorf("failed to get package name: %s", err)
	}

	if f.Name == nil {
		return "", fmt.Errorf("failed to get package name: no package name found")
	}

	return f.Name.Name, nil
}

// getPackageNameFromDir gets the package name from specified directory
func getPackageNameFromDir(dirPath string) (string, error) {
	dirItems, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to get package name: %s", err)
	}

	fileName := ""
	for _, item := range dirItems {
		if !item.IsDir() && fp.Ext(item.Name()) == ".go" {
			fileName = fp.Join(dirPath, item.Name())
			break
		}
	}

	if fileName == "" {
		return "", fmt.Errorf("failed to get package name: no go file exists in directory")
	}

	return getPackageName(fileName)
}
