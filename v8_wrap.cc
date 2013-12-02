#include <cstdlib>
#include <cstring>
#include <sstream>
#include <string>
#include "v8.h"
#include "v8_wrap.h"

extern "C" {

#include "_cgo_export.h"

using namespace v8;

#define ISOLATE_SCOPE(isolate_ptr) \
	Isolate* isolate = isolate_ptr; \
	Locker locker(isolate); \
	Isolate::Scope isolate_scope(isolate) \

#define ENGINE_SCOPE(engine) \
	V8_Context* the_engine = static_cast<V8_Context*>(engine); \
	ISOLATE_SCOPE(the_engine->GetIsolate()) \

#define CONTEXT_SCOPE(engine) \
	V8_Context* the_context = static_cast<V8_Context*>(engine); \
	ISOLATE_SCOPE(the_context->GetIsolate()) \

#define VALUE_SCOPE(value) \
	V8_Value* the_value = static_cast<V8_Value*>(value); \
	ISOLATE_SCOPE(the_value->GetIsolate()); \
	Local<Value> local_value = the_value->self \

#define OBJECT_TEMPLATE_SCOPE(tpl) \
	V8_ObjectTemplate* the_template = static_cast<V8_ObjectTemplate*>(tpl); \
	ISOLATE_SCOPE(the_template->GetIsolate()); \
	Local<ObjectTemplate> local_template = Local<ObjectTemplate>::New(isolate, the_template->self) \

#define OBJECT_TEMPLATE_HANDLE_SCOPE(tpl) \
	V8_ObjectTemplate* the_template = static_cast<V8_ObjectTemplate*>(tpl); \
	ISOLATE_SCOPE(the_template->GetIsolate()); \
	HandleScope scope(isolate); \
	Local<ObjectTemplate> local_template = Local<ObjectTemplate>::New(isolate, the_template->self) \

#define FUNCTION_TEMPLATE_SCOPE(tpl) \
	V8_FunctionTemplate* the_template = static_cast<V8_FunctionTemplate*>(tpl); \
	ISOLATE_SCOPE(the_template->GetIsolate()); \
	Local<FunctionTemplate> local_template = Local<FunctionTemplate>::New(isolate, the_template->self) \

#define FUNCTION_TEMPLATE_HANDLE_SCOPE(tpl) \
	V8_FunctionTemplate* the_template = static_cast<V8_FunctionTemplate*>(tpl); \
	ISOLATE_SCOPE(the_template->GetIsolate()); \
	HandleScope scope(isolate); \
	Local<FunctionTemplate> local_template = Local<FunctionTemplate>::New(isolate, the_template->self) \

class V8_Context {
public:
	V8_Context(Isolate* isolate, Handle<Context> context) {
		isolate_ = isolate;
		self.Reset(isolate_, context);
	}

	V8_Context(V8_Context* engine, Handle<Context> context) {
		isolate_ = engine->GetIsolate();
		self.Reset(isolate_, context);
	}

	~V8_Context() {
		Locker locker(isolate_);
		Isolate::Scope isolate_scope(isolate_);

		self.Reset();
	}

	Isolate* GetIsolate() {
		return isolate_;
	}

	Isolate* isolate_;
	Persistent<Context> self;
};

class V8_Script {
public:
	V8_Script(V8_Context* the_engine, Handle<Script> script) {
		engine = the_engine;
		self.Reset(engine->GetIsolate(), script);
	}

	~V8_Script() {
		ISOLATE_SCOPE(GetIsolate());
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
	V8_Value(V8_Context* the_context, Handle<Value> value) : self(value) {
		context = the_context;
		//self.Reset(context->GetIsolate(), value);
		//context_handler.Reset(context->GetIsolate(), context->self);
	}

	~V8_Value() {
		//ISOLATE_SCOPE(GetIsolate());
		//Context::Scope context_scope(Local<Context>::New(isolate, context->self));
		//self.Reset();
		//context_handler.Reset();
	}

	Isolate* GetIsolate() {
		return context->GetIsolate();
	}

	V8_Context* context;
	Handle<Value> self;
	//Persistent<Context> context_handler;
};

typedef struct V8_ReturnValue {
	V8_ReturnValue(V8_Context* the_context, ReturnValue<Value> the_value) :
		context(the_context),
		value(the_value) {
	}

	V8_Context*        context;
	ReturnValue<Value> value;
} V8_ReturnValue;

class V8_ObjectTemplate {
public:
	V8_ObjectTemplate(V8_Context* the_engine, Handle<ObjectTemplate> the_template) {
		engine = the_engine;
		self.Reset(engine->GetIsolate(), the_template);
	}

	~V8_ObjectTemplate() {
		Locker locker(GetIsolate());
		Isolate::Scope isolate_scope(GetIsolate());

		self.Reset();
	}

	Isolate* GetIsolate() {
		return engine->GetIsolate();
	}

	V8_Context* engine;
	Persistent<ObjectTemplate> self;
};

class V8_FunctionTemplate {
public:
	V8_FunctionTemplate(V8_Context* the_engine, Handle<FunctionTemplate> the_template) {
		engine = the_engine;
		self.Reset(engine->GetIsolate(), the_template);
	}

	~V8_FunctionTemplate() {
		Locker locker(GetIsolate());
		Isolate::Scope isolate_scope(GetIsolate());

		self.Reset();
	}

	Isolate* GetIsolate() {
		return engine->GetIsolate();
	}

	V8_Context* engine;
	Persistent<FunctionTemplate> self;
};

void* new_V8_Value(V8_Context* context, Handle<Value> value) {
	if (value.IsEmpty())
		return NULL;
	return (void*)new V8_Value(context, value);
}

/*
engine
*/
void* V8_NewEngine() {
	ISOLATE_SCOPE(Isolate::New());

	HandleScope handle_scope(isolate);
	Handle<Context> context = Context::New(isolate);

	if (context.IsEmpty())
		return NULL;

	context->Enter();

	return (void*)(new V8_Context(isolate, context));
}

void V8_DisposeEngine(void* engine) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);
	ISOLATE_SCOPE(the_engine->GetIsolate());
	Local<Context> local_context = Local<Context>::New(isolate, the_engine->self);

	local_context->Exit();

	delete the_engine;

	isolate->Dispose();
}

