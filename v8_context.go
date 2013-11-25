package v8

/*
#include "v8_wrap.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "runtime"

//import "reflect"

// A sandboxed execution context with its own set of built-in objects
// and functions.
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
			println("v8.Context.Dispose()", c.self)
		}
		C.V8_DisposeContext(c.self)
	})

	return result
}

//export try_catch_callback
func try_catch_callback(callback unsafe.Pointer) {
	(*(*func())(callback))()
}

func (c *Context) ThrowException(err string) {
	Default.Compile([]byte(`throw "`+err+`"`), nil, nil).Run(c)
	//
	// TODO: use ThrowException() will make function template GetFunction() return NULLv8_all_test.go, why?
	//
	//errPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&err)).Data)
	//C.V8_Context_ThrowException(c.self, (*C.char)(errPtr), C.int(len(err)))
}

func (c *Context) TryCatch(simple bool, callback func()) string {
	isSimple := 0
	if simple {
		isSimple = 1
	}
	creport := C.V8_Context_TryCatch(c.self, unsafe.Pointer(&callback), C.int(isSimple))
	if creport == nil {
		return ""
	}
	report := C.GoString(creport)
	C.free(unsafe.Pointer(creport))
	return report
}
