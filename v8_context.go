package v8

/*
#include "v8_warp.h"
*/
import "C"
import "unsafe"
import "runtime"

type Context struct {
	self   unsafe.Pointer
	engine *Engine
}

func (e *Engine) NewContext() *Context {
	self := C.V8_NewContext(e.self)

	if self == nil {
		return nil
	}

	result := &Context{
		self:   self,
		engine: e,
	}

	runtime.SetFinalizer(result, func(c *Context) {
		if traceDispose {
			println("v8.Context.Dispose()")
		}
		C.V8_DisposeContext(c.self)
	})

	return result
}
