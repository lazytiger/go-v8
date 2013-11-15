package v8

/*
#include "v8_wrap.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "reflect"
import "runtime"

// A compiled JavaScript script.
//
type Script struct {
	self unsafe.Pointer
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

// Runs the script returning the resulting value.
//
func (s Script) Run(c *Context) *Value {
	return newValue(C.V8_RunScript(c.self, s.self))
}

// Pre-compilation data that can be associated with a script.  This
// data can be calculated for a script in advance of actually
// compiling it, and can be stored between compilations.  When script
// data is given to the compile method compilation will be faster.
//
type ScriptData struct {
	self unsafe.Pointer
}

func newScriptData(self unsafe.Pointer) *ScriptData {
	if self == nil {
		return nil
	}

	result := &ScriptData{
		self: self,
	}

	runtime.SetFinalizer(result, func(s *ScriptData) {
		if traceDispose {
			println("v8.ScriptData.Dispose()")
		}
		C.V8_DisposeScriptData(s.self)
	})

	return result
}

// Load previous pre-compilation data.
//
func NewScriptData(data []byte) *ScriptData {
	return newScriptData(C.V8_NewScriptData(
		(*C.char)((unsafe.Pointer)(((*reflect.SliceHeader)(unsafe.Pointer(&data))).Data)),
		C.int(len(data)),
	))
}

// Returns the length of Data().
//
func (sd *ScriptData) Length() int {
	return int(C.V8_ScriptDataLength(sd.self))
}

// Returns a serialized representation of this ScriptData that can later be
// passed to New(). NOTE: Serialized data is platform-dependent.
//
func (sd *ScriptData) Data() []byte {
	return C.GoBytes(
		unsafe.Pointer(C.V8_ScriptDataGetData(sd.self)),
		C.V8_ScriptDataLength(sd.self),
	)
}

// Returns true if the source code could not be parsed.
//
func (sd *ScriptData) HasError() bool {
	return C.V8_ScriptDataHasError(sd.self) == 1
}

// The origin, within a file, of a script.
//
type ScriptOrigin struct {
	self         unsafe.Pointer
	Name         string
	LineOffset   int
	ColumnOffset int
}

func (e *Engine) NewScriptOrigin(name string, lineOffset, columnOffset int) *ScriptOrigin {
	namePtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&name)).Data)
	self := C.V8_NewScriptOrigin(e.self, (*C.char)(namePtr), C.int(len(name)), C.int(lineOffset), C.int(columnOffset))

	if self == nil {
		return nil
	}

	result := &ScriptOrigin{
		self:         self,
		Name:         name,
		LineOffset:   lineOffset,
		ColumnOffset: columnOffset,
	}

	runtime.SetFinalizer(result, func(so *ScriptOrigin) {
		if traceDispose {
			println("v8.ScriptOrigin.Dispose()")
		}
		C.V8_DisposeScriptOrigin(so.self)
	})

	return result
}
