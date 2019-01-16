package generator

import (
	"fmt"
	"os"
	"os/exec"
	fp "path/filepath"

	"github.com/RadhiFadlillah/qamel/qamel/config"
)

// ErrNoIcon is error that fired when app icon doesn't exist
var ErrNoIcon = fmt.Errorf("icon file doesn't exist")

// CreateSysoFile creates syso file from ICO file in the specified source.
func CreateSysoFile(profile config.Profile, dstDir string) error {
	// Check if icon is exist
	iconPath := fp.Join(dstDir, "icon.ico")
	if !fileExists(iconPath) {
		return ErrNoIcon
	}

	// Create temporary rc file
	rcFilePath := fp.Join(dstDir, "qamel-icon.rc")
	defer os.Remove(rcFilePath)

	err := saveToFile(rcFilePath, `IDI_ICON1 ICON DISCARDABLE "icon.ico"`)
	if err != nil {
		return err
	}

	// Create syso file
	sysoFilePath := fp.Join(dstDir, "qamel-icon.syso")
	cmdWindres := exec.Command(profile.Windres, "-i", rcFilePath, "-o", sysoFilePath)
	return cmdWindres.Run()
}
