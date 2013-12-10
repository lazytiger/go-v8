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
	embedable
	self   unsafe.Pointer
	engine *Engine
}

type ContextScope struct {
	context *Context
}

func (cs ContextScope) GetEngine() *Engine {
	return cs.context.engine
}

func (cs ContextScope) GetPrivateData() interface{} {
	return cs.context.GetPrivateData()
}

func (cs ContextScope) SetPrivateData(data interface{}) {
	cs.context.SetPrivateData(data)
}

func (e *Engine) NewContext(globalTemplate *ObjectTemplate) *Context {
	var globalTemplatePtr unsafe.Pointer
	if globalTemplate != nil {
		globalTemplatePtr = globalTemplate.self
	}
	self := C.V8_NewContext(e.self, globalTemplatePtr)

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

//export context_scope_callback
func context_scope_callback(c unsafe.Pointer, callback unsafe.Pointer) {
	(*(*func(ContextScope))(callback))(ContextScope{(*Context)(c)})
}

func (c *Context) Scope(callback func(ContextScope)) {
	C.V8_Context_Scope(c.self, unsafe.Pointer(c), unsafe.Pointer(&callback))
}

//export try_catch_callback
func try_catch_callback(callback unsafe.Pointer) {
	(*(*func())(callback))()
}

func (cs ContextScope) ThrowException(err string) {
	cs.context.engine.Compile([]byte(`throw "`+err+`"`), nil, nil).Run()
	//
	// TODO: use Isolate::ThrowException() will make FunctionTemplate::GetFunction() returns NULL, why?
	//
	//errPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&err)).Data)
	//C.V8_Context_ThrowException(c.self, (*C.char)(errPtr), C.int(len(err)))
}

func (cs ContextScope) TryCatch(simple bool, callback func()) string {
	isSimple := 0
	if simple {
		isSimple = 1
	}
	creport := C.V8_Context_TryCatch(cs.context.self, unsafe.Pointer(&callback), C.int(isSimple))
	if creport == nil {
		return ""
	}
	report := C.GoString(creport)
	C.free(unsafe.Pointer(creport))
	return report
}

type MessageCallback func(message string, data interface{})

func (cs ContextScope) AddMessageListener(simple bool, callback MessageCallback, data interface{}) {
	var goSimple int
	if simple {
		goSimple = 1
	}

	var callbackPointer unsafe.Pointer
	if callback != nil {
		callbackPointer = unsafe.Pointer(&callback)
	}
	C.V8_AddMessageListener(
		callbackPointer,
		unsafe.Pointer(&data),
		C.int(goSimple))
}

//export go_message_callback
func go_message_callback(message, callback, data unsafe.Pointer) {
	report := C.GoString((*C.char)(message))
	C.free(message)
	(*(*MessageCallback)(callback))(report, *(*interface{})(data))
}

func (cs ContextScope) Global() *Object {
	return newValue(C.V8_Context_Global(cs.context.self)).ToObject()
}
