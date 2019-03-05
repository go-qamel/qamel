package qamel

// #include <stdlib.h>
// #include "quickstyle.h"
import "C"
import "unsafe"

func SetQuickStyle(style string) {
    cStyle := C.CString(style)
    defer C.free(unsafe.Pointer(cStyle))
    C.SetQuickStyle(cStyle)
}

func SetQuickFallbackStyle(style string) {
    cStyle := C.CString(style)
    defer C.free(unsafe.Pointer(cStyle))
    C.SetQuickFallbackStyle(cStyle)
}

func AddQuickStylePath(style string) {
    cStyle := C.CString(style)
    defer C.free(unsafe.Pointer(cStyle))
    C.AddQuickStylePath(cStyle)
}
