using namespace v8;

typedef struct _Args
{
		Isolate* isolate;
		char message[256];
} Args;

void* test(void* data)
{
	Args* args = (Args*) data;
	Isolate* isolate = args->isolate;

	if (isolate == NULL)
	{
		std::cout << "null isolate found" << std::endl;
		delete args;
		return NULL;
	}

	Isolate* current = Isolate::GetCurrent();
	if (current == NULL)
	{
		std::cout << "current isolate is null before locker" << std::endl;
	}
	else
	{
		std::cerr << "current isolate is not null before locker" << std::endl;
	}
	Locker locker(isolate);

	current = Isolate::GetCurrent();
	if (current != NULL)
	{
		std::cerr << "current isolate is not null after locker" << std::endl;
	}
	else
	{
		std::cout << "current isolate is null after locker" << std::endl;
	}

	Isolate::Scope isolate_scope(isolate);

	current = Isolate::GetCurrent();
	if (current == NULL)
	{
		std::cerr << "current isolate is null after enter" << std::endl;
	}
	else
	{
		std::cout << "current isolate is not null after enter" << std::endl;
	}

	// Create a stack-allocated handle scope.
	HandleScope handle_scope(isolate);

	// Create a new context.
	Handle<Context> context = Context::New(isolate);

	// Here's how you could create a Persistent handle to the context, if needed.
	Persistent<Context> persistent_context(isolate, context);

	// Enter the created context for compiling and
	// running the hello world script.
	Context::Scope context_scope(context);

	// Create a string containing the JavaScript source code.
	Handle<String> source = String::New(args->message);

	// Compile the source code.
	Handle<Script> script = Script::Compile(source);

	// Run the script to get the result.
	Handle<Value> result = script->Run();

	// The persistent handle needs to be eventually disposed.
	persistent_context.Dispose();

	// Convert the result to an ASCII string and print it.
	String::AsciiValue ascii(result);

	for (int i = 0; i < 3; i++)
	{
		std::cout << *ascii << std::endl;
		sleep(1);
	}

	delete args;
	return NULL;
}
