#include <stdlib.h>
#include <stdint.h>
#include "v8.h"
#include "v8_wrap.h"

extern "C" {

using namespace v8;

class V8_Context {
public:
	V8_Context(Isolate* isolate, Handle<Context> context) {
		isolate_ = isolate;
		self.Reset(isolate_, context);
	}

	V8_Context(V8_Context* ownerEngine, Handle<Context> context) {
		engine = ownerEngine;
		isolate_ = engine->GetIsolate();
		self.Reset(isolate_, context);
	}

	~V8_Context() {
		Locker locker(isolate_);
		Isolate::Scope isolate_scope(isolate_);

		self.Dispose();
		self.Reset();
	}

	Isolate* GetIsolate() {
		return isolate_;
	}

	Isolate* isolate_;
	V8_Context* engine;
	Persistent<Context> self;
};

class V8_Script {
public:
	V8_Script(V8_Context* ownerEngine, Handle<Script> script) {
		engine = ownerEngine;
		self.Reset(engine->GetIsolate(), script);
	}

	~V8_Script() {
		Locker locker(GetIsolate());
		Isolate::Scope isolate_scope(GetIsolate());

		self.Dispose();
		self.Reset();
	}

	Isolate* GetIsolate() {
		return engine->GetIsolate();
	}

	V8_Context* engine;
	Persistent<Script> self;
};

class V8_Value {
public:
	V8_Value(V8_Context* ownerEngine, Handle<Value> value) {
		engine = ownerEngine;
		self.Reset(engine->GetIsolate(), value);
	}

	~V8_Value() {
		Locker locker(GetIsolate());
		Isolate::Scope isolate_scope(GetIsolate());

		self.Dispose();
		self.Reset();
	}

	Isolate* GetIsolate() {
		return engine->GetIsolate();
	}

	V8_Context* engine;
	Persistent<Value> self;
};

typedef struct V8_ReturnValue {
	V8_ReturnValue(V8_Context* the_engine, ReturnValue<Value> the_value) : 
		engine(the_engine),
		value(the_value) {
	}

	V8_Context*        engine;
	ReturnValue<Value> value;
} V8_ReturnValue;

#define ISOLATE_SCOPE(isolate_ptr) \
	Isolate* isolate = isolate_ptr; \
	Locker locker(isolate); \
	Isolate::Scope isolate_scope(isolate); \
	HandleScope handle_scope(isolate) \

#define ENGINE_SCOPE(engine) \
	V8_Context* the_engine = static_cast<V8_Context*>(engine); \
	ISOLATE_SCOPE(the_engine->GetIsolate()); \
	Local<Context> local_context = Local<Context>::New(isolate, the_engine->self); \
	Context::Scope context_scope(local_context) \

#define VALUE_SCOPE(value) \
	V8_Value* the_value = static_cast<V8_Value*>(value); \
	ISOLATE_SCOPE(the_value->GetIsolate()); \
	Local<Context> context = Local<Context>::New(isolate, the_value->engine->self); \
	Context::Scope context_scope(context); \
	Local<Value> local_value = Local<Value>::New(isolate, the_value->self) \

/*
engine
*/
void* V8_NewEngine() {
	ISOLATE_SCOPE(Isolate::New());

	Handle<Context> context = Context::New(isolate);

	if (context.IsEmpty())
		return NULL;

	return (void*)(new V8_Context(isolate, context));
}

void V8_DisposeEngine(void* engine) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);
	Isolate* isolate = the_engine->GetIsolate();
	delete the_engine;
	isolate->Dispose();
}

/*
context
*/
void* V8_NewContext(void* engine) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);

	ISOLATE_SCOPE(the_engine->GetIsolate());
	
	Handle<Context> context = Context::New(isolate);

	if (context.IsEmpty())
		return NULL;

	return (void*)(new V8_Context(the_engine, context));
}

void V8_DisposeContext(void* context) {
	delete static_cast<V8_Context*>(context);
}

