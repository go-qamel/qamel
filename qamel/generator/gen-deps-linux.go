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

		err := copyFileDir(srcPath, dstPath, nil)
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

	// The plugins and qml libraries (*.so) which copied before
	// must be checked as well
	outputDir := fp.Dir(outputPath)
	fp.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		fileExt := fp.Ext(info.Name())
		if strings.HasPrefix(fileExt, ".so") {
			filesToCheck = append(filesToCheck, path)
		}

		return nil
	})

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

func createLinuxScript(outputPath string) error {
	// Prepare content of the script
	scriptContent := "" +
		"#!/bin/bash\n" +
		"appname=`basename $0 | sed s,\\.sh$,,`\n\n" +
		"dirname=`dirname $0`\n" +
		"tmp=\"${dirname#?}\"\n\n" +
		"if [ \"${dirname%$tmp}\" != \"/\" ]; then\n" +
		"dirname=$PWD/$dirname\n" +
		"fi\n" +
		"export LD_LIBRARY_PATH=\"$dirname/libs\"\n" +
		"export QT_PLUGIN_PATH=\"$dirname/plugins\"\n" +
		"export QML_IMPORT_PATH=\"$dirname/qml\"\n" +
		"export QML2_IMPORT_PATH=\"$dirname/qml\"\n" +
		"$dirname/$appname \"$@\"\n"

	// If scripts already exists, remove it
	scriptPath := strings.TrimSuffix(outputPath, fp.Ext(outputPath))
	scriptPath += ".sh"

	err := os.Remove(scriptPath)
	if err != nil {
		return err
	}

	// Write script to file
	dstFile, err := os.OpenFile(scriptPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = dstFile.WriteString(scriptContent)
	if err != nil {
		return err
	}

	return dstFile.Sync()
}
