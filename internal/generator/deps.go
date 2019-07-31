package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	fp "path/filepath"
	"strings"

	"github.com/RadhiFadlillah/qamel/internal/config"
)

// CopyDependencies copies dependencies to output directory
func CopyDependencies(profile config.Profile, projectDir, outputPath string) error {
	switch profile.OS {
	case "linux":
		return copyLinuxDependencies(profile, projectDir, outputPath)
	case "windows":
		return copyWindowsDependencies(profile, projectDir, outputPath)
	default:
		return nil
	}
}

// copyLinuxDependencies copies dependencies for Linux
func copyLinuxDependencies(profile config.Profile, projectDir, outputPath string) error {
	// Get qmake variables from `qmake -query`
	qmakeVars, err := getQmakeVars(profile.Qmake)
	if err != nil {
		return err
	}

	// Get dirs
	outputDir := fp.Dir(outputPath)
	qtQmlDir := qmakeVars["QT_INSTALL_QML"]
	qtLibsDir := qmakeVars["QT_INSTALL_LIBS"]
	qtPluginsDir := qmakeVars["QT_INSTALL_PLUGINS"]

	// Copy QML
	err = copyQmlDependencies(qtQmlDir, profile, projectDir, outputDir)
	if err != nil {
		return err
	}

	// Copy plugins
	err = copyLinuxPlugins(qtPluginsDir, outputDir)
	if err != nil {
		return err
	}

	// Copy libs
	err = copyLinuxLibs(qtLibsDir, outputPath)
	if err != nil {
		return err
	}

	// Create script
	return createLinuxScript(outputPath)
}

// copyWindowsDependencies copies dependencies for Windows
func copyWindowsDependencies(profile config.Profile, projectDir, outputPath string) error {
	// Get qmake variables from `qmake -query`
	qmakeVars, err := getQmakeVars(profile.Qmake)
	if err != nil {
		return err
	}

	// Get dirs
	outputDir := fp.Dir(outputPath)
	qtQmlDir := qmakeVars["QT_INSTALL_QML"]
	qtPluginsDir := qmakeVars["QT_INSTALL_PLUGINS"]

	// Copy QML
	err = copyQmlDependencies(qtQmlDir, profile, projectDir, outputDir)
	if err != nil {
		return err
	}

	// Copy plugins
	err = copyWindowsPlugins(qtPluginsDir, outputDir)
	if err != nil {
		return err
	}

	// Copy libs
	return copyWindowsLibs(profile, outputPath)
}

// getQmakeVars get the qmake properties by running `qmake -query`
func getQmakeVars(qmakePath string) (map[string]string, error) {
	// Run qmake
	cmdQmake := exec.Command(qmakePath, "-query")
	btOutput, err := cmdQmake.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%v\n%s", err, btOutput)
	}

	// Parse output
	buffer := bytes.NewBuffer(btOutput)
	scanner := bufio.NewScanner(buffer)
	qmakeVars := map[string]string{}
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ":", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		qmakeVars[name] = value
	}

	return qmakeVars, nil
}
