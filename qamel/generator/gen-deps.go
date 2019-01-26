package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	fp "path/filepath"
	"strings"

	"github.com/RadhiFadlillah/qamel/qamel/config"
)

// CopyDependencies copies dependecies files to output directory
func CopyDependencies(profile config.Profile, projectDir, outputPath string) error {
	// Get qmake variables from `qmake -query`
	cmdQmake := exec.Command(profile.Qmake, "-query")
	btOutput, err := cmdQmake.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n%s", err, btOutput)
	}

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

	// Copy QML
	outputDir := fp.Dir(outputPath)
	err = copyQmlDependencies(qmakeVars, profile, projectDir, outputDir)
	if err != nil {
		return err
	}

	// Copy OS specific libraries
	switch profile.OS {
	case "linux":
		return copyLinuxDependencies(qmakeVars, outputPath)
	case "windows":
		return copyWindowsDependencies(qmakeVars, profile, outputPath)
	}
	return nil
}

// copyLinuxDependencies copies dependecies files for Linux target
func copyLinuxDependencies(qmakeVars map[string]string, outputPath string) error {
	// Copy plugins
	err := copyLinuxPlugins(qmakeVars, fp.Dir(outputPath))
	if err != nil {
		return err
	}

	err = copyLinuxLibs(qmakeVars, outputPath)
	if err != nil {
		return err
	}

	return createLinuxScript(outputPath)
}

// copyWindowsDependencies copies dependecies files for Windows target
func copyWindowsDependencies(qmakeVars map[string]string, profile config.Profile, outputPath string) error {
	// Copy plugins
	err := copyWindowsPlugins(qmakeVars, fp.Dir(outputPath))
	if err != nil {
		return err
	}

	return copyWindowsLibs(qmakeVars, profile, outputPath)
}