void* V8_ParseJSON(void* context, const char* json, int json_length) {
	CONTEXT_SCOPE(context);

	Handle<Value> value = JSON::Parse(
		String::NewFromOneByte(isolate, (uint8_t*)json, String::kNormalString, json_length)
	);

	if (value.IsEmpty())
		return NULL;

	return new_V8_Value(the_context, value);
}

/*
context
*/
void* V8_NewContext(void* engine, void* global_template) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);

	ISOLATE_SCOPE(the_engine->GetIsolate());

	HandleScope handle_scope(isolate);
	Handle<Context> context = Context::New(
		isolate, NULL,
		global_template == NULL ? Handle<ObjectTemplate>() : Local<ObjectTemplate>::New(
			isolate, static_cast<V8_ObjectTemplate*>(global_template)->self
		)
	);

	if (context.IsEmpty())
		return NULL;

	return (void*)(new V8_Context(the_engine, context));
}

void V8_DisposeContext(void* context) {
	delete static_cast<V8_Context*>(context);
}

typedef struct scope_data {
	void* context;
	void* context_ptr;
} scope_data;

void V8_Context_Scope(void* context, void* context_ptr, void* callback) {
	V8_Context* ctx = static_cast<V8_Context*>(context);
	ISOLATE_SCOPE(ctx->GetIsolate());

	void* prev_context = isolate->GetData();
	scope_data data;
	data.context = context;
	data.context_ptr = context_ptr;
	isolate->SetData(&data);

	// Make nested context scropt use the outermost HandleScope
	if (prev_context == NULL) {
		HandleScope handle_scope(isolate);

		Context::Scope scope(Local<Context>::New(isolate, ctx->self));
		context_scope_callback(context_ptr, callback);
	} else {
		Context::Scope scope(Local<Context>::New(isolate, ctx->self));
		context_scope_callback(context_ptr, callback);
	}

	isolate->SetData(prev_context);
}

V8_Context* V8_Current_Context(Isolate* isolate) {
	void* data = isolate->GetData();
	if (data == NULL)
		v8_panic((char*)"Please call this API in a context scope");
	return static_cast<V8_Context*>(static_cast<scope_data*>(isolate->GetData())->context);
}

void* V8_Current_ContextPtr(Isolate* isolate) {
	void* data = isolate->GetData();
	if (data == NULL)
		v8_panic((char*)"Please call this API in a context scope");
	return static_cast<scope_data*>(isolate->GetData())->context_ptr;
}

void* V8_Context_Global(void* context) {
	CONTEXT_SCOPE(context);
	Local<Context> local_context = Local<Context>::New(isolate, the_context->self);
	return new_V8_Value(the_context, local_context->Global()); 
}

// Extracts a C string from a V8 Utf8Value.
const char* ToCString(const v8::String::Utf8Value& value) {
  return *value ? *value : "<string conversion failed>";
}

void V8_Context_ThrowException(void* context, const char* err, int err_length) {
	V8_Context* ctx = static_cast<V8_Context*>(context);
	ISOLATE_SCOPE(ctx->GetIsolate());

	isolate->ThrowException(
		String::NewFromOneByte(isolate, (uint8_t*)err, String::kNormalString, err_length)
	);
}

