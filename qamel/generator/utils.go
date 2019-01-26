package generator

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	fp "path/filepath"
	"strings"
	"unicode"
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

// fileDirExists checks if the file or dir in specified path is exists
func fileDirExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// dirEmpty checks if the directory in specified path is EMPTY
func dirEmpty(dirPath string) bool {
	f, err := os.Open(dirPath)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}

	return false
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

// getSubDirs fetch all sub directories, including the root dir
func getSubDirs(rootDir string) ([]string, error) {
	subDirs := []string{}
	err := fp.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		if strings.HasPrefix(info.Name(), ".") {
			return fp.SkipDir
		}

		subDirs = append(subDirs, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return subDirs, nil
}

// upperChar change char in pos to uppercase.
func upperChar(str string, pos int) string {
	if pos < 0 || pos >= len(str) {
		return str
	}

	tmp := []byte(str)
	upper := unicode.ToUpper(rune(tmp[pos]))
	tmp[pos] = byte(upper)
	return string(tmp)
}

// saveToFile saves a string content to file
func saveToFile(dstPath string, content string) error {
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = dstFile.WriteString(content)
	if err != nil {
		return err
	}

	return dstFile.Sync()
}

// copyFile copies file from srcPath to dstPath
func copyFile(srcPath, dstPath string) error {
	os.MkdirAll(fp.Dir(dstPath), os.ModePerm)

	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// copyDir copies dir from srcPath to dstPath
func copyDir(srcPath, dstPath string, isSkipped func(string, os.FileInfo) bool) error {
	os.MkdirAll(fp.Dir(dstPath), os.ModePerm)

	fds, err := ioutil.ReadDir(srcPath)
	if err != nil {
		return err
	}

	for _, fd := range fds {
		srcfp := fp.Join(srcPath, fd.Name())
		dstfp := fp.Join(dstPath, fd.Name())

		if isSkipped != nil && isSkipped(srcfp, fd) {
			continue
		}

		if fd.IsDir() {
			err = copyDir(srcfp, dstfp, isSkipped)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcfp, dstfp)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFileDir is wrapper for both copyFile and copyDir.
// Useful when you don't care whether src is file or directory
func copyFileDir(srcPath, dstPath string, isSkipped func(string, os.FileInfo) bool) error {
	stat, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return copyDir(srcPath, dstPath, isSkipped)
	}

	return copyFile(srcPath, dstPath)
}
