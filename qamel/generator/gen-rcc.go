package generator

import (
	"fmt"
	"os"
	"os/exec"
	fp "path/filepath"
)

// CreateRccFile creates rcc.cpp file from resource directory at `dstDir/res`
func CreateRccFile(rccPath string, dstDir string) error {
	// Check if resource directory is exist
	resDir := fp.Join(dstDir, "res")
	if !dirExists(resDir) {
		return fmt.Errorf("resource directory doesn't exist")
	}

	// Get list of file inside resource dir
	resFiles := []string{}
	fp.Walk(resDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		path, _ = fp.Rel(dstDir, path)
		resFiles = append(resFiles, path)
		return nil
	})

	if len(resFiles) == 0 {
		return fmt.Errorf("no resource available")
	}

	// Create temp qrc file
	qrcPath := fp.Join(dstDir, "qamel.qrc")
	qrcFile, err := os.Create(qrcPath)
	if err != nil {
		return err
	}

	defer func() {
		qrcFile.Close()
		os.Remove(qrcPath)
	}()

	qrcContent := fmt.Sprintln(`<!DOCTYPE RCC><RCC version="1.0">`)
	qrcContent += fmt.Sprintln(`<qresource>`)
	for _, resFile := range resFiles {
		qrcContent += fmt.Sprintf("<file>%s</file>\n", resFile)
	}
	qrcContent += fmt.Sprintln(`</qresource>`)
	qrcContent += fmt.Sprintln(`</RCC>`)

	_, err = qrcFile.WriteString(qrcContent)
	if err != nil {
		return err
	}
	qrcFile.Sync()

	// Run rcc
	dst := fp.Join(dstDir, "qamel-rcc.cpp")
	cmdRcc := exec.Command(rccPath, "-o", dst, qrcPath)
	return cmdRcc.Run()
}
