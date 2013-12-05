package v8

/*
#include "v8_wrap.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "reflect"

type PropertyAttribute int

const (
	PA_None       PropertyAttribute = 0
	PA_ReadOnly                     = 1 << 0
	PA_DontEnum                     = 1 << 1
	PA_DontDelete                   = 1 << 2
)

// A JavaScript object (ECMA-262, 4.3.3)
//
type Object struct {
	*Value
}

func (cs *ContextScope) NewObject() *Value {
	return newValue(C.V8_NewObject(cs.context.self))
}

func (o *Object) SetProperty(key string, value *Value, attribs PropertyAttribute) bool {
	keyPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&key)).Data)
	return C.V8_Object_SetProperty(
		o.self, (*C.char)(keyPtr), C.int(len(key)), value.self, C.int(attribs),
	) == 1
}

func (o *Object) GetProperty(key string) *Value {
	keyPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&key)).Data)
	return newValue(C.V8_Object_GetProperty(
		o.self, (*C.char)(keyPtr), C.int(len(key)),
	))
}

func (o *Object) SetElement(index int, value *Value) bool {
	return C.V8_Object_SetElement(
		o.self, C.uint32_t(index), value.self,
	) == 1
}

func (o *Object) GetElement(index int) *Value {
	return newValue(C.V8_Object_GetElement(o.self, C.uint32_t(index)))
}

func (o *Object) GetPropertyAttributes(key string) PropertyAttribute {
	keyPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&key)).Data)
	return PropertyAttribute(C.V8_Object_GetPropertyAttributes(
		o.self, (*C.char)(keyPtr), C.int(len(key)),
	))
}

func (o *Object) InternalFieldCount() int {
	return int(C.V8_Object_InternalFieldCount(o.self))
}

func (o *Object) GetInternalField(index int) interface{} {
	data := C.V8_Object_GetInternalField(o.self, C.int(index))
	return *(*interface{})(data)
}

func (o *Object) SetInternalField(index int, value interface{}) {
	C.V8_Object_SetInternalField(
		o.self,
		C.int(index),
		unsafe.Pointer(&value),
	)
}

// Sets a local property on this object bypassing interceptors and
// overriding accessors or read-only properties.
//
// Note that if the object has an interceptor the property will be set
// locally, but since the interceptor takes precedence the local property
// will only be returned if the interceptor doesn't return a value.
//
// Note also that this only works for named properties.
func (o *Object) ForceSetProperty(key string, value *Value, attribs PropertyAttribute) bool {
	keyPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&key)).Data)
	return C.V8_Object_ForceSetProperty(o.self,
		(*C.char)(keyPtr), C.int(len(key)), value.self, C.int(attribs),
	) == 1
}

func (o *Object) HasProperty(key string) bool {
	keyPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&key)).Data)
	return C.V8_Object_HasProperty(
		o.self, (*C.char)(keyPtr), C.int(len(key)),
	) == 1
}

func (o *Object) DeleteProperty(key string) bool {
	keyPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&key)).Data)
	return C.V8_Object_DeleteProperty(
		o.self, (*C.char)(keyPtr), C.int(len(key)),
	) == 1
}

// Delete a property on this object bypassing interceptors and
// ignoring dont-delete attributes.
func (o *Object) ForceDeleteProperty(key string) bool {
	keyPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&key)).Data)
	return C.V8_Object_ForceDeleteProperty(
		o.self, (*C.char)(keyPtr), C.int(len(key)),
	) == 1
}

func (o *Object) HasElement(index int) bool {
	return C.V8_Object_HasElement(
		o.self, C.uint32_t(index),
	) == 1
}

func (o *Object) DeleteElement(index int) bool {
	return C.V8_Object_DeleteElement(
		o.self, C.uint32_t(index),
	) == 1
}

// Returns an array containing the names of the enumerable properties
// of this object, including properties from prototype objects.  The
// array returned by this method contains the same values as would
// be enumerated by a for-in statement over this object.
//
func (o *Object) GetPropertyNames() *Array {
	return newValue(C.V8_Object_GetPropertyNames(o.self)).ToArray()
}

// This function has the same functionality as GetPropertyNames but
// the returned array doesn't contain the names of properties from
// prototype objects.
//
func (o *Object) GetOwnPropertyNames() *Array {
	return newValue(C.V8_Object_GetOwnPropertyNames(o.self)).ToArray()
}

// Get the prototype object.  This does not skip objects marked to
// be skipped by __proto__ and it does not consult the security
// handler.
//
func (o *Object) GetPrototype() *Object {
	return newValue(C.V8_Object_GetPrototype(o.self)).ToObject()
}

// Set the prototype object.  This does not skip objects marked to
// be skipped by __proto__ and it does not consult the security
// handler.
//
func (o *Object) SetPrototype(proto *Object) bool {
	return C.V8_Object_SetPrototype(o.self, proto.self) == 1
}

// An instance of the built-in array constructor (ECMA-262, 15.4.2).
//
type Array struct {
	*Object
}

func (cs ContextScope) NewArray(length int) *Array {
	return newValue(C.V8_NewArray(
		cs.context.self, C.int(length),
	)).ToArray()
}

func (a *Array) Length() int {
	return int(C.V8_Array_Length(a.self))
}

type RegExpFlags int

// Regular expression flag bits. They can be or'ed to enable a set
// of flags.
//
const (
	RF_None       RegExpFlags = 0
	RF_Global                 = 1
	RF_IgnoreCase             = 2
	RF_Multiline              = 4
)

type RegExp struct {
	*Object
	pattern       string
	patternCached bool
	flags         RegExpFlags
	flagsCached   bool
}

// Creates a regular expression from the given pattern string and
// the flags bit field. May throw a JavaScript exception as
// described in ECMA-262, 15.10.4.1.
//
// For example,
//   NewRegExp("foo", RF_Global | RF_Multiline)
//
// is equivalent to evaluating "/foo/gm".
//
func (cs ContextScope) NewRegExp(pattern string, flags RegExpFlags) *Value {
	patternPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&pattern)).Data)

	return newValue(C.V8_NewRegExp(
		cs.context.self, (*C.char)(patternPtr), C.int(len(pattern)), C.int(flags),
	))
}

// Returns the value of the source property: a string representing
// the regular expression.
func (r *RegExp) Pattern() string {
	if !r.patternCached {
		cstring := C.V8_RegExp_Pattern(r.self)
		r.pattern = C.GoString(cstring)
		r.patternCached = true
		C.free(unsafe.Pointer(cstring))
	}
	return r.pattern
}

// Returns the flags bit field.
//
func (r *RegExp) Flags() RegExpFlags {
	if !r.flagsCached {
		r.flags = RegExpFlags(C.V8_RegExp_Flags(r.self))
		r.flagsCached = true
	}
	return r.flags
}
