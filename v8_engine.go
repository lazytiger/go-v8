package v8

/*
#include "v8_wrap.h"
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

func (e *Engine) ParseJSON(json string) *Value {
	jsonPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&json)).Data)
	return newValue(C.V8_ParseJSON(e.self, (*C.char)(jsonPtr), C.int(len(json))))
}