/*
script
*/
void* V8_Compile(void* engine, const char* code, int length, void* script_origin,void* script_data) {
	ENGINE_SCOPE(engine);

	Handle<Script> script = Script::New(
		String::NewFromOneByte(isolate, (uint8_t*)code, String::kNormalString, length), 
		static_cast<ScriptOrigin*>(script_origin), 
		static_cast<ScriptData*>(script_data),
		Handle<String>()
	);

	if (script.IsEmpty())
		return NULL;

	return (void*)(new V8_Script(the_engine, script));
}

void V8_DisposeScript(void* script) {
	delete static_cast<V8_Script*>(script);
}

void* V8_RunScript(void* context, void* script) {
	V8_Context* ctx = static_cast<V8_Context*>(context);
	ISOLATE_SCOPE(ctx->GetIsolate());
	Local<Context> local_context = Local<Context>::New(isolate, ctx->self);
	Context::Scope context_scope(local_context);

	V8_Script* spt = static_cast<V8_Script*>(script);
	Local<Script> local_script = Local<Script>::New(isolate, spt->self);
	
	Handle<Value> result = local_script->Run();

	if (result.IsEmpty())
		return NULL;

	return (void*)(new V8_Value(spt->engine, result));
}

/*
script data
*/
void* V8_PreCompile(void* engine, const char* code, int length) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);
	ISOLATE_SCOPE(the_engine->GetIsolate());

	return (void*)ScriptData::PreCompile(
		String::NewFromOneByte(isolate, (uint8_t*)code, String::kNormalString, length)
	);
}

void* V8_NewScriptData(const char* data, int length) {
	return (void*)ScriptData::New(data, length);
}

void V8_DisposeScriptData(void* script_data) {
	delete static_cast<ScriptData*>(script_data);
}

int V8_ScriptData_Length(void* script_data) {
	return static_cast<ScriptData*>(script_data)->Length();
}

const char* V8_ScriptData_Data(void* script_data) {
	return static_cast<ScriptData*>(script_data)->Data();
}

int V8_ScriptData_HasError(void* script_data) {
	return static_cast<ScriptData*>(script_data)->HasError();
}

/*
script origin
*/
void* V8_NewScriptOrigin(void* engine, const char* name, int name_length, int line_offset, int column_offset) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);
	ISOLATE_SCOPE(the_engine->GetIsolate());

	return (void*)(new ScriptOrigin(
		String::NewFromOneByte(isolate, (uint8_t*)name, String::kNormalString, name_length),
		Integer::New(line_offset),
		Integer::New(line_offset)
	));
}

void V8_DisposeScriptOrigin(void* script_origin) {
	delete static_cast<ScriptOrigin*>(script_origin);
}

/*
Value wrappers
*/
void V8_DisposeValue(void* value) {
	delete static_cast<V8_Value*>(value);
}

int V8_Value_IsUndefined(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsUndefined();
}

int V8_Value_IsNull(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsNull();
}

int V8_Value_IsTrue(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsTrue();
}

int V8_Value_IsFalse(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsFalse();
}

int V8_Value_IsString(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsString();
}

int V8_Value_IsFunction(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsFunction();
}

int V8_Value_IsArray(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsArray();
}

int V8_Value_IsObject(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsObject();
}

int V8_Value_IsBoolean(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsBoolean();
}

int V8_Value_IsNumber(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsNumber();
}

int V8_Value_IsExternal(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsExternal();
}

int V8_Value_IsInt32(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsInt32();
}

int V8_Value_IsUint32(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsUint32();
}

int V8_Value_IsDate(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsDate();
}

int V8_Value_IsBooleanObject(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsBooleanObject();
}

int V8_Value_IsNumberObject(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsNumberObject();
}

int V8_Value_IsStringObject(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsStringObject();
}

int V8_Value_IsNativeError(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsNativeError();
}

int V8_Value_IsRegExp(void* value) {
	VALUE_SCOPE(value);
	return local_value->IsRegExp();
}

int V8_Value_ToBoolean(void* value) {
	VALUE_SCOPE(value);
	return local_value->BooleanValue();
}
  
double V8_Value_ToNumber(void* value) {
	VALUE_SCOPE(value);
	return local_value->NumberValue();
}

int64_t V8_Value_ToInteger(void* value) {
	VALUE_SCOPE(value);
	return local_value->IntegerValue();
}

uint32_t V8_Value_ToUint32(void* value) {
	VALUE_SCOPE(value);
	return local_value->Uint32Value();
}

