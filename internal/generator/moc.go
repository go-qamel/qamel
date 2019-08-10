package generator

import (
	"fmt"
	"os/exec"
	fp "path/filepath"
	"strings"
)

// CreateMocFile creates moc file from the specified source.
// If destination is not specified, the moc file will be saved in
// file "moc-" + source.name + ".h"
func CreateMocFile(mocPath string, src string) error {
	// Make sure source is exist
	if !fileExists(src) {
		return fmt.Errorf("source file doesn't exists")
	}

	// Create destination name
	dst := "moc-" + fp.Base(src)
	dst = strings.TrimSuffix(dst, fp.Ext(dst)) + ".h"
	dst = fp.Join(fp.Dir(src), dst)

	// Run moc
	cmdMoc := exec.Command(mocPath, "-o", dst, src)
	if btOutput, err := cmdMoc.CombinedOutput(); err != nil {
		return fmt.Errorf("%v: %s", err, btOutput)
	}

	return nil
}
