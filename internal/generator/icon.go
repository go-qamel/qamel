package generator

import (
	"fmt"
	"os"
	"os/exec"
	fp "path/filepath"

	"github.com/go-qamel/qamel/internal/config"
)

// ErrNoIcon is error that fired when app icon doesn't exist
var ErrNoIcon = fmt.Errorf("icon file doesn't exist")

// CreateSysoFile creates syso file from ICO file in the specified source.
func CreateSysoFile(profile config.Profile, projectDir string) error {
	// Check if icon is exist
	iconPath := fp.Join(projectDir, "icon.ico")
	if !fileExists(iconPath) {
		return ErrNoIcon
	}

	// Create temporary rc file
	rcFilePath := fp.Join(projectDir, "qamel-icon.rc")
	defer os.Remove(rcFilePath)

	err := saveToFile(rcFilePath, `IDI_ICON1 ICON DISCARDABLE "icon.ico"`)
	if err != nil {
		return err
	}

	// Create syso file
	sysoFilePath := fp.Join(projectDir, "qamel-icon.syso")
	cmdWindres := exec.Command(profile.Windres, "-i", rcFilePath, "-o", sysoFilePath)
	if btOutput, err := cmdWindres.CombinedOutput(); err != nil {
		return fmt.Errorf("%v: %s", err, btOutput)
	}

	return nil
}