int32_t V8_Value_ToInt32(void* value) {
	VALUE_SCOPE(value);
	return local_value->Int32Value();
}

char* V8_Value_ToString(void* value) {
	VALUE_SCOPE(value);

	Handle<String> string = local_value->ToString();
	uint8_t* str = (uint8_t*)malloc(string->Length() + 1);
	string->WriteOneByte(str);

	return (char*)str;
}

void* V8_Undefined(void* engine) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);
	ISOLATE_SCOPE(the_engine->GetIsolate());
	return (void*)(new V8_Value(the_engine, Undefined(isolate)));
}

void* V8_Null(void* engine) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);
	ISOLATE_SCOPE(the_engine->GetIsolate());
	return (void*)(new V8_Value(the_engine, Null(isolate)));
}

void* V8_True(void* engine) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);
	ISOLATE_SCOPE(the_engine->GetIsolate());
	return (void*)(new V8_Value(the_engine, True(isolate)));
}

void* V8_False(void* engine) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);
	ISOLATE_SCOPE(the_engine->GetIsolate());
	return (void*)(new V8_Value(the_engine, False(isolate)));
}

void* V8_NewNumber(void* engine, double val) {
	ENGINE_SCOPE(engine);
	
	return (void*)(new V8_Value(the_engine, 
		Number::New(isolate, val)
	));
}

void* V8_NewString(void* engine, const char* val, int val_length) {
	ENGINE_SCOPE(engine);
	
	return (void*)(new V8_Value(the_engine, 
		String::NewFromOneByte(isolate, (uint8_t*)val, String::kNormalString, val_length)
	));
}

/*
object
*/
void* V8_NewObject(void* engine) {
	ENGINE_SCOPE(engine);

	return (void*)(new V8_Value(the_engine, 
		Object::New()
	));
}

int V8_Object_SetProperty(void* value, const char* key, int key_length, void* prop_value, int attribs) {
	VALUE_SCOPE(value);

	return Local<Object>::Cast(local_value)->Set(
		String::NewFromOneByte(isolate, (uint8_t*)key, String::kNormalString, key_length),
		Local<Value>::New(isolate, static_cast<V8_Value*>(prop_value)->self),
		(v8::PropertyAttribute)attribs
	);
}

void* V8_Object_GetProperty(void* value, const char* key, int key_length) {
	VALUE_SCOPE(value);

	return (void*)(new V8_Value(the_value->engine,
		Local<Object>::Cast(local_value)->Get(
			String::NewFromOneByte(isolate, (uint8_t*)key, String::kNormalString, key_length)
		)
	));
}

int V8_Object_SetElement(void* value, uint32_t index, void* elem_value) {
	VALUE_SCOPE(value);

	return Local<Object>::Cast(local_value)->Set(
		index,
		Local<Value>::New(isolate, static_cast<V8_Value*>(elem_value)->self)
	);
}

void* V8_Object_GetElement(void* value, uint32_t index) {
	VALUE_SCOPE(value);

	return (void*)(new V8_Value(the_value->engine,
		Local<Object>::Cast(local_value)->Get(index)
	));
}

int V8_Object_GetPropertyAttributes(void* value, const char* key, int key_length) {
	VALUE_SCOPE(value);

	return Local<Object>::Cast(local_value)->GetPropertyAttributes(
		String::NewFromOneByte(isolate, (uint8_t*)key, String::kNormalString, key_length)
	);
}

int V8_Object_ForceSetProperty(void* value, const char* key, int key_length, void* prop_value, int attribs) {
	VALUE_SCOPE(value);

	return Local<Object>::Cast(local_value)->ForceSet(
		String::NewFromOneByte(isolate, (uint8_t*)key, String::kNormalString, key_length),
		Local<Value>::New(isolate, static_cast<V8_Value*>(prop_value)->self),
		(v8::PropertyAttribute)attribs
	);
}

int V8_Object_HasProperty(void* value, const char* key, int key_length) {
	VALUE_SCOPE(value);

	return Local<Object>::Cast(local_value)->Has(
		String::NewFromOneByte(isolate, (uint8_t*)key, String::kNormalString, key_length)
	);
}

