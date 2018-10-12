package qamel

// #include <stdlib.h>
// #include <stdbool.h>
// #include "engine.h"
import "C"
import "unsafe"

// Engine is the wrapper for QQMLApplicationEngine
type Engine struct {
	ptr unsafe.Pointer
}

// NewEngine creates a new QQmlApplicationEngine with the given parent.
// You will have to call load() later in order to load a QML file.
func NewEngine() Engine {
	ptr := C.Engine_NewEngine()
	return Engine{ptr: ptr}
}

// NewEngineWithSource constructs a QQmlApplicationEngine with the given QML source.
func NewEngineWithSource(source string) Engine {
	cSource := C.CString(source)
	defer C.free(unsafe.Pointer(cSource))

	ptr := C.Engine_NewEngineWithSource(cSource)
	return Engine{ptr: ptr}
}

// Load loads the root QML file located at url. The object tree defined by the file is
// created immediately for local file urls.
func (engine Engine) Load(url string) {
	if engine.ptr == nil {
		return
	}

	cURL := C.CString(url)
	defer C.free(unsafe.Pointer(cURL))
	C.Engine_Load(engine.ptr, cURL)
}

// ClearComponentCache clears the engine's internal component cache. This function causes the property
// metadata of all components previously loaded by the engine to be destroyed. All previously loaded
// components and the property bindings for all extant objects created from those components will cease
// to function. This function returns the engine to a state where it does not contain any loaded
// component data. This may be useful in order to reload a smaller subset of the previous component set,
// or to load a new version of a previously loaded component. Once the component cache has been cleared,
// components must be loaded before any new objects can be created.
func (engine Engine) ClearComponentCache() {
	if engine.ptr == nil {
		return
	}

	C.Engine_ClearComponentCache(engine.ptr)
}
