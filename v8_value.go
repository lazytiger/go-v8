package v8

/*
#include "v8_warp.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "runtime"

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

func (v *Value) IsUndefined() bool {
	if v.isType&isUndefined == isUndefined {
		return true
	}
	if v.notType&isUndefined == isUndefined {
		return false
	}
	if C.V8_ValueIsUndefined(v.self) == 1 {
		v.isType |= isUndefined
		return true
	} else {
		v.notType |= isUndefined
		return false
	}
}

func (v *Value) IsNull() bool {
	if v.isType&isNull == isNull {
		return true
	}
	if v.notType&isNull == isNull {
		return false
	}
	if C.V8_ValueIsNull(v.self) == 1 {
		v.isType |= isNull
		return true
	} else {
		v.notType |= isNull
		return false
	}
}

func (v *Value) IsTrue() bool {
	if v.isType&isTrue == isTrue {
		return true
	}
	if v.notType&isTrue == isTrue {
		return false
	}
	if C.V8_ValueIsTrue(v.self) == 1 {
		v.isType |= isTrue
		return true
	} else {
		v.notType |= isTrue
		return false
	}
}

func (v *Value) IsFalse() bool {
	if v.isType&isFalse == isFalse {
		return true
	}
	if v.notType&isFalse == isFalse {
		return false
	}
	if C.V8_ValueIsFalse(v.self) == 1 {
		v.isType |= isFalse
		return true
	} else {
		v.notType |= isFalse
		return false
	}
}

func (v *Value) IsString() bool {
	if v.isType&isString == isString {
		return true
	}
	if v.notType&isString == isString {
		return false
	}
	if C.V8_ValueIsString(v.self) == 1 {
		v.isType |= isString
		return true
	} else {
		v.notType |= isString
		return false
	}
}

func (v *Value) IsFunction() bool {
	if v.isType&isFunction == isFunction {
		return true
	}
	if v.notType&isFunction == isFunction {
		return false
	}
	if C.V8_ValueIsFunction(v.self) == 1 {
		v.isType |= isFunction
		return true
	} else {
		v.notType |= isFunction
		return false
	}
}

func (v *Value) IsArray() bool {
	if v.isType&isArray == isArray {
		return true
	}
	if v.notType&isArray == isArray {
		return false
	}
	if C.V8_ValueIsArray(v.self) == 1 {
		v.isType |= isArray
		return true
	} else {
		v.notType |= isArray
		return false
	}
}

func (v *Value) IsObject() bool {
	if v.isType&isObject == isObject {
		return true
	}
	if v.notType&isObject == isObject {
		return false
	}
	if C.V8_ValueIsObject(v.self) == 1 {
		v.isType |= isObject
		return true
	} else {
		v.notType |= isObject
		return false
	}
}

func (v *Value) IsBoolean() bool {
	if v.isType&isBoolean == isBoolean {
		return true
	}
	if v.notType&isBoolean == isBoolean {
		return false
	}
	if C.V8_ValueIsBoolean(v.self) == 1 {
		v.isType |= isBoolean
		return true
	} else {
		v.notType |= isBoolean
		return false
	}
}

func (v *Value) IsNumber() bool {
	if v.isType&isNumber == isNumber {
		return true
	}
	if v.notType&isNumber == isNumber {
		return false
	}
	if C.V8_ValueIsNumber(v.self) == 1 {
		v.isType |= isNumber
		return true
	} else {
		v.notType |= isNumber
		return false
	}
}

func (v *Value) IsExternal() bool {
	if v.isType&isExternal == isExternal {
		return true
	}
	if v.notType&isExternal == isExternal {
		return false
	}
	if C.V8_ValueIsExternal(v.self) == 1 {
		v.isType |= isExternal
		return true
	} else {
		v.notType |= isExternal
		return false
	}
}

func (v *Value) IsInt32() bool {
	if v.isType&isInt32 == isInt32 {
		return true
	}
	if v.notType&isInt32 == isInt32 {
		return false
	}
	if C.V8_ValueIsInt32(v.self) == 1 {
		v.isType |= isInt32
		return true
	} else {
		v.notType |= isInt32
		return false
	}
}

func (v *Value) IsUint32() bool {
	if v.isType&isUint32 == isUint32 {
		return true
	}
	if v.notType&isUint32 == isUint32 {
		return false
	}
	if C.V8_ValueIsUint32(v.self) == 1 {
		v.isType |= isUint32
		return true
	} else {
		v.notType |= isUint32
		return false
	}
}

func (v *Value) IsDate() bool {
	if v.isType&isDate == isDate {
		return true
	}
	if v.notType&isDate == isDate {
		return false
	}
	if C.V8_ValueIsDate(v.self) == 1 {
		v.isType |= isDate
		return true
	} else {
		v.notType |= isDate
		return false
	}
}

func (v *Value) IsBooleanObject() bool {
	if v.isType&isBooleanObject == isBooleanObject {
		return true
	}
	if v.notType&isBooleanObject == isBooleanObject {
		return false
	}
	if C.V8_ValueIsBooleanObject(v.self) == 1 {
		v.isType |= isBooleanObject
		return true
	} else {
		v.notType |= isBooleanObject
		return false
	}
}

func (v *Value) IsNumberObject() bool {
	if v.isType&isNumberObject == isNumberObject {
		return true
	}
	if v.notType&isNumberObject == isNumberObject {
		return false
	}
	if C.V8_ValueIsNumberObject(v.self) == 1 {
		v.isType |= isNumberObject
		return true
	} else {
		v.notType |= isNumberObject
		return false
	}
}

func (v *Value) IsStringObject() bool {
	if v.isType&isStringObject == isStringObject {
		return true
	}
	if v.notType&isStringObject == isStringObject {
		return false
	}
	if C.V8_ValueIsStringObject(v.self) == 1 {
		v.isType |= isStringObject
		return true
	} else {
		v.notType |= isStringObject
		return false
	}
}

func (v *Value) IsNativeError() bool {
	if v.isType&isNativeError == isNativeError {
		return true
	}
	if v.notType&isNativeError == isNativeError {
		return false
	}
	if C.V8_ValueIsNativeError(v.self) == 1 {
		v.isType |= isNativeError
		return true
	} else {
		v.notType |= isNativeError
		return false
	}
}

func (v *Value) IsRegExp() bool {
	if v.isType&isRegExp == isRegExp {
		return true
	}
	if v.notType&isRegExp == isRegExp {
		return false
	}
	if C.V8_ValueIsRegExp(v.self) == 1 {
		v.isType |= isRegExp
		return true
	} else {
		v.notType |= isRegExp
		return false
	}
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
