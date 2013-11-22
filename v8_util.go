package v8

/*
#include "v8_wrap.h"
*/
import "C"
import "unsafe"
import "reflect"

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

func (e *Engine) ParseJSON(json string) *Value {
	jsonPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&json)).Data)
	return newValue(C.V8_ParseJSON(e.self, (*C.char)(jsonPtr), C.int(len(json))))
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
