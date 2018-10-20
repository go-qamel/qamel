package qamel

// #include <stdlib.h>
// #include <stdbool.h>
// #include "viewer.h"
import "C"
import "unsafe"

// Viewer is the QML viewer which wraps QQuickView
type Viewer struct {
	ptr unsafe.Pointer
}

// NewViewer constructs a QQuickView.
func NewViewer() Viewer {
	ptr := C.Viewer_NewViewer()
	return Viewer{ptr: ptr}
}

// NewViewerWithSource constructs a QQuickView with the given QML source.
func NewViewerWithSource(source string) Viewer {
	view := NewViewer()
	view.SetSource(source)
	return view
}

// SetSource sets the source to the url, loads the QML component and instantiates it.
// The source could be a Qt resource path (qrc://icon) or a file path (file://path/to/icon).
// However, it must be a valid path.
func (view Viewer) SetSource(url string) {
	if view.ptr == nil {
		return
	}

	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))
	C.Viewer_SetSource(view.ptr, cURL)
}

// SetResizeMode sets whether the view should resize the window contents.
// If this property is set to SizeViewToRootObject (the default), the view resizes
// to the size of the root item in the QML. If this property is set to
// SizeRootObjectToView, the view will automatically resize the root item to the
// size of the view.
func (view Viewer) SetResizeMode(resizeMode ResizeMode) {
	if view.ptr == nil {
		return
	}

	C.Viewer_SetResizeMode(view.ptr, C.int(resizeMode))
}

// SetFlags sets the flags of the window. The window flags control the window's appearance
// in the windowing system, whether it's a dialog, popup, or a regular window, and whether it
// should have a title bar, etc. The actual window flags might differ from the flags set with
// setFlags() if the requested flags could not be fulfilled.
func (view Viewer) SetFlags(flags WindowFlags) {
	if view.ptr == nil {
		return
	}

	C.Viewer_SetFlags(view.ptr, C.int(flags))
}

// SetHeight sets the height of the window.
func (view Viewer) SetHeight(height int) {
	if view.ptr == nil {
		return
	}

	C.Viewer_SetHeight(view.ptr, C.int(int32(height)))
}

// SetWidth sets the width of the window.
func (view Viewer) SetWidth(width int) {
	if view.ptr == nil {
		return
	}

	C.Viewer_SetWidth(view.ptr, C.int(int32(width)))
}

// SetMaximumHeight sets the maximum height of the window.
func (view Viewer) SetMaximumHeight(height int) {
	if view.ptr == nil {
		return
	}

	C.Viewer_SetMaximumHeight(view.ptr, C.int(int32(height)))
}

// SetMaximumWidth sets the maximum width of the window.
func (view Viewer) SetMaximumWidth(width int) {
	if view.ptr == nil {
		return
	}

	C.Viewer_SetMaximumWidth(view.ptr, C.int(int32(width)))
}

// SetMinimumHeight sets the minimum height of the window.
func (view Viewer) SetMinimumHeight(height int) {
	if view.ptr == nil {
		return
	}

	C.Viewer_SetMinimumHeight(view.ptr, C.int(int32(height)))
}

// SetMinimumWidth sets the minimum width of the window.
func (view Viewer) SetMinimumWidth(width int) {
	if view.ptr == nil {
		return
	}

	C.Viewer_SetMinimumWidth(view.ptr, C.int(int32(width)))
}

// SetOpacity sets the opacity of the window in the windowing system. If the windowing system supports
// window opacity, this can be used to fade the window in and out, or to make it semitransparent.
// A value of 1.0 or above is treated as fully opaque, whereas a value of 0.0 or below is treated as
// fully transparent. Values inbetween represent varying levels of translucency between the two extremes.
// The default value is 1.0.
func (view Viewer) SetOpacity(opacity float64) {
	if view.ptr == nil {
		return
	}

	C.Viewer_SetOpacity(view.ptr, C.double(opacity))
}