int V8_Object_DeleteProperty(void* value, const char* key, int key_length) {
	VALUE_SCOPE(value);

	return Local<Object>::Cast(local_value)->Delete(
		String::NewFromOneByte(isolate, (uint8_t*)key, String::kNormalString, key_length)
	);
}

int V8_Object_ForceDeleteProperty(void* value, const char* key, int key_length) {
	VALUE_SCOPE(value);

	return Local<Object>::Cast(local_value)->ForceDelete(
		String::NewFromOneByte(isolate, (uint8_t*)key, String::kNormalString, key_length)
	);
}

int V8_Object_HasElement(void* value, uint32_t index) {
	VALUE_SCOPE(value);

	return Local<Object>::Cast(local_value)->Has(index);
}

int V8_Object_DeleteElement(void* value, uint32_t index) {
	VALUE_SCOPE(value);

	return Local<Object>::Cast(local_value)->Delete(index);
}

void* V8_Object_GetPropertyNames(void* value) {
	VALUE_SCOPE(value);

	return (void*)new V8_Value(the_value->engine,
		Local<Object>::Cast(local_value)->GetPropertyNames()
	);
}

void* V8_Object_GetOwnPropertyNames(void* value) {
	VALUE_SCOPE(value);

	return (void*)new V8_Value(the_value->engine,
		Local<Object>::Cast(local_value)->GetOwnPropertyNames()
	);
}

void* V8_Object_GetPrototype(void* value) {
	VALUE_SCOPE(value);

	return (void*)new V8_Value(the_value->engine,
		Local<Object>::Cast(local_value)->GetPrototype()
	);
}

int V8_Object_SetPrototype(void* value, void* proto) {
	VALUE_SCOPE(value);

	return Local<Object>::Cast(local_value)->SetPrototype(
		Local<Value>::New(isolate, static_cast<V8_Value*>(proto)->self)
	);
}

int V8_Object_IsCallable(void* value) {
	VALUE_SCOPE(value);

	return Local<Object>::Cast(local_value)->IsCallable();
}

typedef struct {
	V8_Context*                        engine;
	const PropertyCallbackInfo<Value>* getter_info;
	const PropertyCallbackInfo<void>*  setter_info;
	V8_ReturnValue*                    returnValue;
} V8_AccessorCallbackInfo;

extern void go_getter_callback(char* key, int key_length, void* info, void* callback);
extern void go_setter_callback(char* key, int key_length, void* value, void* info, void* callback);

void V8_GetterCallback(Local<String> property, const PropertyCallbackInfo<Value>& info) {
	v8::Isolate* isolate = v8::Isolate::GetCurrent();
	HandleScope handle_scope(isolate);

	Local<Array> callback_data = Local<Array>::Cast(info.Data());

	V8_AccessorCallbackInfo callback_info;
	callback_info.engine = (V8_Context*)Local<External>::Cast(callback_data->Get(0))->Value();
	callback_info.getter_info = &info;
	callback_info.setter_info = NULL;
	callback_info.returnValue = NULL;

	void* callback = Local<External>::Cast(callback_data->Get(1))->Value();

	void* key = Local<External>::Cast(callback_data->Get(3))->Value();

	int key_length = Local<Integer>::Cast(callback_data->Get(4))->Value();

	go_getter_callback(
		(char*)key, key_length,
		&callback_info, 
		callback
	);

	if (callback_info.returnValue != NULL)
		delete callback_info.returnValue;
}

void V8_SetterCallback(Local<String> property, Local<Value> value, const PropertyCallbackInfo<void>& info) {
	v8::Isolate* isolate = v8::Isolate::GetCurrent();
	HandleScope handle_scope(isolate);

	Local<Array> callback_data = Local<Array>::Cast(info.Data());

	V8_AccessorCallbackInfo callback_info;
	callback_info.engine = (V8_Context*)Local<External>::Cast(callback_data->Get(0))->Value();
	callback_info.getter_info = NULL;
	callback_info.setter_info = &info;
	callback_info.returnValue = NULL;

	void* callback = Local<External>::Cast(callback_data->Get(2))->Value();

	void* key = Local<External>::Cast(callback_data->Get(3))->Value();

	int key_length = Local<Integer>::Cast(callback_data->Get(4))->Value();

	go_setter_callback(
		(char*)key, key_length,
		(void*)new V8_Value(callback_info.engine, value), 
		&callback_info, 
		callback
	);

	if (callback_info.returnValue != NULL)
		delete callback_info.returnValue;
}