char* V8_Context_TryCatch(void* context, void* callback, int simple) {
	V8_Context* ctx = static_cast<V8_Context*>(context);
	ISOLATE_SCOPE(ctx->GetIsolate());

	v8::TryCatch try_catch;

	try_catch_callback(callback);

	if (!try_catch.HasCaught()) {
		return NULL;
	}

	v8::String::Utf8Value exception(try_catch.Exception());
	const char* exception_string = ToCString(exception);
	v8::Handle<v8::Message> message = try_catch.Message();

	if (message.IsEmpty() || simple) {
		// V8 didn't provide any extra information about this error; just
		// print the exception.
		char *cstr = (char*)malloc(exception.length() + 1);
		std::strcpy(cstr, exception_string);
		return cstr;
	}

	std::stringstream report;

	// Print (filename):(line number): (message).
	v8::String::Utf8Value filename(message->GetScriptResourceName());
	const char* filename_string = ToCString(filename);
	int linenum = message->GetLineNumber();
	report << filename_string << ":" << linenum << ": " << exception_string << std::endl;

	// Print line of source code.
	v8::String::Utf8Value sourceline(message->GetSourceLine());
	const char* sourceline_string = ToCString(sourceline);
	report << sourceline_string << std::endl;

	// Print wavy underline (GetUnderline is deprecated).
	int start = message->GetStartColumn();
	for (int i = 0; i < start; i++) {
		report << " ";
	}

	int end = message->GetEndColumn();
	for (int i = start; i < end; i++) {
		report << "^";
	}

	report << std::endl;

	v8::String::Utf8Value stack_trace(try_catch.StackTrace());
	if (stack_trace.length() > 0) {
		const char* stack_trace_string = ToCString(stack_trace);
		report << stack_trace_string << std::endl;
	}

	std::string report_string = report.str();
	char *cstr = (char*)malloc(report_string.length() +1);
	std::strcpy(cstr, report_string.c_str());

	return cstr;
}

/*
script
*/
void* V8_Compile(void* engine, const char* code, int length, void* script_origin,void* script_data) {
	ENGINE_SCOPE(engine);

	HandleScope handle_scope(isolate);

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

void* V8_Script_Run(void* script) {
	V8_Script* the_script = static_cast<V8_Script*>(script);
	ISOLATE_SCOPE(the_script->engine->GetIsolate());
	V8_Context* the_context = V8_Current_Context(isolate);
	Local<Script> local_script = Local<Script>::New(isolate, the_script->self);

	return new_V8_Value(the_context, local_script->Run());
}

/*
script data
*/
void* V8_PreCompile(void* engine, const char* code, int length) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);
	ISOLATE_SCOPE(the_engine->GetIsolate());
	HandleScope handle_scope(isolate);

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
	HandleScope handle_scope(isolate);

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
	return new_V8_Value(the_engine, Undefined(isolate));
}

void* V8_Null(void* engine) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);
	ISOLATE_SCOPE(the_engine->GetIsolate());
	return new_V8_Value(the_engine, Null(isolate));
}

void* V8_True(void* engine) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);
	ISOLATE_SCOPE(the_engine->GetIsolate());
	return new_V8_Value(the_engine, True(isolate));
}

void* V8_False(void* engine) {
	V8_Context* the_engine = static_cast<V8_Context*>(engine);
	ISOLATE_SCOPE(the_engine->GetIsolate());
	return new_V8_Value(the_engine, False(isolate));
}

void* V8_NewNumber(void* context, double val) {
	CONTEXT_SCOPE(context);

	return new_V8_Value(the_context, Number::New(isolate, val));
}

void* V8_NewString(void* context, const char* val, int val_length) {
	CONTEXT_SCOPE(context);

	return new_V8_Value(the_context,
		String::NewFromOneByte(isolate, (uint8_t*)val, String::kNormalString, val_length)
	);
}

/*
object
*/
void* V8_NewObject(void* context) {
	CONTEXT_SCOPE(context);

	return new_V8_Value(the_context, Object::New());
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

	return new_V8_Value(the_value->context,
		Local<Object>::Cast(local_value)->Get(
			String::NewFromOneByte(isolate, (uint8_t*)key, String::kNormalString, key_length)
		)
	);
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

	return new_V8_Value(the_value->context,
		Local<Object>::Cast(local_value)->Get(index)
	);
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

	return new_V8_Value(the_value->context,
		Local<Object>::Cast(local_value)->GetPropertyNames()
	);
}

void* V8_Object_GetOwnPropertyNames(void* value) {
	VALUE_SCOPE(value);

	return new_V8_Value(the_value->context,
		Local<Object>::Cast(local_value)->GetOwnPropertyNames()
	);
}

