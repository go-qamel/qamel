package qamel

// #include <stdint.h>
// #include <stdlib.h>
// #include <string.h>
// #include <stdbool.h>
// #include "application.h"
import "C"
import (
	"runtime"
	"strings"
	"unsafe"
)

func init() {
	runtime.LockOSThread()
}

// Application is the main app which wraps QGuiApplication
type Application struct {
	ptr unsafe.Pointer
}

// NewApplication initializes the window system and constructs
// an QGuiApplication object with argc command line arguments in argv.
func NewApplication(argc int, argv []string) *Application {
	argvC := C.CString(strings.Join(argv, "|"))
	defer C.free(unsafe.Pointer(argvC))

	ptr := C.App_NewApplication(C.int(int32(argc)), argvC)
	return &Application{ptr: ptr}
}

// SetAttribute sets the Application's attribute if on is true;
// otherwise clears the attribute
func SetAttribute(attribute Attribute, on bool) {
	C.App_SetAttribute(C.longlong(attribute), C.bool(on))
}

// SetFont changes the default application font
func (app Application) SetFont(fontFamily string, pointSize int, weight FontWeight, italic bool) {
	cFontFamily := C.CString(fontFamily)
	defer C.free(unsafe.Pointer(cFontFamily))
	C.App_SetFont(cFontFamily, C.int(int32(pointSize)), C.int(weight), C.bool(italic))
}

// SetQuitOnLastWindowClosed set whether the application implicitly quits
// when the last window is closed. The default is true.
func (app Application) SetQuitOnLastWindowClosed(quit bool) {
	C.App_SetQuitOnLastWindowClosed(C.bool(quit))
}

// SetApplicationDisplayName sets the user visible name of this application
func (app Application) SetApplicationDisplayName(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.App_SetApplicationDisplayName(cName)
}

// SetWindowIcon sets the default window icon from the specified fileName.
// The fileName could be a relative path (../icon), Qt resource path (qrc://icon),
// or absolute file path (file://path/to/icon)
func (app Application) SetWindowIcon(fileName string) {
	cFileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))
	C.App_SetWindowIcon(cFileName)
}

// SetApplicationName sets the name of this application
func (app Application) SetApplicationName(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.App_SetApplicationName(cName)
}

// SetApplicationVersion sets the version of this application
func (app Application) SetApplicationVersion(version string) {
	cVersion := C.CString(version)
	defer C.free(unsafe.Pointer(cVersion))
	C.App_SetApplicationVersion(cVersion)
}

// SetOrganizationName sets the name of the organization that
// write this application
func (app Application) SetOrganizationName(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.App_SetOrganizationName(cName)
}

// SetOrganizationDomain sets the internet domain of the
// organization that write this application
func (app Application) SetOrganizationDomain(domain string) {
	cDomain := C.CString(domain)
	defer C.free(unsafe.Pointer(cDomain))
	C.App_SetOrganizationDomain(cDomain)
}

// Exec enters the main event loop and waits until exit() is called, and
// then returns the value that was set to exit() (which is 0 if exit() is
// called via quit()).
func (app Application) Exec() int {
	return int(int32(C.App_Exec()))
}
