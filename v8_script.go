package v8

/*
#include "v8_warp.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "runtime"

type Script struct {
	self    unsafe.Pointer
	context *Context
}

func (c *Context) CompileScript(code string) *Script {
	ccode := C.CString(code)
	self := C.V8_CompileScript(c.self, ccode)
	C.free(unsafe.Pointer(ccode))

	if self == nil {
		return nil
	}

	result := &Script{
		self:    self,
		context: c,
	}

	runtime.SetFinalizer(result, func(s *Script) {
		if traceDispose {
			println("v8.Script.Dispose()")
		}
		C.V8_DisposeScript(s.self)
	})

	return result
}

func (s Script) Run() *Value {
	if v := C.V8_RunScript(s.context.self, s.self); v != nil {
		return newValue(v)
	}
	return nil
}