void* V8_Object_GetPrototype(void* value) {
	VALUE_SCOPE(value);

	return new_V8_Value(the_value->context,
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

void V8_AccessorGetterCallback(Local<String> property, const PropertyCallbackInfo<Value>& info) {
	v8::Isolate* isolate_ptr = info.GetIsolate();
	ISOLATE_SCOPE(isolate_ptr);

	Local<Array> callback_data = Local<Array>::Cast(info.Data());

	V8_AccessorCallbackInfo callback_info;
	callback_info.engine = Local<External>::Cast(callback_data->Get(OTA_Context))->Value();
	callback_info.info = (void*)&info;
	callback_info.returnValue = NULL;
	callback_info.data = Local<External>::Cast(callback_data->Get(OTA_Data))->Value();
	callback_info.callback = Local<External>::Cast(callback_data->Get(OTA_Getter))->Value();
	callback_info.key = (const char*)Local<External>::Cast(callback_data->Get(OTA_KeyString))->Value();
	callback_info.key_length = Local<Integer>::Cast(callback_data->Get(OTA_KeyLength))->Value();

	void* context_ptr = V8_Current_ContextPtr(isolate);

	go_accessor_callback(OTA_Getter, &callback_info, context_ptr);

	if (callback_info.returnValue != NULL)
		delete static_cast<V8_ReturnValue*>(callback_info.returnValue);
}

void V8_AccessorSetterCallback(Local<String> property, Local<Value> value, const PropertyCallbackInfo<void>& info) {
	v8::Isolate* isolate_ptr = info.GetIsolate();
	ISOLATE_SCOPE(isolate_ptr);

	Local<Array> callback_data = Local<Array>::Cast(info.Data());

	V8_AccessorCallbackInfo callback_info;
	callback_info.engine = Local<External>::Cast(callback_data->Get(OTA_Context))->Value();
	callback_info.info = (void*)&info;
	callback_info.returnValue = NULL;
	callback_info.setValue = new_V8_Value(static_cast<V8_Context*>(callback_info.engine), value);
	callback_info.data = Local<External>::Cast(callback_data->Get(OTA_Data))->Value();
	callback_info.callback = Local<External>::Cast(callback_data->Get(OTA_Setter))->Value();
	callback_info.key = (const char*)Local<External>::Cast(callback_data->Get(OTA_KeyString))->Value();
	callback_info.key_length = Local<Integer>::Cast(callback_data->Get(OTA_KeyLength))->Value();

	void* context_ptr = V8_Current_ContextPtr(isolate);

	go_accessor_callback(OTA_Setter, &callback_info, context_ptr);

	if (callback_info.returnValue != NULL)
		delete static_cast<V8_ReturnValue*>(callback_info.returnValue);
}

// sync with V8_ObjectTemplate_SetAccessor
void V8_Object_SetAccessor(void *value, const char* key, int key_length, void* getter, void* setter, void* data, int attribs) {
	VALUE_SCOPE(value);

	Handle<Array> callback_info = Array::New(OTA_Num);
	callback_info->Set(OTA_Context, External::New((void*)the_value->context));
	callback_info->Set(OTA_Getter, External::New(getter));
	callback_info->Set(OTA_Setter, External::New(setter));
	callback_info->Set(OTA_KeyString, External::New((void*)key));
	callback_info->Set(OTA_KeyLength, Integer::New(key_length));
	callback_info->Set(OTA_Data, External::New(data));

	if (callback_info.IsEmpty())
		return;

	Local<Object>::Cast(local_value)->SetAccessor(
		String::NewFromOneByte(isolate, (uint8_t*)key, String::kNormalString, key_length),
		V8_AccessorGetterCallback, setter == NULL ? NULL : V8_AccessorSetterCallback,
 		callback_info
	);
}

void* V8_PropertyCallbackInfo_This(void *info, PropertyDataEnum typ ) {
	V8_PropertyCallbackInfo* the_info = (V8_PropertyCallbackInfo*)info;
	V8_Context* engine = static_cast<V8_Context*>(the_info->engine);
	ENGINE_SCOPE(engine);
	switch(typ) {
	case OTP_Getter:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<Value>*>(the_info->info)->This());
	case OTP_Setter:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<Value>*>(the_info->info)->This());
	case OTP_Deleter:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<Boolean>*>(the_info->info)->This());
	case OTP_Query:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<Integer>*>(the_info->info)->This());
	case OTP_Enumerator:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<Array>*>(the_info->info)->This());
	default:
		//impossible, should panic
		return NULL;
	}
}

void* V8_PropertyCallbackInfo_Holder(void *info, PropertyDataEnum typ) {
	V8_PropertyCallbackInfo* the_info = (V8_PropertyCallbackInfo*)info;
	V8_Context* engine = static_cast<V8_Context*>(the_info->engine);
	ENGINE_SCOPE(engine);
	switch(typ) {
	case OTP_Getter:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<Value>*>(the_info->info)->Holder());
	case OTP_Setter:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<Value>*>(the_info->info)->Holder());
	case OTP_Deleter:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<Boolean>*>(the_info->info)->Holder());
	case OTP_Query:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<Integer>*>(the_info->info)->Holder());
	case OTP_Enumerator:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<Array>*>(the_info->info)->Holder());
	default:
		//impossible, should panic
		return NULL;
	}
}

