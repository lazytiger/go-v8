package v8

/*
#include "v8_warp.h"
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
