package v8

/*
#include "v8_wrap.h"
*/
import "C"
import "unsafe"
import "reflect"
import "sync"

type AccessControl int

// Access control specifications.
//
// Some accessors should be accessible across contexts.  These
// accessors have an explicit access control parameter which specifies
// the kind of cross-context access that should be allowed.
//
// Additionally, for security, accessors can prohibit overwriting by
// accessors defined in JavaScript.  For objects that have such
// accessors either locally or in their prototype chain it is not
// possible to overwrite the accessor by using __defineGetter__ or
// __defineSetter__ from JavaScript code.
//
const (
	AC_DEFAULT               AccessControl = 0
	AC_ALL_CAN_READ                        = 1
	AC_ALL_CAN_WRITE                       = 1 << 1
	AC_PROHIBITS_OVERWRITING               = 1 << 2
)

type ObjectTemplate struct {
	sync.Mutex
	id         int
	engine     *Engine
	accessors  map[string]*accessorInfo
	properties map[string]*propertyInfo
	self       unsafe.Pointer
}

type accessorInfo struct {
	key     string
	getter  GetterCallback
	setter  SetterCallback
	data	interface{}
	attribs PropertyAttribute
}

type propertyInfo struct {
	key     string
	value   *Value
	attribs PropertyAttribute
}

func (e *Engine) NewObjectTemplate() *ObjectTemplate {
	self := C.V8_NewObjectTemplate(e.self)

	if self == nil {
		return nil
	}

	ot := &ObjectTemplate{
		id:         e.objectTemplateId + 1,
		engine:     e,
		accessors:  make(map[string]*accessorInfo),
		properties: make(map[string]*propertyInfo),
		self:       self,
	}

	e.objectTemplateId += 1
	e.objectTemplates[ot.id] = ot

	return ot
}

func (ot *ObjectTemplate) Dispose() {
	ot.Lock()
	defer ot.Unlock()

	if ot.id > 0 {
		delete(ot.engine.objectTemplates, ot.id)
		ot.id = 0
		ot.engine = nil
		C.V8_DisposeObjectTemplate(ot.self)
	}
}

func (ot *ObjectTemplate) NewObject() *Value {
	ot.Lock()
	defer ot.Unlock()

	if ot.engine == nil {
		return nil
	}

	return newValue(C.V8_ObjectTemplate_NewObject(ot.self))
}

func (ot *ObjectTemplate) WrapObject(value *Value) {
	ot.Lock()
	defer ot.Unlock()

	object := value.ToObject()

	for _, info := range ot.accessors {
		object.setAccessor(info)
	}

	for _, info := range ot.properties {
		object.SetProperty(info.key, info.value, info.attribs)
	}
}

func (ot *ObjectTemplate) SetProperty(key string, value *Value, attribs PropertyAttribute) {
	info := &propertyInfo{
		key:     key,
		value:   value,
		attribs: attribs,
	}

	ot.properties[key] = info

	keyPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&info.key)).Data)

	C.V8_ObjectTemplate_SetProperty(
		ot.self, (*C.char)(keyPtr), C.int(len(key)), value.self, C.int(attribs),
	)
}

func (ot *ObjectTemplate) SetAccessor(key string, getter GetterCallback, setter SetterCallback, data interface{}, attribs PropertyAttribute) {
	info := &accessorInfo{
		key:     key,
		getter:  getter,
		setter:  setter,
		data:	 data,
		attribs: attribs,
	}

	ot.accessors[key] = info

	keyPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&info.key)).Data)

	C.V8_ObjectTemplate_SetAccessor(
		ot.self,
		(*C.char)(keyPtr), C.int(len(info.key)),
		unsafe.Pointer(&(info.getter)),
		unsafe.Pointer(&(info.setter)),
		unsafe.Pointer(&(info.data)),
		C.int(info.attribs),
	)
}

// Property getter callback info
//
type GetterCallbackInfo struct {
	self        unsafe.Pointer
	returnValue ReturnValue
}

func (g GetterCallbackInfo) This() *Object {
	return newValue(C.V8_GetterCallbackInfo_This(g.self)).ToObject()
}

func (g GetterCallbackInfo) Holder() *Object {
	return newValue(C.V8_GetterCallbackInfo_Holder(g.self)).ToObject()
}

func (g GetterCallbackInfo) Data() interface{} {
	return *(*interface{})(C.V8_GetterCallbackInfo_Data(g.self))
}

func (g *GetterCallbackInfo) ReturnValue() ReturnValue {
	if g.returnValue.self == nil {
		g.returnValue.self = C.V8_GetterCallbackInfo_ReturnValue(g.self)
	}
	return g.returnValue
}

// Property setter callback info
//
type SetterCallbackInfo struct {
	self unsafe.Pointer
}

func (s SetterCallbackInfo) This() *Object {
	return newValue(C.V8_SetterCallbackInfo_This(s.self)).ToObject()
}

func (s SetterCallbackInfo) Holder() *Object {
	return newValue(C.V8_SetterCallbackInfo_Holder(s.self)).ToObject()
}

func (s SetterCallbackInfo) Data() interface{} {
	return *(*interface{})(C.V8_SetterCallbackInfo_Data(s.self))
}

type GetterCallback func(name string, info GetterCallbackInfo)

type SetterCallback func(name string, value *Value, info SetterCallbackInfo)

