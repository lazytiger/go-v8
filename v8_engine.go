package v8

/*
#include "v8_wrap.h"
*/
import "C"
import "unsafe"
import "runtime"

var traceDispose = false

// Represents an isolated instance of the V8 engine.
// Objects from one engine must not be used in other engine.
type Engine struct {
	self       unsafe.Pointer
	_undefined *Value
	_null      *Value
	_true      *Value
	_false     *Value
}

func NewEngine() *Engine {
	self := C.V8_NewEngine()

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
		C.V8_DisposeEngine(i.self)
	})

	return result
}
