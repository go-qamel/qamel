package application

// #include <stdlib.h>
// #include <stdbool.h>
// #include "application.h"
import "C"
import "unsafe"

// Application is the main app which wraps QGuiApplication
type Application struct {
	ptr unsafe.Pointer
}

// NewApplication initializes the window system and constructs
// an QGuiApplication object with argc command line arguments in argv.
func NewApplication(argc int, argv []string) Application {
	argvC := sliceToChars(argv)
	defer C.free(unsafe.Pointer(argvC))

	ptr := C.App_NewApplication(C.int(argc), argvC)
	return Application{ptr: ptr}
}

// SetAttribute sets the Application's attribute if on is true;
// otherwise clears the attribute
func SetAttribute(attribute Attribute, on bool) {
	C.App_SetAttribute(C.int(attribute), C.bool(on))
}

// SetFont changes the default application font
func (app Application) SetFont(fontFamily string, pointSize int, weight FontWeight, italic bool) {
	if app.ptr == nil {
		return
	}

	cFontFamily := C.CString(fontFamily)
	defer C.free(unsafe.Pointer(cFontFamily))
	C.App_SetFont(app.ptr, cFontFamily, C.int(pointSize), C.int(weight), C.bool(italic))
}

// SetQuitOnLastWindowClosed set whether the application implicitly quits
// when the last window is closed. The default is true.
func (app Application) SetQuitOnLastWindowClosed(quit bool) {
	if app.ptr == nil {
		return
	}

	C.App_SetQuitOnLastWindowClosed(app.ptr, C.bool(quit))
}

// SetApplicationDisplayName sets the user visible name of this application
func (app Application) SetApplicationDisplayName(name string) {
	if app.ptr == nil {
		return
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.App_SetApplicationDisplayName(app.ptr, cName)
}

// SetWindowIcon sets the default window icon from the specified fileName.
// The fileName could be a relative path (../icon), Qt resource path (qrc://icon),
// or absolute file path (file://path/to/icon)
func (app Application) SetWindowIcon(fileName string) {
	if app.ptr == nil {
		return
	}

	cFileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))
	C.App_SetWindowIcon(app.ptr, cFileName)
}

// SetApplicationName sets the name of this application
func (app Application) SetApplicationName(name string) {
	if app.ptr == nil {
		return
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.App_SetApplicationName(app.ptr, cName)
}

// SetApplicationVersion sets the version of this application
func (app Application) SetApplicationVersion(version string) {
	if app.ptr == nil {
		return
	}

	cVersion := C.CString(version)
	defer C.free(unsafe.Pointer(cVersion))
	C.App_SetApplicationVersion(app.ptr, cVersion)
}

// SetOrganizationName sets the name of the organization that
// write this application
func (app Application) SetOrganizationName(name string) {
	if app.ptr == nil {
		return
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.App_SetOrganizationName(app.ptr, cName)
}

// SetOrganizationDomain sets the internet domain of the
// organization that write this application
func (app Application) SetOrganizationDomain(domain string) {
	if app.ptr == nil {
		return
	}

	cDomain := C.CString(domain)
	defer C.free(unsafe.Pointer(cDomain))
	C.App_SetOrganizationDomain(app.ptr, cDomain)
}

// Exec enters the main event loop and waits until exit() is called, and
// then returns the value that was set to exit() (which is 0 if exit() is
// called via quit()).
func (app Application) Exec() int {
	return int(int32(C.App_Exec()))
}

func sliceToChars(src []string) **C.char {
	cArray := C.malloc(C.size_t(len(src)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	a := (*[1<<30 - 1]*C.char)(cArray)

	for i, str := range src {
		a[i] = C.CString(str)
	}

	return (**C.char)(cArray)
}
