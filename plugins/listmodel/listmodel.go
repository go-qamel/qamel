package listmodel

// #include <stdlib.h>
// #include <stdint.h>
// #include <string.h>
// #include "listmodel.h"
import "C"
import "unsafe"

// RegisterQML registers QamelListModel as QML object
func RegisterQML(uri string, versionMajor int, versionMinor int, qmlName string) {
	cURI := C.CString(uri)
	cQmlName := C.CString(qmlName)
	cVersionMajor := C.int(int32(versionMajor))
	cVersionMinor := C.int(int32(versionMinor))
	defer func() {
		C.free(unsafe.Pointer(cURI))
		C.free(unsafe.Pointer(cQmlName))
	}()

	C.QamelListModel_RegisterQML(cURI, cVersionMajor, cVersionMinor, cQmlName)
}
