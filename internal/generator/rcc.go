package generator

import (
	"fmt"
	"os"
	"os/exec"
	fp "path/filepath"

	"github.com/RadhiFadlillah/qamel/internal/config"
)

// ErrNoResourceDir is error that fired when resource directory doesn't exist
var ErrNoResourceDir = fmt.Errorf("resource directory doesn't exist")

// CreateRccFile creates rcc.cpp file from resource directory at `projectDir/res`
func CreateRccFile(profile config.Profile, projectDir string) error {
	// Create cgo file
	err := CreateCgoFile(profile, projectDir, "main")
	if err != nil {
		return err
	}

	// Check if resource directory is exist
	resDir := fp.Join(projectDir, "res")
	if !dirExists(resDir) {
		return ErrNoResourceDir
	}

	// Get list of file inside resource dir
	resFiles := []string{}
	fp.Walk(resDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		path, _ = fp.Rel(projectDir, path)
		resFiles = append(resFiles, path)
		return nil
	})

	if len(resFiles) == 0 {
		return fmt.Errorf("no resource available")
	}

	// Create temp qrc file
	qrcPath := fp.Join(projectDir, "qamel.qrc")
	defer os.Remove(qrcPath)

	qrcContent := fmt.Sprintln(`<!DOCTYPE RCC><RCC version="1.0">`)
	qrcContent += fmt.Sprintln(`<qresource>`)
	for _, resFile := range resFiles {
		qrcContent += fmt.Sprintf("<file>%s</file>\n", resFile)
	}
	qrcContent += fmt.Sprintln(`</qresource>`)
	qrcContent += fmt.Sprintln(`</RCC>`)

	err = saveToFile(qrcPath, qrcContent)
	if err != nil {
		return err
	}

	// Run rcc
	dst := fp.Join(projectDir, "qamel-rcc.cpp")
	cmdRcc := exec.Command(profile.Rcc, "-o", dst, qrcPath)
	btOutput, err := cmdRcc.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n%s", err, btOutput)
	}

	return nil
}