// SetTitle sets the window's title in the windowing system. The window title might appear in the title
// area of the window decorations, depending on the windowing system and the window flags. It might also
// be used by the windowing system to identify the window in other contexts, such as in the task switcher.
func (view Viewer) SetTitle(title string) {
	if view.ptr == nil {
		return
	}

	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.Viewer_SetTitle(view.ptr, cTitle)
}

// SetVisible sets whether the window is visible or not. This property controls the visibility of the
// window in the windowing system. By default, the window is not visible, you must call setVisible(true),
// or show() or similar to make it visible.
func (view Viewer) SetVisible(visible bool) {
	if view.ptr == nil {
		return
	}

	C.Viewer_SetVisible(view.ptr, C.bool(visible))
}

// SetPosition sets the position of the window on the desktop to x, y.
func (view Viewer) SetPosition(x int, y int) {
	if view.ptr == nil {
		return
	}

	C.Viewer_SetPosition(view.ptr, C.int(int32(x)), C.int(int32(y)))
}

// SetIcon sets the window's icon in the windowing system. The window icon might be used by the windowing
// system for example to decorate the window, and/or in the task switcher.
// Note: On macOS, the window title bar icon is meant for windows representing documents, and will only
// show up if a file path is also set.
func (view Viewer) SetIcon(fileName string) {
	if view.ptr == nil {
		return
	}

	cFileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))
	C.Viewer_SetIcon(view.ptr, cFileName)
}

// Show shows the window. This is equivalent to calling showFullScreen(), showMaximized(), or
// showNormal(), depending on the platform's default behavior for the window type and flags.
func (view Viewer) Show() {
	if view.ptr == nil {
		return
	}

	C.Viewer_Show(view.ptr)
}

// ShowMaximized shows the window as maximized.
// Equivalent to calling setWindowStates(WindowMaximized) and then setVisible(true).
func (view Viewer) ShowMaximized() {
	if view.ptr == nil {
		return
	}

	C.Viewer_ShowMaximized(view.ptr)
}

// ShowMinimized shows the window as minimized.
// Equivalent to calling setWindowStates(WindowMinimized) and then setVisible(true).
func (view Viewer) ShowMinimized() {
	if view.ptr == nil {
		return
	}

	C.Viewer_ShowMinimized(view.ptr)
}

// ShowFullScreen shows the window as fullscreen.
// Equivalent to calling setWindowStates(WindowFullScreen) and then setVisible(true).
func (view Viewer) ShowFullScreen() {
	if view.ptr == nil {
		return
	}

	C.Viewer_ShowFullScreen(view.ptr)
}

// ShowNormal shows the window as normal, i.e. neither maximized, minimized, nor fullscreen.
// Equivalent to calling setWindowStates(WindowNoState) and then setVisible(true).
func (view Viewer) ShowNormal() {
	if view.ptr == nil {
		return
	}

	C.Viewer_ShowNormal(view.ptr)
}

// SetWindowStates sets the screen-occupation state of the window. The window state represents whether
// the window appears in the windowing system as maximized, minimized and/or fullscreen. The window can
// be in a combination of several states. For example, if the window is both minimized and maximized,
// the window will appear minimized, but clicking on the task bar entry will restore it to the
// maximized state.
func (view Viewer) SetWindowStates(state WindowStates) {
	if view.ptr == nil {
		return
	}

	C.Viewer_SetWindowStates(view.ptr, C.int(state))
}

// ClearComponentCache clears the engine's internal component cache. This function causes the property
// metadata of all components previously loaded by the engine to be destroyed. All previously loaded
// components and the property bindings for all extant objects created from those components will cease
// to function. This function returns the engine to a state where it does not contain any loaded
// component data. This may be useful in order to reload a smaller subset of the previous component set,
// or to load a new version of a previously loaded component. Once the component cache has been cleared,
// components must be loaded before any new objects can be created.
func (view Viewer) ClearComponentCache() {
	if view.ptr == nil {
		return
	}

	C.Viewer_ClearComponentCache(view.ptr)
}

// Reload reloads the active QML view.
func (view Viewer) Reload() {
	if view.ptr == nil {
		return
	}

	C.Viewer_Reload(view.ptr)
}