int V8_Object_SetAccessor(void *value, const char* key, int key_length, void* getter, void* setter, int attribs) {
	VALUE_SCOPE(value);

	Handle<Array> callback_info = Array::New(3);
	callback_info->Set(0, External::New((void*)the_value->engine));
	callback_info->Set(1, External::New(getter));
	callback_info->Set(2, External::New(setter));
	callback_info->Set(3, External::New((void*)key));
	callback_info->Set(4, Integer::New(key_length));

	if (callback_info.IsEmpty())
		return 0;

	return Local<Object>::Cast(local_value)->SetAccessor(
		String::NewFromOneByte(isolate, (uint8_t*)key, String::kNormalString, key_length),
		V8_GetterCallback, setter == NULL ? NULL : V8_SetterCallback,
 		callback_info
	);
}

void* V8_GetterCallbackInfo_This(void *info) {
	V8_AccessorCallbackInfo* the_info = (V8_AccessorCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return (void*)new V8_Value(the_info->engine, the_info->getter_info->This());
}

void* V8_GetterCallbackInfo_Holder(void *info) {
	V8_AccessorCallbackInfo* the_info = (V8_AccessorCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return (void*)new V8_Value(the_info->engine, the_info->getter_info->Holder());
}

void* V8_GetterCallbackInfo_ReturnValue(void *info) {
	V8_AccessorCallbackInfo* the_info = (V8_AccessorCallbackInfo*)info;
	if (the_info->returnValue == NULL) {
		the_info->returnValue = new V8_ReturnValue(
			the_info->engine, 
			the_info->getter_info->GetReturnValue()
		);;
	}
	return (void*)the_info->returnValue;
}

void* V8_SetterCallbackInfo_This(void *info) {
	V8_AccessorCallbackInfo* the_info = (V8_AccessorCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return (void*)new V8_Value(the_info->engine, the_info->setter_info->This());
}

void* V8_SetterCallbackInfo_Holder(void *info) {
	V8_AccessorCallbackInfo* the_info = (V8_AccessorCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return (void*)new V8_Value(the_info->engine, the_info->setter_info->Holder());
}

/*
array
*/
void* V8_NewArray(void* engine, int length) {
	ENGINE_SCOPE(engine);

	return (void*)(new V8_Value(the_engine, 
		Array::New(length)
	));
}

int V8_Array_Length(void* value) {
	VALUE_SCOPE(value);
	return Local<Array>::Cast(local_value)->Length();
}

/*
regexp
*/
void* V8_NewRegExp(void* engine, const char* pattern, int length, int flags) {
	ENGINE_SCOPE(engine);

	return (void*)(new V8_Value(the_engine, RegExp::New(
		String::NewFromOneByte(isolate, (uint8_t*)pattern, String::kNormalString, length), 
		(RegExp::Flags)flags
	)));
}

char* V8_RegExp_Pattern(void* value) {
	VALUE_SCOPE(value);

	Local<String> pattern = Local<RegExp>::Cast(local_value)->GetSource();

	uint8_t* str = (uint8_t*)malloc(pattern->Length() + 1);
	pattern->WriteOneByte(str);

	return (char*)str;
}

int V8_RegExp_Flags(void* value) {
	VALUE_SCOPE(value);
	return Local<RegExp>::Cast(local_value)->GetFlags();
}

/*
return value
*/
void V8_ReturnValue_Set(void* rv, void* result) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->engine);
	the_rv->value.Set(static_cast<V8_Value*>(result)->self);
}

void V8_ReturnValue_SetBoolean(void* rv, int v) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->engine);
	the_rv->value.Set((bool)v);
}

void V8_ReturnValue_SetNumber(void* rv, double v) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->engine);
	the_rv->value.Set(v);
}

void V8_ReturnValue_SetInt32(void* rv, int32_t v) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->engine);
	the_rv->value.Set(v);
}