//export go_getter_callback
func go_getter_callback(key *C.char, length C.int, info, callback unsafe.Pointer) {
	name := reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(key)),
		Len:  int(length),
	}
	gname := *((*string)(unsafe.Pointer(&name)))
	(*(*GetterCallback)(callback))(gname, GetterCallbackInfo{info, ReturnValue{}})
}

//export go_setter_callback
func go_setter_callback(key *C.char, length C.int, value, info, callback unsafe.Pointer) {
	name := reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(key)),
		Len:  int(length),
	}
	gname := *((*string)(unsafe.Pointer(&name)))
	(*(*SetterCallback)(callback))(gname, newValue(value), SetterCallbackInfo{info})
}

func (o *Object) setAccessor(info *accessorInfo) bool {
	keyPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&info.key)).Data)
	return C.V8_Object_SetAccessor(
		o.self,
		(*C.char)(keyPtr), C.int(len(info.key)),
		unsafe.Pointer(&(info.getter)),
		unsafe.Pointer(&(info.setter)),
		unsafe.Pointer(&(info.data)),
		C.int(info.attribs),
	) == 1
}

// A JavaScript function object (ECMA-262, 15.3).
//
type Function struct {
	*Object
}

type FunctionCallback func(FunctionCallbackInfo)

type FunctionTemplate struct {
	sync.Mutex
	id       int
	engine   *Engine
	callback FunctionCallback
	self     unsafe.Pointer
}

func (e *Engine) NewFunctionTemplate(callback FunctionCallback) *FunctionTemplate {
	ft := &FunctionTemplate{
		id:       e.funcTemplateId + 1,
		engine:   e,
		callback: callback,
	}

	self := C.V8_NewFunctionTemplate(e.self, unsafe.Pointer(&(ft.callback)))
	if self == nil {
		return nil
	}
	ft.self = self

	e.funcTemplateId += 1
	e.funcTemplates[ft.id] = ft

	return ft
}

func (ft *FunctionTemplate) Dispose() {
	ft.Lock()
	defer ft.Unlock()

	if ft.id > 0 {
		delete(ft.engine.funcTemplates, ft.id)
		ft.id = 0
		ft.engine = nil
		C.V8_DisposeFunctionTemplate(ft.self)
	}
}

func (ft *FunctionTemplate) NewFunction() *Value {
	ft.Lock()
	defer ft.Unlock()

	if ft.engine == nil {
		return nil
	}

	return newValue(C.V8_FunctionTemplate_GetFunction(ft.self))
}

//export go_function_callback
func go_function_callback(info, callback unsafe.Pointer) {
	callbackFunc := *(*func(FunctionCallbackInfo))(callback)
	callbackFunc(FunctionCallbackInfo{info, ReturnValue{}})
}

func (f *Function) Call(args ...*Value) *Value {
	argv := make([]unsafe.Pointer, len(args))
	for i, arg := range args {
		argv[i] = arg.self
	}
	return newValue(C.V8_Function_Call(
		f.self, C.int(len(args)),
		unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&argv)).Data),
	))
}

// Function and property return value
//
type ReturnValue struct {
	self unsafe.Pointer
}

func (rv ReturnValue) Set(value *Value) {
	C.V8_ReturnValue_Set(rv.self, value.self)
}

func (rv ReturnValue) SetBoolean(value bool) {
	valueInt := 0
	if value {
		valueInt = 1
	}
	C.V8_ReturnValue_SetBoolean(rv.self, C.int(valueInt))
}

func (rv ReturnValue) SetNumber(value float64) {
	C.V8_ReturnValue_SetNumber(rv.self, C.double(value))
}

func (rv ReturnValue) SetInt32(value int32) {
	C.V8_ReturnValue_SetInt32(rv.self, C.int32_t(value))
}

func (rv ReturnValue) SetUint32(value uint32) {
	C.V8_ReturnValue_SetUint32(rv.self, C.uint32_t(value))
}

func (rv ReturnValue) SetString(value string) {
	valuePtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&value)).Data)
	C.V8_ReturnValue_SetString(rv.self, (*C.char)(valuePtr), C.int(len(value)))
}

func (rv ReturnValue) SetNull() {
	C.V8_ReturnValue_SetNull(rv.self)
}

func (rv ReturnValue) SetUndefined() {
	C.V8_ReturnValue_SetUndefined(rv.self)
}

// Function callback info
//
type FunctionCallbackInfo struct {
	self        unsafe.Pointer
	returnValue ReturnValue
}

func (fc FunctionCallbackInfo) Get(i int) *Value {
	return newValue(C.V8_FunctionCallbackInfo_Get(fc.self, C.int(i)))
}

func (fc FunctionCallbackInfo) Length() int {
	return int(C.V8_FunctionCallbackInfo_Length(fc.self))
}

func (fc FunctionCallbackInfo) Callee() *Function {
	return newValue(C.V8_FunctionCallbackInfo_Callee(fc.self)).ToFunction()
}

func (fc FunctionCallbackInfo) This() *Object {
	return newValue(C.V8_FunctionCallbackInfo_This(fc.self)).ToObject()
}

func (fc FunctionCallbackInfo) Holder() *Object {
	return newValue(C.V8_FunctionCallbackInfo_Holder(fc.self)).ToObject()
}

func (fc *FunctionCallbackInfo) ReturnValue() ReturnValue {
	if fc.returnValue.self == nil {
		fc.returnValue.self = C.V8_FunctionCallbackInfo_ReturnValue(fc.self)
	}
	return fc.returnValue
}
