package v8

/*
#include "v8_wrap.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "runtime"

type Value struct {
	self    unsafe.Pointer
	isType  int
	notType int
}

func newValue(self unsafe.Pointer) *Value {
	if self == nil {
		return nil
	}

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

func (v *Value) GetBoolean() bool {
	return C.V8_ValueGetBoolean(v.self) == 1
}

func (v *Value) GetNumber() float64 {
	return float64(C.V8_ValueGetNumber(v.self))
}

func (v *Value) GetInteger() int64 {
	return int64(C.V8_ValueGetInteger(v.self))
}

func (v *Value) GetUint32() uint32 {
	return uint32(C.V8_ValueGetUint32(v.self))
}

func (v *Value) GetInt32() int32 {
	return int32(C.V8_ValueGetInt32(v.self))
}

const (
	isUndefined     = 1 << iota
	isNull          = 1 << iota
	isTrue          = 1 << iota
	isFalse         = 1 << iota
	isString        = 1 << iota
	isFunction      = 1 << iota
	isArray         = 1 << iota
	isObject        = 1 << iota
	isBoolean       = 1 << iota
	isNumber        = 1 << iota
	isExternal      = 1 << iota
	isInt32         = 1 << iota
	isUint32        = 1 << iota
	isDate          = 1 << iota
	isBooleanObject = 1 << iota
	isNumberObject  = 1 << iota
	isStringObject  = 1 << iota
	isNativeError   = 1 << iota
	isRegExp        = 1 << iota
)

func (v *Value) checkJsType(typeCode int, check func(unsafe.Pointer) bool) bool {
	if (v.isType & typeCode) == typeCode {
		return true
	}

	if (v.notType & typeCode) == typeCode {
		return false
	}

	if check(v.self) {
		v.isType |= typeCode
		return true
	} else {
		v.notType |= typeCode
		return false
	}
}

func (v *Value) IsUndefined() bool {
	return v.checkJsType(isUndefined, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsUndefined(self) == 1
	})
}

func (v *Value) IsNull() bool {
	return v.checkJsType(isNull, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsNull(self) == 1
	})
}

func (v *Value) IsTrue() bool {
	return v.checkJsType(isTrue, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsTrue(self) == 1
	})
}

func (v *Value) IsFalse() bool {
	return v.checkJsType(isFalse, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsFalse(self) == 1
	})
}

func (v *Value) IsString() bool {
	return v.checkJsType(isString, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsString(self) == 1
	})
}

func (v *Value) IsFunction() bool {
	return v.checkJsType(isFunction, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsFunction(self) == 1
	})
}

func (v *Value) IsArray() bool {
	return v.checkJsType(isArray, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsArray(self) == 1
	})
}

func (v *Value) IsObject() bool {
	return v.checkJsType(isObject, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsObject(self) == 1
	})
}

func (v *Value) IsBoolean() bool {
	return v.checkJsType(isBoolean, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsBoolean(self) == 1
	})
}

func (v *Value) IsNumber() bool {
	return v.checkJsType(isNumber, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsNumber(self) == 1
	})
}

func (v *Value) IsExternal() bool {
	return v.checkJsType(isExternal, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsExternal(self) == 1
	})
}

func (v *Value) IsInt32() bool {
	return v.checkJsType(isInt32, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsInt32(self) == 1
	})
}

func (v *Value) IsUint32() bool {
	return v.checkJsType(isUint32, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsUint32(self) == 1
	})
}

func (v *Value) IsDate() bool {
	return v.checkJsType(isDate, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsDate(self) == 1
	})
}

func (v *Value) IsBooleanObject() bool {
	return v.checkJsType(isBooleanObject, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsBooleanObject(self) == 1
	})
}

func (v *Value) IsNumberObject() bool {
	return v.checkJsType(isNumberObject, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsNumberObject(self) == 1
	})
}

func (v *Value) IsStringObject() bool {
	return v.checkJsType(isStringObject, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsStringObject(self) == 1
	})
}

func (v *Value) IsNativeError() bool {
	return v.checkJsType(isNativeError, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsNativeError(self) == 1
	})
}

func (v *Value) IsRegExp() bool {
	return v.checkJsType(isRegExp, func(self unsafe.Pointer) bool {
		return C.V8_ValueIsRegExp(self) == 1
	})
}
