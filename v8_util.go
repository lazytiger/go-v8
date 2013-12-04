package v8

/*
#include "v8_wrap.h"
*/
import "C"
import "unsafe"
import "reflect"
import "runtime"

var (
	jsonObjectBegin = []byte("{")
	jsonObjectEnd   = []byte("}")
	jsonColon       = []byte(":")
	jsonComma       = []byte(",")
	jsonQuote       = []byte("\"")
	jsonArrayBegin  = []byte("[")
	jsonArrayEnd    = []byte("]")
	jsonTrue        = []byte("true")
	jsonFalse       = []byte("false")
	jsonNull        = []byte("null")
)

func (cs ContextScope) Eval(code string) *Value {
	if script := cs.context.engine.Compile([]byte(code), nil, nil); script != nil {
		return script.Run()
	}
	return nil
}

func (cs ContextScope) ParseJSON(json string) *Value {
	jsonPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&json)).Data)
	return newValue(C.V8_ParseJSON(cs.context.self, (*C.char)(jsonPtr), C.int(len(json))))
}

func ToJSON(value *Value) []byte {
	return AppendJSON(make([]byte, 0, 1024), value)
}

func AppendJSON(dst []byte, value *Value) []byte {
	switch {
	case value.IsArray():
		dst = append(dst, jsonArrayBegin...)
		array := value.ToArray()
		length := array.Length()
		for i := 0; i < length; i++ {
			dst = AppendJSON(dst, array.GetElement(i))
			if i < length-1 {
				dst = append(dst, jsonComma...)
			}
		}
		dst = append(dst, jsonArrayEnd...)
	case value.IsObject():
		dst = append(dst, jsonObjectBegin...)
		object := value.ToObject()
		names := object.GetOwnPropertyNames()
		length := names.Length()
		for i := 0; i < length; i++ {
			name := names.GetElement(i).ToString()
			dst = append(dst, jsonQuote...)
			dst = append(dst, name...)
			dst = append(dst, jsonQuote...)
			dst = append(dst, jsonColon...)
			dst = AppendJSON(dst, object.GetProperty(name))
			if i < length-1 {
				dst = append(dst, jsonComma...)
			}
		}
		dst = append(dst, jsonObjectEnd...)
	case value.IsString():
		dst = append(dst, jsonQuote...)
		str := value.ToString()
		for i := 0; i < len(str); i++ {
			c := str[i]
			switch c {
			case '"':
				dst = append(dst, '\\', '"')
			case '\\':
				dst = append(dst, '\\', '\\')
			case '/':
				dst = append(dst, '\\', '/')
			case '\n':
				dst = append(dst, '\\', 'n')
			case '\r':
				dst = append(dst, '\\', 'r')
			case '\t':
				dst = append(dst, '\\', 't')
			case '\b':
				dst = append(dst, '\\', 'b')
			case '\f':
				dst = append(dst, '\\', 'f')
			default:
				dst = append(dst, c)
			}
		}
		dst = append(dst, jsonQuote...)
	case value.IsNumber():
		dst = append(dst, value.ToString()...)
	case value.IsTrue():
		dst = append(dst, jsonTrue...)
	case value.IsFalse():
		dst = append(dst, jsonFalse...)
	case value.IsNull():
		dst = append(dst, jsonNull...)
	}

	return dst
}

func GetVersion() string {
	return C.GoString(C.V8_GetVersion())
}

func SetFlagsFromString(cmd string) {
	cs := C.CString(cmd)
	defer C.free(unsafe.Pointer(cs))
	C.V8_SetFlagsFromString(cs, C.int(len(cmd)))
}

type ArrayBufferAllocateCallback func(int, bool) unsafe.Pointer
type ArrayBufferFreeCallback func(unsafe.Pointer, int)

type ArrayBufferAllocator struct {
	self unsafe.Pointer
}

// Call this to get a new ArrayBufferAllocator
func newArrayBufferAllocator() *ArrayBufferAllocator {
	allocator := &ArrayBufferAllocator{}
	runtime.SetFinalizer(allocator, func(allocator *ArrayBufferAllocator) {
		if allocator.self == nil {
			return
		}
		if traceDispose {
			println("dispose array buffer allocator", allocator.self)
		}
		C.V8_Dispose_Allocator(allocator.self)
	})
	return allocator
}

// Call SetArrayBufferAllocator first if you want use any of
// ArrayBuffer, ArrayBufferView, Int8Array...
// Please be sure to call this function once and keep allocator
// Please set ac and fc to nil if you don't want a custom one
func SetArrayBufferAllocator(
	ac ArrayBufferAllocateCallback,
	fc ArrayBufferFreeCallback) {
	var acPointer, fcPointer unsafe.Pointer
	if ac != nil {
		acPointer = unsafe.Pointer(&ac)
	}
	if fc != nil {
		fcPointer = unsafe.Pointer(&fc)
	}

	gMutex.Lock()
	defer gMutex.Unlock()
	gAllocator.self = C.V8_SetArrayBufferAllocator(
		gAllocator.self,
		acPointer,
		fcPointer)
}

//export go_array_buffer_allocate
func go_array_buffer_allocate(callback unsafe.Pointer, length C.size_t, initialized C.int) unsafe.Pointer {
	return (*(*ArrayBufferAllocateCallback)(callback))(int(length), initialized != 0)
}

//export go_array_buffer_free
func go_array_buffer_free(callback unsafe.Pointer, data unsafe.Pointer, length C.size_t) {
	(*(*ArrayBufferFreeCallback)(callback))(data, int(length))
}
