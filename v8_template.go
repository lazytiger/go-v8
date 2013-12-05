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
	id                 int
	engine             *Engine
	accessors          map[string]*accessorInfo
	namedInfo          *namedPropertyInfo
	indexedInfo        *indexedPropertyInfo
	properties         map[string]*propertyInfo
	self               unsafe.Pointer
	internalFieldCount int
}

type namedPropertyInfo struct {
	getter     NamedPropertyGetterCallback
	setter     NamedPropertySetterCallback
	deleter    NamedPropertyDeleterCallback
	query      NamedPropertyQueryCallback
	enumerator NamedPropertyEnumeratorCallback
	data       interface{}
}

type indexedPropertyInfo struct {
	getter     IndexedPropertyGetterCallback
	setter     IndexedPropertySetterCallback
	deleter    IndexedPropertyDeleterCallback
	query      IndexedPropertyQueryCallback
	enumerator IndexedPropertyEnumeratorCallback
	data       interface{}
}

type accessorInfo struct {
	key     string
	getter  AccessorGetterCallback
	setter  AccessorSetterCallback
	data    interface{}
	attribs PropertyAttribute
}

type NamedPropertyGetterCallback func(string, PropertyCallbackInfo)
type NamedPropertySetterCallback func(string, *Value, PropertyCallbackInfo)
type NamedPropertyDeleterCallback func(string, PropertyCallbackInfo)
type NamedPropertyQueryCallback func(string, PropertyCallbackInfo)
type NamedPropertyEnumeratorCallback func(PropertyCallbackInfo)

type IndexedPropertyGetterCallback func(uint32, PropertyCallbackInfo)
type IndexedPropertySetterCallback func(uint32, *Value, PropertyCallbackInfo)
type IndexedPropertyDeleterCallback func(uint32, PropertyCallbackInfo)
type IndexedPropertyQueryCallback func(uint32, PropertyCallbackInfo)
type IndexedPropertyEnumeratorCallback func(PropertyCallbackInfo)

type propertyInfo struct {
	key     string
	value   *Value
	attribs PropertyAttribute
}

func newObjectTemplate(e *Engine, self unsafe.Pointer) *ObjectTemplate {
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

func (e *Engine) NewObjectTemplate() *ObjectTemplate {
	self := C.V8_NewObjectTemplate(e.self)

	return newObjectTemplate(e, self)
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

func (ot *ObjectTemplate) SetInternalFieldCount(count int) {
	C.V8_ObjectTemplate_SetInternalFieldCount(ot.self, C.int(count))
	ot.internalFieldCount = count
}

func (ot *ObjectTemplate) InternalFieldCount() int {
	return ot.internalFieldCount
}

func (ot *ObjectTemplate) SetAccessor(
	key string,
	getter AccessorGetterCallback,
	setter AccessorSetterCallback,
	data interface{},
	attribs PropertyAttribute,
) {
	info := &accessorInfo{
		key:     key,
		getter:  getter,
		setter:  setter,
		data:    data,
		attribs: attribs,
	}

	ot.accessors[key] = info

	var getterPointer, setterPointer unsafe.Pointer
	if info.getter != nil {
		getterPointer = unsafe.Pointer(&info.getter)
	}

	if info.setter != nil {
		setterPointer = unsafe.Pointer(&info.setter)
	}

	keyPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&info.key)).Data)

	C.V8_ObjectTemplate_SetAccessor(
		ot.self,
		(*C.char)(keyPtr), C.int(len(info.key)),
		getterPointer,
		setterPointer,
		unsafe.Pointer(&(info.data)),
		C.int(info.attribs),
	)
}

func (ot *ObjectTemplate) SetNamedPropertyHandler(
	getter NamedPropertyGetterCallback,
	setter NamedPropertySetterCallback,
	query NamedPropertyQueryCallback,
	deleter NamedPropertyDeleterCallback,
	enumerator NamedPropertyEnumeratorCallback,
	data interface{},
) {
	info := &namedPropertyInfo{
		getter:     getter,
		setter:     setter,
		query:      query,
		deleter:    deleter,
		enumerator: enumerator,
		data:       data,
	}

	ot.namedInfo = info
	var getterPointer, setterPointer, queryPointer, deleterPointer, enumeratorPointer unsafe.Pointer
	if info.getter != nil {
		getterPointer = unsafe.Pointer(&info.getter)
	}

	if info.setter != nil {
		setterPointer = unsafe.Pointer(&info.setter)
	}

	if info.query != nil {
		queryPointer = unsafe.Pointer(&info.query)
	}

	if info.deleter != nil {
		deleterPointer = unsafe.Pointer(&info.deleter)
	}

	if info.enumerator != nil {
		enumeratorPointer = unsafe.Pointer(&info.enumerator)
	}

	C.V8_ObjectTemplate_SetNamedPropertyHandler(
		ot.self,
		getterPointer,
		setterPointer,
		queryPointer,
		deleterPointer,
		enumeratorPointer,
		unsafe.Pointer(&data))
}

