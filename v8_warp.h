#ifndef _V8_WARP_H_
#define _V8_WARP_H_

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/*
Isolate warppers
*/
extern void *V8_NewIsolate();

extern void V8_DisposeIsolate(void *isolate);

/*
Context warppers
*/
extern void *V8_NewContext(void *isolate_ptr);

extern void V8_DisposeContext(void *context);

/*
Script warppers
*/
extern void *V8_CompileScript(void *context, const char *code);

extern void V8_DisposeScript(void *script);

extern void *V8_RunScript(void *context, void *script);

/*
Value warppers
*/
extern void V8_DisposeValue(void *value);

extern char *V8_ValueToString(void *value);

extern int V8_ValueIsUndefined(void *value);

extern int V8_ValueIsNull(void *value);

extern int V8_ValueIsTrue(void *value);

extern int V8_ValueIsFalse(void *value);

extern int V8_ValueIsString(void *value);

extern int V8_ValueIsFunction(void *value);

extern int V8_ValueIsArray(void *value);

extern int V8_ValueIsObject(void *value);

extern int V8_ValueIsBoolean(void *value);

extern int V8_ValueIsNumber(void *value);

extern int V8_ValueIsExternal(void *value);

extern int V8_ValueIsInt32(void *value);

extern int V8_ValueIsUint32(void *value);

extern int V8_ValueIsDate(void *value);

extern int V8_ValueIsBooleanObject(void *value);

extern int V8_ValueIsNumberObject(void *value);

extern int V8_ValueIsStringObject(void *value);

extern int V8_ValueIsNativeError(void *value);

extern int V8_ValueIsRegExp(void *value);

extern int V8_ValueGetBoolean();
  
extern double V8_ValueGetNumber();

extern int64_t V8_ValueGetInteger();

extern uint32_t V8_ValueGetUint32();

extern int32_t V8_ValueGetInt32();

#ifdef __cplusplus
} // extern "C"
#endif

#endif