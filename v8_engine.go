package v8

/*
#include "v8_wrap.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "runtime"
import "reflect"

var traceDispose = false

// Represents an isolated instance of the V8 engine.
// Objects from one engine must not be used in other engine.
type Engine struct {
	self             unsafe.Pointer
	_undefined       *Value
	_null            *Value
	_true            *Value
	_false           *Value
	funcTemplateId   int
	funcTemplates    map[int]*FunctionTemplate
	objectTemplateId int
	objectTemplates  map[int]*ObjectTemplate
}

func NewEngine() *Engine {
	self := C.V8_NewEngine()

	if self == nil {
		return nil
	}

	result := &Engine{
		self:            self,
		funcTemplates:   make(map[int]*FunctionTemplate),
		objectTemplates: make(map[int]*ObjectTemplate),
	}

	runtime.SetFinalizer(result, func(e *Engine) {
		if traceDispose {
			println("v8.Engine.Dispose()", e.self)
		}
		C.V8_DisposeEngine(e.self)
	})

	return result
}

//export try_catch_callback
func try_catch_callback(callback unsafe.Pointer) {
	(*(*func())(callback))()
}

func (e *Engine) ThrowException(err string) {
	errPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&err)).Data)
	C.V8_ThrowException(e.self, (*C.char)(errPtr), C.int(len(err)))
}

func (e *Engine) TryCatch(simple bool, callback func()) string {
	isSimple := 0
	if simple {
		isSimple = 1
	}
	creport := C.V8_TryCatch(e.self, unsafe.Pointer(&callback), C.int(isSimple))
	if creport == nil {
		return ""
	}
	report := C.GoString(creport)
	C.free(unsafe.Pointer(creport))
	return report
}
