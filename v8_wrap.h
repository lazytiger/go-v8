#ifndef _V8_WARP_H_
#define _V8_WARP_H_

#include <stdint.h>
#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef enum {
        OTP_Context = 0,
        OTP_Getter,
        OTP_Setter,
        OTP_Query,
        OTP_Deleter,
        OTP_Enumerator,
        OTP_Data,
        OTP_Num
} PropertyDataEnum;

typedef enum {
        OTA_Context = 0,
        OTA_Getter,
        OTA_Setter,
        OTA_KeyString,
        OTA_KeyLength,
        OTA_Data,
        OTA_Num
} AccessorDataEnum;

typedef struct {
        void*        engine;
        void*        info;
        void*        setValue;
        const char*  key;
        int          key_length;
        void*        data;
        void*        callback;
        void*        returnValue;
} V8_AccessorCallbackInfo;

typedef struct {
        void*     engine;
        void*     info;
        void*     callback;
        void*     setValue;
        void*     data;
        char*	  key;
	uint32_t  index;
        void*     returnValue;
} V8_PropertyCallbackInfo;

/*
V8
*/
extern const char* V8_GetVersion();

extern void V8_SetFlagsFromString(const char* str, int length);

/*
engine
*/
extern void* V8_NewEngine();

extern void V8_DisposeEngine(void* engine);

extern void* V8_ParseJSON(void* context, const char* json, int json_length);

/*
context
*/
extern void* V8_NewContext(void* engine, void* global_template);

extern void V8_DisposeContext(void* context);

extern void V8_Context_Enter(void* context);

extern void V8_Context_Exit(void* context);

extern void V8_Context_Scope(void* context, void* context_ptr, void* callback);

extern void V8_Context_ThrowException(void* context, const char* err, int err_length);

extern char* V8_Context_TryCatch(void* context, void* callback, int simple);

/*
script
*/
extern void* V8_Compile(void* engine, const char* code, int length, void* script_origin, void* script_data);

extern void V8_DisposeScript(void* script);

extern void* V8_Script_Run(void* script);

/*
script data
*/
extern void* V8_PreCompile(void* engine, const char* code, int length);

extern void* V8_NewScriptData(const char* data, int length);

extern void V8_DisposeScriptData(void* script_data);

extern int V8_ScriptData_Length(void* script_data);

extern const char* V8_ScriptData_Data(void* script_data);

extern int V8_ScriptData_HasError(void* script_data);

/*
script origin
*/
extern void* V8_NewScriptOrigin(void* engine, const char* name, int name_length, int line_offset, int column_offset);

extern void V8_DisposeScriptOrigin(void* script_origin);

/*
value
*/
extern void V8_DisposeValue(void* value);

extern int V8_Value_IsUndefined(void* value);

extern int V8_Value_IsNull(void* value);

extern int V8_Value_IsTrue(void* value);

extern int V8_Value_IsFalse(void* value);

extern int V8_Value_IsString(void* value);

extern int V8_Value_IsFunction(void* value);

extern int V8_Value_IsArray(void* value);

extern int V8_Value_IsObject(void* value);

extern int V8_Value_IsBoolean(void* value);

extern int V8_Value_IsNumber(void* value);

extern int V8_Value_IsExternal(void* value);

extern int V8_Value_IsInt32(void* value);

extern int V8_Value_IsUint32(void* value);

extern int V8_Value_IsDate(void* value);

extern int V8_Value_IsBooleanObject(void* value);

extern int V8_Value_IsNumberObject(void* value);

extern int V8_Value_IsStringObject(void* value);

extern int V8_Value_IsNativeError(void* value);

extern int V8_Value_IsRegExp(void* value);

extern int V8_Value_ToBoolean(void* value);
  
extern double V8_Value_ToNumber(void* value);

extern int64_t V8_Value_ToInteger(void* value);

extern uint32_t V8_Value_ToUint32(void* value);

extern int32_t V8_Value_ToInt32(void* value);

extern char* V8_Value_ToString(void* value);

extern void* V8_Undefined(void* engine);

extern void* V8_Null(void* engine);

extern void* V8_True(void* engine);

extern void* V8_False(void* engine);

extern void* V8_NewNumber(void* context, double val);

extern void* V8_NewString(void* context, const char* val, int val_length);

/*
object
*/
extern void* V8_NewObject(void* context);

extern int V8_Object_SetProperty(void* value, const char* key, int key_length, void* prop_value, int attribs);

