#include <stdlib.h>
#include <stdint.h>
#include "v8.h"

extern "C" {

using namespace v8;

class V8_Context {
public:
	V8_Context(Isolate *isolate, Handle<Context> context) {
		self.Reset(isolate, context);
		isolate_ = isolate;
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
	Persistent<Context> self;
};

class V8_Script {
public:
	V8_Script(Isolate *isolate, Handle<Script> script) {
		self.Reset(isolate, script);
		isolate_ = isolate;
	}

	~V8_Script() {
		self.Dispose();
		self.Reset();
	}

	Isolate *GetIsolate() {
		return isolate_;
	}

	Isolate *isolate_;
	Persistent<Script> self;
};

class V8_Value {
public:
	V8_Value(Isolate* isolate, Handle<Value> value) {
		self.Reset(isolate, value);
		isolate_ = isolate;
	}

	~V8_Value() {
		self.Dispose();
		self.Reset();
	}

	Isolate *GetIsolate() {
		return isolate_;
	}

	Isolate *isolate_;
	Persistent<Value> self;
};

void *V8_NewIsolate() {
	return (void *)Isolate::New();
}

void V8_DisposeIsolate(void *isolate) {
	static_cast<Isolate *>(isolate)->Dispose();
}

/*
Context warppers
*/
void *V8_NewContext(void *isolate_ptr) {
	Isolate* isolate = static_cast<Isolate *>(isolate_ptr);
	Locker locker(isolate);
	Isolate::Scope isolate_scope(isolate);
	HandleScope handle_scope(isolate);
	
	Handle<Context> context = Context::New(isolate);

	if (context.IsEmpty())
		return NULL;

	return (void *)(new V8_Context(isolate, context));
}

void V8_DisposeContext(void *context) {
	delete static_cast<V8_Context *>(context);
}

#define ISOLATE_SCOPE(v) \
	Isolate* isolate = v->GetIsolate(); \
	Locker locker(isolate); \
	Isolate::Scope isolate_scope(isolate); \
	HandleScope handle_scope(isolate) \

/*
Script warppers
*/
void *V8_CompileScript(void *context, const char *code) {
	V8_Context *ctx = static_cast<V8_Context *>(context);

	ISOLATE_SCOPE(ctx);

	Local<Context> local_context = Local<Context>::New(isolate, ctx->self);

	Context::Scope context_scope(local_context);

	Handle<String> source = String::New(code);
	Handle<Script> script = Script::Compile(source);

	if (script.IsEmpty())
		return NULL;

	return (void *)(new V8_Script(isolate, script));
}

void V8_DisposeScript(void *script) {
	delete static_cast<V8_Script *>(script);
}

void *V8_RunScript(void *context, void *script) {
	V8_Context *ctx = static_cast<V8_Context *>(context);
	V8_Script *spt = static_cast<V8_Script *>(script);

	ISOLATE_SCOPE(ctx);

	Local<Context> local_context = Local<Context>::New(isolate, ctx->self);
	Local<Script> local_script = Local<Script>::New(isolate, spt->self);

	Context::Scope context_scope(local_context);
	
	Handle<Value> result = local_script->Run();

	if (result.IsEmpty())
		return NULL;

	return (void *)(new V8_Value(isolate, result));
}

/*
Value warppers
*/
void V8_DisposeValue(void *value) {
	delete static_cast<V8_Value *>(value);
}

#define VALUE_TO_LOCAL(value, local_value) \
	V8_Value *val = static_cast<V8_Value *>(value); \
	ISOLATE_SCOPE(val); \
	Local<Value> local_value = Local<Value>::New(isolate, val->self)

char *V8_ValueToString(void *value) {
	VALUE_TO_LOCAL(value, local_value);

	Handle<String> string = local_value->ToString();
	uint8_t *str = (uint8_t *)malloc(string->Length() + 1);
	string->WriteOneByte(str);
	return (char *)str;
}

int V8_ValueIsUndefined(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsUndefined();
}

int V8_ValueIsNull(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsNull();
}

int V8_ValueIsTrue(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsTrue();
}

int V8_ValueIsFalse(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsFalse();
}

int V8_ValueIsString(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsString();
}

int V8_ValueIsFunction(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsFunction();
}

int V8_ValueIsArray(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsArray();
}

int V8_ValueIsObject(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsObject();
}

int V8_ValueIsBoolean(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsBoolean();
}

int V8_ValueIsNumber(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsNumber();
}

int V8_ValueIsExternal(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsExternal();
}

int V8_ValueIsInt32(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsInt32();
}

int V8_ValueIsUint32(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsUint32();
}

int V8_ValueIsDate(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsDate();
}

int V8_ValueIsBooleanObject(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsBooleanObject();
}

int V8_ValueIsNumberObject(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsNumberObject();
}

int V8_ValueIsStringObject(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsStringObject();
}

int V8_ValueIsNativeError(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsNativeError();
}

int V8_ValueIsRegExp(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IsRegExp();
}

int V8_ValueGetBoolean(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->BooleanValue();
}
  
double V8_ValueGetNumber(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->NumberValue();
}

int64_t V8_ValueGetInteger(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->IntegerValue();
}

uint32_t V8_ValueGetUint32(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->Uint32Value();
}

int32_t V8_ValueGetInt32(void *value) {
	VALUE_TO_LOCAL(value, local_value);
	return local_value->Int32Value();
}

} // extern "C"