func (ot *ObjectTemplate) SetIndexedPropertyHandler(
	getter IndexedPropertyGetterCallback,
	setter IndexedPropertySetterCallback,
	query IndexedPropertyQueryCallback,
	deleter IndexedPropertyDeleterCallback,
	enumerator IndexedPropertyEnumeratorCallback,
	data interface{},
) {
	info := &indexedPropertyInfo{
		getter:     getter,
		setter:     setter,
		query:      query,
		deleter:    deleter,
		enumerator: enumerator,
		data:       data,
	}

	var getterPointer, setterPointer, queryPointer, deleterPointer, enumeratorPointer unsafe.Pointer
	if info.getter != nil {
		getterPointer = unsafe.Pointer(&info.getter)
	}

	if info.setter != nil {
		setterPointer = unsafe.Pointer(&info.setter)
	}

	if info.query != nil {
		queryPointer = unsafe.Pointer(&info.query)
	}

	if info.deleter != nil {
		deleterPointer = unsafe.Pointer(&info.deleter)
	}

	if info.enumerator != nil {
		enumeratorPointer = unsafe.Pointer(&info.enumerator)
	}

	ot.indexedInfo = info

	C.V8_ObjectTemplate_SetIndexedPropertyHandler(
		ot.self,
		getterPointer,
		setterPointer,
		queryPointer,
		deleterPointer,
		enumeratorPointer,
		unsafe.Pointer(&data))
}

type PropertyCallbackInfo struct {
	self        unsafe.Pointer
	typ         C.PropertyDataEnum
	data        interface{}
	returnValue ReturnValue
	context     *Context
}

func (p PropertyCallbackInfo) CurrentScope() ContextScope {
	return ContextScope{p.context}
}

func (p PropertyCallbackInfo) This() *Object {
	return newValue(C.V8_PropertyCallbackInfo_This(p.self, p.typ)).ToObject()
}

func (p PropertyCallbackInfo) Holder() *Object {
	return newValue(C.V8_PropertyCallbackInfo_Holder(p.self, p.typ)).ToObject()
}

func (p PropertyCallbackInfo) Data() interface{} {
	return p.data
}

func (p PropertyCallbackInfo) ReturnValue() ReturnValue {
	if p.returnValue.self == nil {
		p.returnValue.self = C.V8_PropertyCallbackInfo_ReturnValue(p.self, p.typ)
	}
	return p.returnValue
}

// Property getter callback info
//
type AccessorCallbackInfo struct {
	self        unsafe.Pointer
	data        interface{}
	returnValue ReturnValue
	context     *Context
	typ         C.AccessorDataEnum
}

func (ac AccessorCallbackInfo) CurrentScope() ContextScope {
	return ContextScope{ac.context}
}

func (ac AccessorCallbackInfo) This() *Object {
	return newValue(C.V8_AccessorCallbackInfo_This(ac.self, ac.typ)).ToObject()
}

func (ac AccessorCallbackInfo) Holder() *Object {
	return newValue(C.V8_AccessorCallbackInfo_Holder(ac.self, ac.typ)).ToObject()
}

func (ac AccessorCallbackInfo) Data() interface{} {
	return ac.data
}

func (ac *AccessorCallbackInfo) ReturnValue() ReturnValue {
	if ac.returnValue.self == nil {
		ac.returnValue.self = C.V8_AccessorCallbackInfo_ReturnValue(ac.self, ac.typ)
	}
	return ac.returnValue
}

type AccessorGetterCallback func(name string, info AccessorCallbackInfo)

type AccessorSetterCallback func(name string, value *Value, info AccessorCallbackInfo)

//export go_accessor_callback
func go_accessor_callback(typ C.AccessorDataEnum, info *C.V8_AccessorCallbackInfo, context unsafe.Pointer) {
	name := reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(info.key)),
		Len:  int(info.key_length),
	}
	gname := *((*string)(unsafe.Pointer(&name)))
	switch typ {
	case C.OTA_Getter:
		(*(*AccessorGetterCallback)(info.callback))(
			gname,
			AccessorCallbackInfo{unsafe.Pointer(info), *(*interface{})(info.data), ReturnValue{}, (*Context)(context), typ})
	case C.OTA_Setter:
		(*(*AccessorSetterCallback)(info.callback))(
			gname,
			newValue(info.setValue),
			AccessorCallbackInfo{unsafe.Pointer(info), *(*interface{})(info.data), ReturnValue{}, (*Context)(context), typ})
	default:
		panic("impossible type")
	}
}