extern void* V8_Object_GetProperty(void* value, const char* key, int key_length);

extern int V8_Object_SetElement(void* value, uint32_t index, void* elem_value);

extern void* V8_Object_GetElement(void* value, uint32_t index);

extern int V8_Object_GetPropertyAttributes(void *value, const char* key, int key_length);

extern int V8_Object_ForceSetProperty(void* value, const char* key, int key_length, void* prop_value, int attribs);

extern int V8_Object_HasProperty(void *value, const char* key, int key_length);

extern int V8_Object_DeleteProperty(void *value, const char* key, int key_length);

extern int V8_Object_ForceDeleteProperty(void *value, const char* key, int key_length);

extern int V8_Object_HasElement(void* value, uint32_t index);

extern int V8_Object_DeleteElement(void* value, uint32_t index);

extern void* V8_Object_GetPropertyNames(void *value);

extern void* V8_Object_GetOwnPropertyNames(void *value);

extern void* V8_Object_GetPrototype(void *value);

extern int V8_Object_SetPrototype(void *value, void *proto);

extern void V8_Object_SetAccessor(void *value, const char* key, int key_length, void* getter, void* setter, void* data, int attribs);

extern void* V8_AccessorCallbackInfo_This(void *info, AccessorDataEnum type);

extern void* V8_AccessorCallbackInfo_Holder(void *info, AccessorDataEnum type);

extern void* V8_AccessorCallbackInfo_ReturnValue(void *info, AccessorDataEnum typ);

extern void* V8_PropertyCallbackInfo_This(void *info, PropertyDataEnum typ );

extern void* V8_PropertyCallbackInfo_Holder(void *info, PropertyDataEnum typ );

extern void* V8_PropertyCallbackInfo_ReturnValue(void *info, PropertyDataEnum typ );

/*
array
*/
extern void* V8_NewArray(void* context, int length);

extern int V8_Array_Length(void* value);

/*
regexp
*/
extern void* V8_NewRegExp(void* context, const char* pattern, int length, int flags);

extern char* V8_RegExp_Pattern(void* value);

extern int V8_RegExp_Flags(void* value);

/*
return value
*/
extern void V8_ReturnValue_Set(void* rv, void* result);

extern void V8_ReturnValue_SetBoolean(void* rv, int v);

extern void V8_ReturnValue_SetNumber(void* rv, double v);

extern void V8_ReturnValue_SetInt32(void* rv, int32_t v);

extern void V8_ReturnValue_SetUint32(void* rv, uint32_t v);

extern void V8_ReturnValue_SetString(void* rv, const char* str, int str_length);

extern void V8_ReturnValue_SetNull(void* rv);

extern void V8_ReturnValue_SetUndefined(void* rv);

/*
function
*/
extern void* V8_Function_Call(void* value, int argc, void* argv);

extern void* V8_FunctionCallbackInfo_Get(void* info, int i);

extern int V8_FunctionCallbackInfo_Length(void* info);

extern void* V8_FunctionCallbackInfo_Callee(void* info);

extern void* V8_FunctionCallbackInfo_This(void* info);

extern void* V8_FunctionCallbackInfo_Holder(void* info);

extern void* V8_FunctionCallbackInfo_ReturnValue(void* info);

/*
object template
*/
extern void* V8_NewObjectTemplate(void* engine);

extern void V8_DisposeObjectTemplate(void* tpl);

extern void V8_ObjectTemplate_SetProperty(void* tpl, const char* key, int key_length, void* prop_value, int attribs);

extern void* V8_ObjectTemplate_NewObject(void* tpl);

extern void V8_ObjectTemplate_SetAccessor(void *tpl, const char* key, int key_length, void* getter, void* setter, void* data, int attribs);

extern void V8_ObjectTemplate_SetNamedPropertyHandler(
        void* tpl, 
        void* getter, 
        void* setter, 
        void* query, 
        void* deleter, 
        void* enumerator, 
        void* data
);

extern void V8_ObjectTemplate_SetIndexedPropertyHandler(
        void* tpl, 
        void* getter, 
        void* setter, 
        void* query, 
        void* deleter, 
        void* enumerator, 
        void* data
);

/*
function template
*/
extern void* V8_NewFunctionTemplate(void* engine, void* callback);

extern void V8_DisposeFunctionTemplate(void* tpl);

extern void* V8_FunctionTemplate_GetFunction(void* tpl);

#ifdef __cplusplus
} // extern "C"
#endif

#endif