void* V8_PropertyCallbackInfo_ReturnValue(void *info,  PropertyDataEnum typ) {
	V8_PropertyCallbackInfo* the_info = (V8_PropertyCallbackInfo*)info;
	V8_Context* engine = static_cast<V8_Context*>(the_info->engine);
	ENGINE_SCOPE(engine);
	if (the_info->returnValue == NULL) {
		switch(typ) {
		case OTP_Getter:
			the_info->returnValue = new V8_ReturnValue(
				engine,
				static_cast<PropertyCallbackInfo<Value>*>(the_info->info)->GetReturnValue()
			);
			break;
		case OTP_Setter:
			the_info->returnValue = new V8_ReturnValue(
				engine,
				static_cast<PropertyCallbackInfo<Value>*>(the_info->info)->GetReturnValue()
			);
			break;
		case OTP_Deleter:
			the_info->returnValue = new V8_ReturnValue(
				engine,
				ReturnValue<Value>(static_cast<PropertyCallbackInfo<Boolean>*>(the_info->info)->GetReturnValue())
			);
			break;
		case OTP_Query:
			the_info->returnValue = new V8_ReturnValue(
				engine,
				ReturnValue<Value>(static_cast<PropertyCallbackInfo<Integer>*>(the_info->info)->GetReturnValue())
			);
			break;
		case OTP_Enumerator:
			the_info->returnValue = new V8_ReturnValue(
				engine,
				ReturnValue<Value>(static_cast<PropertyCallbackInfo<Array>*>(the_info->info)->GetReturnValue())
			);
			break;
		}
	}
	return the_info->returnValue;
}

void* V8_AccessorCallbackInfo_This(void *info, AccessorDataEnum typ) {
	V8_AccessorCallbackInfo* the_info = (V8_AccessorCallbackInfo*)info;
	V8_Context* engine = static_cast<V8_Context*>(the_info->engine);
	ENGINE_SCOPE(engine);
	switch(typ) {
	case OTA_Getter:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<Value>*>(the_info->info)->This());
	case OTA_Setter:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<void>*>(the_info->info)->This());
	default:
		return NULL;
	}
}

void* V8_AccessorCallbackInfo_Holder(void *info, AccessorDataEnum typ) {
	V8_AccessorCallbackInfo* the_info = (V8_AccessorCallbackInfo*)info;
	V8_Context* engine = static_cast<V8_Context*>(the_info->engine);
	ENGINE_SCOPE(engine);
	switch(typ) {
	case OTA_Getter:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<Value>*>(the_info->info)->Holder());
	case OTA_Setter:
		return new_V8_Value(engine, static_cast<PropertyCallbackInfo<void>*>(the_info->info)->Holder());
	default:
		return NULL;
	}
}

void* V8_AccessorCallbackInfo_ReturnValue(void *info, AccessorDataEnum typ) {
	V8_AccessorCallbackInfo* the_info = (V8_AccessorCallbackInfo*)info;
	V8_Context* engine = static_cast<V8_Context*>(the_info->engine);
	if (the_info->returnValue == NULL) {
		switch(typ) {
		case OTA_Getter:
			the_info->returnValue = new V8_ReturnValue(
				engine,
				static_cast<PropertyCallbackInfo<Value>*>(the_info->info)->GetReturnValue()
			);
		}
	}
	return the_info->returnValue;
}

/*
array
*/
void* V8_NewArray(void* context, int length) {
	CONTEXT_SCOPE(context);
	return new_V8_Value(the_context, Array::New(length));
}

int V8_Array_Length(void* value) {
	VALUE_SCOPE(value);
	return Local<Array>::Cast(local_value)->Length();
}

/*
regexp
*/
void* V8_NewRegExp(void* context, const char* pattern, int length, int flags) {
	CONTEXT_SCOPE(context);

	return new_V8_Value(the_context, RegExp::New(
		String::NewFromOneByte(isolate, (uint8_t*)pattern, String::kNormalString, length),
		(RegExp::Flags)flags
	));
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
	ENGINE_SCOPE(the_rv->context);
	the_rv->value.Set(static_cast<V8_Value*>(result)->self);
}

void V8_ReturnValue_SetBoolean(void* rv, int v) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->context);
	the_rv->value.Set((bool)v);
}

void V8_ReturnValue_SetNumber(void* rv, double v) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->context);
	the_rv->value.Set(v);
}

void V8_ReturnValue_SetInt32(void* rv, int32_t v) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->context);
	the_rv->value.Set(v);
}

void V8_ReturnValue_SetUint32(void* rv, uint32_t v) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->context);
	the_rv->value.Set(v);
}

void V8_ReturnValue_SetString(void* rv, const char* str, int str_length) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->context);
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
	ENGINE_SCOPE(the_rv->context);
	the_rv->value.SetNull();
}

void V8_ReturnValue_SetUndefined(void* rv) {
	V8_ReturnValue* the_rv = (V8_ReturnValue*)rv;
	ENGINE_SCOPE(the_rv->context);
	the_rv->value.SetUndefined();
}

/*
function
*/
typedef struct {
	V8_Context*                            engine;
	const v8::FunctionCallbackInfo<Value>* info;
	V8_ReturnValue*                        returnValue;
} V8_FunctionCallbackInfo;