//export go_named_property_callback
func go_named_property_callback(typ C.PropertyDataEnum, info *C.V8_PropertyCallbackInfo, context unsafe.Pointer) {
	gname := ""
	if info.key != nil {
		gname = C.GoString(info.key)
	}
	switch typ {
	case C.OTP_Getter:
		(*(*NamedPropertyGetterCallback)(info.callback))(
			gname, PropertyCallbackInfo{unsafe.Pointer(info), typ, *(*interface{})(info.data), ReturnValue{}, (*Context)(context)})
	case C.OTP_Setter:
		(*(*NamedPropertySetterCallback)(info.callback))(
			gname,
			newValue(info.setValue),
			PropertyCallbackInfo{unsafe.Pointer(info), typ, *(*interface{})(info.data), ReturnValue{}, (*Context)(context)})
	case C.OTP_Deleter:
		(*(*NamedPropertyDeleterCallback)(info.callback))(
			gname, PropertyCallbackInfo{unsafe.Pointer(info), typ, *(*interface{})(info.data), ReturnValue{}, (*Context)(context)})
	case C.OTP_Query:
		(*(*NamedPropertyQueryCallback)(info.callback))(
			gname, PropertyCallbackInfo{unsafe.Pointer(info), typ, *(*interface{})(info.data), ReturnValue{}, (*Context)(context)})
	case C.OTP_Enumerator:
		(*(*NamedPropertyEnumeratorCallback)(info.callback))(
			PropertyCallbackInfo{unsafe.Pointer(info), typ, *(*interface{})(info.data), ReturnValue{}, (*Context)(context)})
	}
}

//export go_indexed_property_callback
func go_indexed_property_callback(typ C.PropertyDataEnum, info *C.V8_PropertyCallbackInfo, context unsafe.Pointer) {
	switch typ {
	case C.OTP_Getter:
		(*(*IndexedPropertyGetterCallback)(info.callback))(
			uint32(info.index), PropertyCallbackInfo{unsafe.Pointer(info), typ, *(*interface{})(info.data), ReturnValue{}, (*Context)(context)})
	case C.OTP_Setter:
		(*(*IndexedPropertySetterCallback)(info.callback))(
			uint32(info.index),
			newValue(info.setValue),
			PropertyCallbackInfo{unsafe.Pointer(info), typ, *(*interface{})(info.data), ReturnValue{}, (*Context)(context)})
	case C.OTP_Deleter:
		(*(*IndexedPropertyDeleterCallback)(info.callback))(
			uint32(info.index), PropertyCallbackInfo{unsafe.Pointer(info), typ, *(*interface{})(info.data), ReturnValue{}, (*Context)(context)})
	case C.OTP_Query:
		(*(*IndexedPropertyQueryCallback)(info.callback))(
			uint32(info.index), PropertyCallbackInfo{unsafe.Pointer(info), typ, *(*interface{})(info.data), ReturnValue{}, (*Context)(context)})
	case C.OTP_Enumerator:
		(*(*IndexedPropertyEnumeratorCallback)(info.callback))(
			PropertyCallbackInfo{unsafe.Pointer(info), typ, *(*interface{})(info.data), ReturnValue{}, (*Context)(context)})
	}
}

func (o *Object) setAccessor(info *accessorInfo) {
	keyPtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&info.key)).Data)
	var getterPointer, setterPointer unsafe.Pointer
	if info.getter != nil {
		getterPointer = unsafe.Pointer(&info.getter)
	}

	if info.setter != nil {
		setterPointer = unsafe.Pointer(&info.setter)
	}
	C.V8_Object_SetAccessor(
		o.self,
		(*C.char)(keyPtr), C.int(len(info.key)),
		getterPointer,
		setterPointer,
		unsafe.Pointer(&(info.data)),
		C.int(info.attribs),
	)
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
	data     interface{}
	self     unsafe.Pointer
}

func (e *Engine) NewFunctionTemplate(callback FunctionCallback, data interface{}) *FunctionTemplate {
	ft := &FunctionTemplate{
		id:       e.funcTemplateId + 1,
		engine:   e,
		callback: callback,
		data:     data,
	}

	var callbackPtr unsafe.Pointer

	if callback != nil {
		callbackPtr = unsafe.Pointer(&(ft.callback))
	}

	self := C.V8_NewFunctionTemplate(e.self, callbackPtr, unsafe.Pointer(&data))
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

func (ft *FunctionTemplate) SetClassName(name string) {
	ft.Lock()
	defer ft.Unlock()

	if ft.engine == nil {
		panic("engine can't be nil")
	}

	namePtr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&name)).Data)
	C.V8_FunctionTemplate_SetClassName(ft.self, (*C.char)(namePtr), C.int(len(name)))
}

func (ft *FunctionTemplate) InstanceTemplate() *ObjectTemplate {
	ft.Lock()
	defer ft.Unlock()

	if ft.engine == nil {
		panic("engine can't be nil")
	}

	self := C.V8_FunctionTemplate_InstanceTemplate(ft.self)
	return newObjectTemplate(ft.engine, self)
}

//export go_function_callback
func go_function_callback(info, callback, context, data unsafe.Pointer) {
	callbackFunc := *(*func(FunctionCallbackInfo))(callback)
	callbackFunc(FunctionCallbackInfo{info, ReturnValue{}, (*Context)(context), *(*interface{})(data)})
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
	context     *Context
	data        interface{}
}

func (fc FunctionCallbackInfo) CurrentScope() ContextScope {
	return ContextScope{fc.context}
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

func (fc FunctionCallbackInfo) Data() interface{} {
	return fc.data
}

func (fc *FunctionCallbackInfo) ReturnValue() ReturnValue {
	if fc.returnValue.self == nil {
		fc.returnValue.self = C.V8_FunctionCallbackInfo_ReturnValue(fc.self)
	}
	return fc.returnValue
}
