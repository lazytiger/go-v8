#ifndef _V8_WARP_H_
#define _V8_WARP_H_

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/*
isolate wrappers
*/
extern void* V8_NewIsolate();

extern void V8_DisposeIsolate(void* isolate);

/*
context wrappers
*/
extern void* V8_NewContext(void* isolate_ptr);

extern void V8_DisposeContext(void* context);

/*
script wrappers
*/
extern void* V8_Compile(void* context, const char* code, void* script_origin, void* script_data);

extern void V8_DisposeScript(void* script);

extern void* V8_RunScript(void* context, void* script);

/*
script data wrappers
*/
extern void* V8_PreCompile(void* isolate_ptr, const char* code);

extern void* V8_NewScriptData(const char* data, int length);

extern void V8_DisposeScriptData(void* script_data);

extern int V8_ScriptDataLength(void* script_data);

extern const char* V8_ScriptDataGetData(void* script_data);

extern int V8_ScriptDataHasError(void* script_data);

/*
script origin wrappers
*/
extern void* V8_NewScriptOrigin(void* isolate_ptr, const char* name, int line_offset, int column_offset);

extern void V8_DisposeScriptOrigin(void* script_origin);

/*
value wrappers
*/
extern void V8_DisposeValue(void* value);

extern char* V8_ValueToString(void* value);

extern int V8_ValueIsUndefined(void* value);

extern int V8_ValueIsNull(void* value);

extern int V8_ValueIsTrue(void* value);

extern int V8_ValueIsFalse(void* value);

extern int V8_ValueIsString(void* value);

extern int V8_ValueIsFunction(void* value);

extern int V8_ValueIsArray(void* value);

extern int V8_ValueIsObject(void* value);

extern int V8_ValueIsBoolean(void* value);

extern int V8_ValueIsNumber(void* value);

extern int V8_ValueIsExternal(void* value);

extern int V8_ValueIsInt32(void* value);

extern int V8_ValueIsUint32(void* value);

extern int V8_ValueIsDate(void* value);

extern int V8_ValueIsBooleanObject(void* value);

extern int V8_ValueIsNumberObject(void* value);

extern int V8_ValueIsStringObject(void* value);

extern int V8_ValueIsNativeError(void* value);

extern int V8_ValueIsRegExp(void* value);

extern int V8_ValueGetBoolean(void* value);
  
extern double V8_ValueGetNumber(void* value);

extern int64_t V8_ValueGetInteger(void* value);

extern uint32_t V8_ValueGetUint32(void* value);

extern int32_t V8_ValueGetInt32(void* value);

/*
special values
*/
extern void* V8_Undefined(void* isolate_ptr);

extern void* V8_Null(void* isolate_ptr);

extern void* V8_True(void* isolate_ptr);

extern void* V8_False(void* isolate_ptr);

#ifdef __cplusplus
} // extern "C"
#endif

#endif