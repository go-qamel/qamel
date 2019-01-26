package generator

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	fp "path/filepath"
	"runtime"
	"strings"

	"github.com/RadhiFadlillah/qamel/qamel/config"
)

func copyWindowsPlugins(qmakeVars map[string]string, outputDir string) error {
	qtPluginsDir := qmakeVars["QT_INSTALL_PLUGINS"]
	plugins := []string{
		"platforms/qwindows.dll",
		"imageformats",
		"iconengines/qsvgicon.dll"}

	for _, plugin := range plugins {
		srcPath := fp.Join(qtPluginsDir, plugin)
		dstPath := fp.Join(outputDir, plugin)

		if !fileDirExists(srcPath) {
			continue
		}

		err := copyFileDir(srcPath, dstPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyWindowsLibs(qmakeVars map[string]string, profile config.Profile, outputPath string) error {
	// Get list of files to check.
	// The first one is the output binary file
	filesToCheck := []string{outputPath}

	// The plugins library (*.dll) must be checked as well
	outputDir := fp.Dir(outputPath)
	if dirExists(outputDir) {
		fp.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if fp.Ext(info.Name()) == ".dll" {
				filesToCheck = append(filesToCheck, path)
			}

			return nil
		})
	}

	// Get directory of Qt and Gcc libs
	qtLibsDir := fp.Dir(profile.Qmake)
	gccLibsDir := fp.Dir(profile.Gcc)
	if runtime.GOOS != "windows" {
		// If it's not Windows, it must be using MXE,
		// so find GCC using relative path from qmake location
		gccLibsDir = fp.Join(qtLibsDir, "..", "..", "bin")
	}

	// Get list of dependency of the files
	mapDependencies := map[string]string{}
	filesAlreadyChecked := map[string]struct{}{}

	for {
		// Keep checking until there are no file left
		if len(filesToCheck) == 0 {
			break
		}
		fileName := filesToCheck[0]

		// If this file has been checked before, skip
		if _, checked := filesAlreadyChecked[fileName]; checked {
			filesToCheck = filesToCheck[1:]
			continue
		}

		// Fetch dependencies using objdump
		cmdObjdump := exec.Command(profile.Objdump, "-p", fileName)
		objdumpResult, err := cmdObjdump.CombinedOutput()
		if err != nil {
			continue
		}

		// Parse objdump results
		buffer := bytes.NewBuffer(objdumpResult)
		scanner := bufio.NewScanner(buffer)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(line, "DLL Name:") {
				continue
			}

			libName := strings.TrimPrefix(line, "DLL Name:")
			libName = strings.TrimSpace(libName)
			if libName == "" {
				continue
			}

			libPath := fp.Join(qtLibsDir, libName)
			if !fileExists(libPath) {
				libPath = fp.Join(gccLibsDir, libName)
				if !fileExists(libPath) {
					continue
				}
			}

			filesToCheck = append(filesToCheck, libPath)
			mapDependencies[libName] = libPath
		}

		// Save this files as been already checked, then move to next
		filesAlreadyChecked[fileName] = struct{}{}
		filesToCheck = filesToCheck[1:]
	}

	// Copy all dependency libs to output dir
	var err error
	for libName, libPath := range mapDependencies {
		err = copyFile(libPath, fp.Join(outputDir, libName))
		if err != nil {
			return err
		}
	}

	return nil
}
