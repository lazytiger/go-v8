#include "v8.h"

using namespace v8;

int main(int argc, char* argv[]) {
	// Get the default Isolate created at startup.
	Isolate* isolate = Isolate::GetCurrent();

	// Create a stack-allocated handle scope.
	HandleScope handle_scope(isolate);

	// Create a new context.
	Handle<Context> context = Context::New(isolate);

	// Here's how you could create a Persistent handle to the context, if needed.
	Persistent<Context> persistent_context(isolate, context);

	// Enter the created context for compiling and
	// running the hello world script. 
	Context::Scope context_scope(context);

	Handle<FunctionTemplate> ctor = FunctionTemplate::New();
	context->Global()->Set(String::New("T"), ctor->GetFunction());

	Local<Value> result1 = Script::Compile(String::New("(new T) instanceof T"))->Run();
	String::AsciiValue ascii1(result1);
	printf("(new T) instanceof T => %s\n", *ascii1);

	Local<Value> result2 = Script::Compile(String::New("typeof(new T)"))->Run();
	String::AsciiValue ascii2(result2);
	printf("typeof(new T) => %s\n", *ascii2);

	return 0;
}