void V8_FunctionCallback(const v8::FunctionCallbackInfo<Value>& info) {
	v8::Isolate* isolate_ptr = info.GetIsolate();
	ISOLATE_SCOPE(isolate_ptr);

	Local<Array> callback_data = Local<Array>::Cast(info.Data());

	V8_FunctionCallbackInfo callback_info;
	callback_info.engine = (V8_Context*)Local<External>::Cast(callback_data->Get(0))->Value();
	callback_info.info = &info;
	callback_info.returnValue = NULL;

	void* callback = Local<External>::Cast(callback_data->Get(1))->Value();

	void* context_ptr = V8_Current_ContextPtr(isolate);

	go_function_callback(&callback_info, callback, context_ptr);

	if (callback_info.returnValue != NULL)
		delete callback_info.returnValue;
}

void* V8_Function_Call(void* value, int argc, void* argv) {
	VALUE_SCOPE(value);

	Handle<Value>* real_argv = new Handle<Value>[argc];
	V8_Value* *argv_ptr = (V8_Value**)argv;

	for (int i = 0; i < argc; i ++) {
		real_argv[i] = Local<Value>::New(isolate, static_cast<V8_Value*>(argv_ptr[i])->self);
	}

	void* result = new_V8_Value(the_value->context,
		Local<Function>::Cast(local_value)->Call(local_value, argc, real_argv)
	);

	delete[] real_argv;

	return result;
}

