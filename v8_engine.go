package v8

/*
#include "v8_warp.h"
*/
import "C"
import "unsafe"
import "runtime"

var traceDispose = false

var Default = NewEngine()

// Represents an isolated instance of the V8 engine.
// Objects from one engine must not be used in other engine.
type Engine struct {
	self unsafe.Pointer
}

func NewEngine() *Engine {
	self := C.V8_NewIsolate()

	if self == nil {
		return nil
	}

	result := &Engine{
		self: self,
	}

	runtime.SetFinalizer(result, func(i *Engine) {
		if traceDispose {
			println("v8.Engine.Dispose()")
		}
		C.V8_DisposeIsolate(i.self)
	})

	return result
}
