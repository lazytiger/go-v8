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

func (e *Engine) Undefined() *Value {
	if e._undefined == nil {
		e._undefined = newValue(C.V8_Undefined(e.self))
	}
	return e._undefined
}

func (e *Engine) Null() *Value {
	if e._null == nil {
		e._null = newValue(C.V8_Null(e.self))
	}
	return e._null
}

func (e *Engine) True() *Value {
	if e._true == nil {
		e._true = newValue(C.V8_True(e.self))
	}
	return e._true
}

func (e *Engine) False() *Value {
	if e._false == nil {
		e._false = newValue(C.V8_False(e.self))
	}
	return e._false
}

func (e *Engine) NewBoolean(value bool) *Value {
	if value {
		return e.True()
	}
	return e.False()
}

func (e *Engine) NewNumber(value float64) *Value {
	return newValue(C.V8_NewNumber(
		e.self, C.double(value),
	))
}

func (e *Engine) NewInteger(value int64) *Value {
	return newValue(C.V8_NewNumber(
		e.self, C.double(value),
	))
}

func (e *Engine) NewString(value string) *Value {
	valPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&value)).Data)
	return newValue(C.V8_NewString(
		e.self, (*C.char)(valPtr), C.int(len(value)),
	))
}

// Pre-compiles the specified script (context-independent).
//
func (e *Engine) PreCompile(code []byte) *ScriptData {
	codePtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&code)).Data)
	return newScriptData(C.V8_PreCompile(
		e.self, (*C.char)(codePtr), C.int(len(code)),
	))
}

// Compiles the specified script (context-independent).
// 'data' is the Pre-parsing data, as obtained by PreCompile()
// using pre_data speeds compilation if it's done multiple times.
//
func (e *Engine) Compile(code []byte, origin *ScriptOrigin, data *ScriptData) *Script {
	var originPtr unsafe.Pointer
	var dataPtr unsafe.Pointer

	if origin != nil {
		originPtr = origin.self
	}

	if data != nil {
		dataPtr = data.self
	}

	codePtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&code)).Data)
	self := C.V8_Compile(e.self, (*C.char)(codePtr), C.int(len(code)), originPtr, dataPtr)

	if self == nil {
		return nil
	}

	result := &Script{
		self: self,
	}

	runtime.SetFinalizer(result, func(s *Script) {
		if traceDispose {
			println("v8.Script.Dispose()")
		}
		C.V8_DisposeScript(s.self)
	})

	return result
}