void* V8_FunctionCallbackInfo_Get(void* info, int i) {
	V8_FunctionCallbackInfo* the_info = (V8_FunctionCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return new_V8_Value(the_info->engine, (*(the_info->info))[i]);
}

int V8_FunctionCallbackInfo_Length(void* info) {
	V8_FunctionCallbackInfo* the_info = (V8_FunctionCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return the_info->info->Length();
}

void* V8_FunctionCallbackInfo_Callee(void* info) {
	V8_FunctionCallbackInfo* the_info = (V8_FunctionCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return new_V8_Value(the_info->engine, the_info->info->Callee());
}

void* V8_FunctionCallbackInfo_This(void* info) {
	V8_FunctionCallbackInfo* the_info = (V8_FunctionCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return new_V8_Value(the_info->engine, the_info->info->This());
}

void* V8_FunctionCallbackInfo_Holder(void* info) {
	V8_FunctionCallbackInfo* the_info = (V8_FunctionCallbackInfo*)info;
	ENGINE_SCOPE(the_info->engine);
	return new_V8_Value(the_info->engine, the_info->info->Holder());
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

/*
object template
*/
void* V8_NewObjectTemplate(void* engine) {
	ENGINE_SCOPE(engine);
	HandleScope handle_scope(isolate);
	Handle<ObjectTemplate> tpl = ObjectTemplate::New();
	if (tpl.IsEmpty())
		return NULL;
	return (void*)new V8_ObjectTemplate(the_engine, tpl);
}

void V8_DisposeObjectTemplate(void* tpl) {
	delete static_cast<V8_ObjectTemplate*>(tpl);
}

void V8_ObjectTemplate_SetProperty(void* tpl, const char* key, int key_length, void* prop_value, int attribs) {
	OBJECT_TEMPLATE_SCOPE(tpl);

	local_template->Set(
		String::NewFromOneByte(isolate, (uint8_t*)key, String::kNormalString, key_length),
		Local<Value>::New(isolate, static_cast<V8_Value*>(prop_value)->self),
		(v8::PropertyAttribute)attribs
	);
}

void* V8_ObjectTemplate_NewObject(void* tpl) {
	OBJECT_TEMPLATE_SCOPE(tpl);
	V8_Context* the_context = V8_Current_Context(isolate);
	return new_V8_Value(the_context, local_template->NewInstance());
}

// sync with V8_Object_SetAccessor
void V8_ObjectTemplate_SetAccessor(void *tpl, const char* key, int key_length, void* getter, void* setter, void* data, int attribs) {
	OBJECT_TEMPLATE_HANDLE_SCOPE(tpl);

	Handle<Array> callback_info = Array::New(5);
	callback_info->Set(0, External::New((void*)the_template->engine));
	callback_info->Set(1, External::New(getter));
	callback_info->Set(2, External::New(setter));
	callback_info->Set(3, External::New((void*)key));
	callback_info->Set(4, Integer::New(key_length));
	callback_info->Set(5, External::New(data));

	if (callback_info.IsEmpty())
		return;

	local_template->SetAccessor(
		String::NewFromOneByte(isolate, (uint8_t*)key, String::kNormalString, key_length),
		V8_AccessorGetterCallback, setter == NULL ? NULL : V8_AccessorSetterCallback,
 		callback_info
	);
}

void V8_NamedPropertyGetterCallbackBase(
	PropertyDataEnum typ, 
	Local<String> property, 
	Local<Value> value, 
	void* info_ptr,
	v8::Isolate* isolate_ptr, 
	Local<Value> callback_data_val
) {
    ISOLATE_SCOPE(isolate_ptr);
    Local<Array> callback_data = Local<Array>::Cast(callback_data_val);
    V8_PropertyCallbackInfo callback_info;
    callback_info.engine = Local<External>::Cast(callback_data->Get(OTP_Context))->Value();
    callback_info.info = info_ptr;
    callback_info.returnValue = NULL;
    callback_info.data = Local<External>::Cast(callback_data->Get(OTP_Data))->Value();
    callback_info.callback = Local<External>::Cast(callback_data->Get(typ))->Value();
    callback_info.key = NULL;

	if (typ != OTP_Enumerator) {
		uint8_t* key = (uint8_t*)malloc(property->Length() + 1);
		property->WriteOneByte(key);
		callback_info.key = (char*)key;
	}

	if (typ == OTP_Setter) {
		callback_info.setValue = new_V8_Value(
			V8_Current_Context(isolate_ptr),
			value
		);
	}

	void* context_ptr = V8_Current_ContextPtr(isolate);

	go_named_property_callback(typ, &callback_info, context_ptr);

	if (typ != OTP_Enumerator) {
		free(callback_info.key);
	}

	if (callback_info.returnValue != NULL)
		delete static_cast<V8_ReturnValue*>(callback_info.returnValue);
}

void V8_NamedPropertyGetterCallback(Local<String> property, const PropertyCallbackInfo<Value> &info) {
	V8_NamedPropertyGetterCallbackBase(OTP_Getter, property, Local<Value>(), (void*)&info, info.GetIsolate(), info.Data());
}

void V8_NamedPropertySetterCallback(Local<String> property, Local<Value> value, const PropertyCallbackInfo<Value> &info) {
	V8_NamedPropertyGetterCallbackBase(OTP_Setter, property, value, (void*)&info, info.GetIsolate(), info.Data());
}

void V8_NamedPropertyDeleterCallback(Local<String> property, const PropertyCallbackInfo<Boolean> &info) {
	V8_NamedPropertyGetterCallbackBase(OTP_Deleter, property, Local<Value>(), (void*)&info, info.GetIsolate(), info.Data());
}

void V8_NamedPropertyQueryCallback(Local<String> property, const PropertyCallbackInfo<Integer> &info) {
	V8_NamedPropertyGetterCallbackBase(OTP_Query, property, Local<Value>(), (void*)&info, info.GetIsolate(), info.Data());
}

void V8_NamedPropertyEnumeratorCallback(const PropertyCallbackInfo<Array> &info) {
	V8_NamedPropertyGetterCallbackBase(OTP_Enumerator, Local<String>(), Local<Value>(), (void*)&info, info.GetIsolate(), info.Data());
}

void V8_ObjectTemplate_SetNamedPropertyHandler(
	void* tpl, 
	void* getter, 
	void* setter, 
	void* query, 
	void* deleter, 
	void* enumerator, 
	void* data
) {
	OBJECT_TEMPLATE_HANDLE_SCOPE(tpl);

	Handle<Array> callback_info = Array::New(OTP_Num);
	callback_info->Set(OTP_Context, External::New((void*)the_template->engine));
	callback_info->Set(OTP_Getter, External::New(getter));
	callback_info->Set(OTP_Setter, External::New(setter));
	callback_info->Set(OTP_Query, External::New(query));
	callback_info->Set(OTP_Deleter, External::New(deleter));
	callback_info->Set(OTP_Enumerator, External::New(enumerator));
	callback_info->Set(OTP_Data, External::New(data));

	if (callback_info.IsEmpty())
		return;

	local_template->SetNamedPropertyHandler(
		V8_NamedPropertyGetterCallback, 
		setter == NULL ? NULL : V8_NamedPropertySetterCallback,
		query == NULL ? NULL : V8_NamedPropertyQueryCallback,
		deleter == NULL ? NULL : V8_NamedPropertyDeleterCallback,
		enumerator == NULL ? NULL : V8_NamedPropertyEnumeratorCallback,
 		callback_info
	);
}

void V8_IndexedPropertyGetterCallbackBase(
	PropertyDataEnum typ, 
	uint32_t index, 
	Local<Value> value, 
	void* info_ptr,
	v8::Isolate* isolate_ptr, 
	Local<Value> callback_data_val
) {
    ISOLATE_SCOPE(isolate_ptr);

    Local<Array> callback_data = Local<Array>::Cast(callback_data_val);
    V8_PropertyCallbackInfo callback_info;
    callback_info.engine = Local<External>::Cast(callback_data->Get(OTP_Context))->Value();
    callback_info.info = info_ptr;
    callback_info.returnValue = NULL;
    callback_info.data = Local<External>::Cast(callback_data->Get(OTP_Data))->Value();
    callback_info.callback = Local<External>::Cast(callback_data->Get(typ))->Value();
	callback_info.index = index;

	if (typ == OTP_Setter) {
		callback_info.setValue = new_V8_Value(
			V8_Current_Context(isolate_ptr),
			value
		);
	}

	void* context_ptr = V8_Current_ContextPtr(isolate);

	go_indexed_property_callback(typ, &callback_info, context_ptr);

	if (callback_info.returnValue != NULL)
		delete static_cast<V8_ReturnValue*>(callback_info.returnValue);
}

void V8_IndexedPropertyGetterCallback(uint32_t index, const PropertyCallbackInfo<Value> &info) {
	V8_IndexedPropertyGetterCallbackBase(OTP_Getter, index, Local<Value>(), (void*)&info, info.GetIsolate(), info.Data());
}

void V8_IndexedPropertySetterCallback(uint32_t index, Local<Value> value, const PropertyCallbackInfo<Value> &info) {
	V8_IndexedPropertyGetterCallbackBase(OTP_Setter, index, value, (void*)&info, info.GetIsolate(), info.Data());
}

void V8_IndexedPropertyDeleterCallback(uint32_t index, const PropertyCallbackInfo<Boolean> &info) {
	V8_IndexedPropertyGetterCallbackBase(OTP_Getter, index, Local<Value>(), (void*)&info, info.GetIsolate(), info.Data());
}

void V8_IndexedPropertyQueryCallback(uint32_t index, const PropertyCallbackInfo<Integer> &info) {
	V8_IndexedPropertyGetterCallbackBase(OTP_Getter, index, Local<Value>(), (void*)&info, info.GetIsolate(), info.Data());
}

void V8_IndexedPropertyEnumeratorCallback(const PropertyCallbackInfo<Array> &info) {
	V8_IndexedPropertyGetterCallbackBase(OTP_Getter, 0, Local<Value>(), (void*)&info, info.GetIsolate(), info.Data());
}

void V8_ObjectTemplate_SetIndexedPropertyHandler(
	void* tpl, 
	void* getter, 
	void* setter, 
	void* query, 
	void* deleter, 
	void* enumerator, 
	void* data
) {
	OBJECT_TEMPLATE_HANDLE_SCOPE(tpl);
	
	Handle<Array> callback_info = Array::New(OTP_Num);
	callback_info->Set(OTP_Context, External::New((void*)the_template->engine));
	callback_info->Set(OTP_Getter, External::New(getter));
	callback_info->Set(OTP_Setter, External::New(setter));
	callback_info->Set(OTP_Query, External::New(query));
	callback_info->Set(OTP_Deleter, External::New(deleter));
	callback_info->Set(OTP_Enumerator, External::New(enumerator));
	callback_info->Set(OTP_Data, External::New(data));

	if (callback_info.IsEmpty())
		return;

	local_template->SetIndexedPropertyHandler(
		V8_IndexedPropertyGetterCallback, 
		setter == NULL ? NULL : V8_IndexedPropertySetterCallback,
		query == NULL ? NULL : V8_IndexedPropertyQueryCallback,
		deleter == NULL ? NULL : V8_IndexedPropertyDeleterCallback,
		enumerator == NULL ? NULL : V8_IndexedPropertyEnumeratorCallback,
 		callback_info
	);
}


/*
function template
*/
void* V8_NewFunctionTemplate(void* engine, void* callback) {
	ENGINE_SCOPE(engine);

	HandleScope scope(isolate);

	Handle<Array> callback_data = Array::New(2);

	if (callback_data.IsEmpty())
		return NULL;

	callback_data->Set(0, External::New(engine));
	callback_data->Set(1, External::New(callback));

	Handle<FunctionTemplate> tpl = callback == NULL ? FunctionTemplate::New() : FunctionTemplate::New(
		V8_FunctionCallback, callback_data
	);

	if (tpl.IsEmpty())
		return NULL;

	return (void*)new V8_FunctionTemplate(the_engine, tpl);
}

void V8_DisposeFunctionTemplate(void* tpl) {
	delete static_cast<V8_FunctionTemplate*>(tpl);
}

void* V8_FunctionTemplate_GetFunction(void* tpl) {
	FUNCTION_TEMPLATE_SCOPE(tpl);
	V8_Context* the_context = V8_Current_Context(isolate);
	return new_V8_Value(the_context, local_template->GetFunction());
}

void V8_FunctionTemplate_SetClassName(void* tpl, const char* name, int name_length) {
	FUNCTION_TEMPLATE_HANDLE_SCOPE(tpl);
	return local_template->SetClassName(
		String::NewFromOneByte(isolate, (uint8_t*)name, String::kNormalString, name_length)
	);
}

void* V8_FunctionTemplate_InstanceTemplate(void* tpl) {
	FUNCTION_TEMPLATE_HANDLE_SCOPE(tpl);
	return new V8_ObjectTemplate(the_template->engine, local_template->InstanceTemplate());
}

const char* V8_GetVersion() {
	return V8::GetVersion();
}

void V8_SetFlagsFromString(const char* str, int length) {
	V8::SetFlagsFromString(str, length);
}

} // extern "C"
