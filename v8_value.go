package v8

/*
#include "v8_warp.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "runtime"

type Value struct {
	self unsafe.Pointer
}

func newValue(self unsafe.Pointer) *Value {
	result := &Value{
		self: self,
	}

	runtime.SetFinalizer(result, func(v *Value) {
		if traceDispose {
			println("v8.Value.Dispose()")
		}
		C.V8_DisposeValue(v.self)
	})

	return result
}

func (v *Value) ToString() string {
	cstring := C.V8_ValueToString(v.self)
	gostring := C.GoString(cstring)
	C.free(unsafe.Pointer(cstring))
	return gostring
}

func (v *Value) IsUndefined() bool {
	return C.V8_ValueIsUndefined(v.self) == 1
}

func (v *Value) IsNull() bool {
	return C.V8_ValueIsNull(v.self) == 1
}

func (v *Value) IsTrue() bool {
	return C.V8_ValueIsTrue(v.self) == 1
}

func (v *Value) IsFalse() bool {
	return C.V8_ValueIsFalse(v.self) == 1
}

func (v *Value) IsString() bool {
	return C.V8_ValueIsString(v.self) == 1
}

func (v *Value) IsFunction() bool {
	return C.V8_ValueIsFunction(v.self) == 1
}

func (v *Value) IsArray() bool {
	return C.V8_ValueIsArray(v.self) == 1
}

func (v *Value) IsObject() bool {
	return C.V8_ValueIsObject(v.self) == 1
}

func (v *Value) IsBoolean() bool {
	return C.V8_ValueIsBoolean(v.self) == 1
}

func (v *Value) IsNumber() bool {
	return C.V8_ValueIsNumber(v.self) == 1
}

func (v *Value) IsExternal() bool {
	return C.V8_ValueIsExternal(v.self) == 1
}

func (v *Value) IsInt32() bool {
	return C.V8_ValueIsInt32(v.self) == 1
}

func (v *Value) IsUint32() bool {
	return C.V8_ValueIsUint32(v.self) == 1
}

func (v *Value) IsDate() bool {
	return C.V8_ValueIsDate(v.self) == 1
}

func (v *Value) IsBooleanObject() bool {
	return C.V8_ValueIsBooleanObject(v.self) == 1
}

func (v *Value) IsNumberObject() bool {
	return C.V8_ValueIsNumberObject(v.self) == 1
}

func (v *Value) IsStringObject() bool {
	return C.V8_ValueIsStringObject(v.self) == 1
}

func (v *Value) IsNativeError() bool {
	return C.V8_ValueIsNativeError(v.self) == 1
}

func (v *Value) IsRegExp() bool {
	return C.V8_ValueIsRegExp(v.self) == 1
}

func (v *Value) GetBoolean() bool {
	return C.V8_ValueGetBoolean() == 1
}

func (v *Value) GetNumber() float64 {
	return float64(C.V8_ValueGetNumber())
}

func (v *Value) GetInteger() int64 {
	return int64(C.V8_ValueGetInteger())
}

func (v *Value) GetUint32() uint32 {
	return uint32(C.V8_ValueGetUint32())
}

func (v *Value) GetInt32() int32 {
	return int32(C.V8_ValueGetInt32())
}
