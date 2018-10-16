package qamel

import (
	"sync"
	"unsafe"
)

var (
	mutex     = sync.Mutex{}
	mapObject = map[unsafe.Pointer]interface{}{}
)

// QmlObject is the base of QML object
type QmlObject struct {
	Ptr unsafe.Pointer
}

// RegisterObject registers the specified pointer to specified object
func RegisterObject(ptr unsafe.Pointer, obj interface{}) {
	if ptr == nil || obj == nil {
		return
	}

	mutex.Lock()
	mapObject[ptr] = obj
	mutex.Unlock()
}

// BorrowObject fetch object for the specified pointer
func BorrowObject(ptr unsafe.Pointer) interface{} {
	if ptr == nil {
		return nil
	}

	mutex.Lock()
	return mapObject[ptr]
}

// ReturnObject returns pointer and lock map again
func ReturnObject(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}

	mutex.Unlock()
}

// DeleteObject remove object for the specified pointer
func DeleteObject(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}

	mutex.Lock()
	delete(mapObject, ptr)
	mutex.Unlock()
}
