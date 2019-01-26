package generator

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	fp "path/filepath"
	"strings"
)

func copyLinuxPlugins(qmakeVars map[string]string, outputDir string) error {
	qtPluginsDir := qmakeVars["QT_INSTALL_PLUGINS"]
	dstPluginsDir := fp.Join(outputDir, "plugins")
	plugins := []string{
		"platforms/libqxcb.so",
		"platforminputcontexts",
		"imageformats",
		"xcbglintegrations",
		"iconengines/libqsvgicon.so"}

	for _, plugin := range plugins {
		srcPath := fp.Join(qtPluginsDir, plugin)
		dstPath := fp.Join(dstPluginsDir, plugin)

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

func copyLinuxLibs(qmakeVars map[string]string, outputPath string) error {
	// Get list of files to check.
	// The first one is the output binary file
	filesToCheck := []string{outputPath}

	// The plugins library (*.so) must be checked as well
	outputDir := fp.Dir(outputPath)
	dstPluginsDir := fp.Join(outputDir, "plugins")
	if dirExists(dstPluginsDir) {
		fp.Walk(dstPluginsDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			fileExt := fp.Ext(info.Name())
			if strings.HasPrefix(fileExt, ".so") {
				filesToCheck = append(filesToCheck, path)
			}

			return nil
		})
	}

	// Get list of dependency of the files
	qtLibsDir := qmakeVars["QT_INSTALL_LIBS"]
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

		// Fetch dependencies using ldd
		cmdLdd := exec.Command("ldd", fileName)
		lddResult, err := cmdLdd.CombinedOutput()
		if err != nil {
			continue
		}

		// Parse ldd results
		// It will look like this: libname.so => /path/to/lib (memorty address)
		buffer := bytes.NewBuffer(lddResult)
		scanner := bufio.NewScanner(buffer)

		for scanner.Scan() {
			lib := strings.TrimSpace(scanner.Text())
			libName := strings.SplitN(lib, " ", 2)[0]
			libName = fp.Base(libName)
			if libName == "" {
				continue
			}

			libPath := fp.Join(qtLibsDir, libName)
			if !fileExists(libPath) {
				continue
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
	dstLibsDir := fp.Join(outputDir, "libs")
	for libName, libPath := range mapDependencies {
		err = copyFile(libPath, fp.Join(dstLibsDir, libName))
		if err != nil {
			return err
		}
	}

	return nil
}