void V8_ReturnValue_SetUint32(void* rv, uint32_t v) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->engine);
	the_rv->value.Set(v);
}

void V8_ReturnValue_SetString(void* rv, const char* str, int str_length) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->engine);
	if (str_length == 0) {
		the_rv->value.SetEmptyString();
	} else {
		the_rv->value.Set(
			String::NewFromOneByte(isolate, (uint8_t*)str, String::kNormalString, str_length)
		);
	}
}

void V8_ReturnValue_SetNull(void* rv) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->engine);
	the_rv->value.SetNull();
}

void V8_ReturnValue_SetUndefined(void* rv) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->engine);
	the_rv->value.SetUndefined();
}

/*
function
*/
typedef struct {
	V8_Context*                        engine;
	const FunctionCallbackInfo<Value>* info;
	V8_ReturnValue*                    returnValue;
} V8_FunctionCallbackInfo;

extern void go_function_callback(void* info, void* callback);

void V8_FunctionCallback(const FunctionCallbackInfo<Value>& info) {
	v8::Isolate* isolate = v8::Isolate::GetCurrent();
	HandleScope handle_scope(isolate);

	Local<Array> callback_data = Local<Array>::Cast(info.Data());

	V8_FunctionCallbackInfo callback_info;
	callback_info.engine = (V8_Context*)Local<External>::Cast(callback_data->Get(0))->Value();
	callback_info.info = &info;
	callback_info.returnValue = NULL;

	void* callback = Local<External>::Cast(callback_data->Get(1))->Value();

	go_function_callback(&callback_info, callback);

	if (callback_info.returnValue != NULL)
		delete callback_info.returnValue;
}

void* V8_NewFunction(void* engine, void* callback) {
	ENGINE_SCOPE(engine);

	Handle<Array> callback_data = Array::New(2);

	if (callback_data.IsEmpty())
		return NULL;

	callback_data->Set(0, External::New(engine));
	callback_data->Set(1, External::New(callback));

	return (void*)(new V8_Value(the_engine,
		Function::New(isolate, V8_FunctionCallback, callback_data)
	));
}

void* V8_Function_Call(void* value, int argc, void* argv) {
	VALUE_SCOPE(value);

	Handle<Value>* real_argv = new Handle<Value>[argc];
	V8_Value* *argv_ptr = (V8_Value**)argv;

	for (int i = 0; i < argc; i ++) {
		real_argv[i] = Local<Value>::New(isolate, static_cast<V8_Value*>(argv_ptr[i])->self);
	}

	void* result = (void*)(new V8_Value(the_value->engine, 
		Local<Function>::Cast(local_value)->Call(local_value, argc, real_argv)
	));

	delete[] real_argv;

	return result;
}

void* V8_FunctionCallbackInfo_Get(void* info, int i) {
	V8_FunctionCallbackInfo* the_info = (V8_FunctionCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return (void*)new V8_Value(the_info->engine, (*(the_info->info))[i]);
}

int V8_FunctionCallbackInfo_Length(void* info) {
	V8_FunctionCallbackInfo* the_info = (V8_FunctionCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return the_info->info->Length();
}

void* V8_FunctionCallbackInfo_Callee(void* info) {
	V8_FunctionCallbackInfo* the_info = (V8_FunctionCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return (void*)new V8_Value(the_info->engine, the_info->info->Callee());
}

void* V8_FunctionCallbackInfo_This(void* info) {
	V8_FunctionCallbackInfo* the_info = (V8_FunctionCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return (void*)new V8_Value(the_info->engine, the_info->info->This());
}

void* V8_FunctionCallbackInfo_Holder(void* info) {
	V8_FunctionCallbackInfo* the_info = (V8_FunctionCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return (void*)new V8_Value(the_info->engine, the_info->info->Holder());
}

void* V8_FunctionCallbackInfo_ReturnValue(void* info) {
	V8_FunctionCallbackInfo* the_info = (V8_FunctionCallbackInfo*)info;
	if (the_info->returnValue == NULL) {
		the_info->returnValue = new V8_ReturnValue(
			the_info->engine, 
			the_info->info->GetReturnValue()
		);;
	}
	return (void*)the_info->returnValue;
}

} // extern "C"