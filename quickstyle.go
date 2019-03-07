package qamel

// #include <stdlib.h>
// #include "quickstyle.h"
import "C"
import "unsafe"

// SetQuickStyle sets the application style to style.
// Note that the style must be configured before loading QML that
// imports Qt Quick Controls 2. It is not possible to change
// the style after the QML types have been registered.
func SetQuickStyle(style string) {
	cStyle := C.CString(style)
	defer C.free(unsafe.Pointer(cStyle))
	C.SetQuickStyle(cStyle)
}

// SetQuickStyleFallback sets the application fallback style
// to style. Note that the fallback style must be the name of
// one of the built-in Qt Quick Controls 2 styles, e.g. "Material".
func SetQuickStyleFallback(style string) {
	cStyle := C.CString(style)
	defer C.free(unsafe.Pointer(cStyle))
	C.SetQuickStyleFallback(cStyle)
}

// AddQuickStylePath adds path as a directory where Qt Quick
// Controls 2 searches for available styles. The path may be
// any local filesystem directory or Qt Resource directory.
func AddQuickStylePath(path string) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.AddQuickStylePath(cPath)
}
