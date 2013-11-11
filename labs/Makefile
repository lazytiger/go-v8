V8_CFLAGS=-I ../v8/include
V8_LFLAGS=-L ../v8/out/native -lv8_base.x64 -lv8_nosnapshot.x64 -licuuc -licudata -licui18n

all: hello_world thread_test1 thread_test2 thread_test3

clean:
	rm hello_world
	rm thread_test1
	rm thread_test2
	rm thread_test3

hello_world: hello_world.cpp
	g++ ${V8_CFLAGS} hello_world.cpp -o hello_world ${V8_LFLAGS}

thread_test1: thread_test1.cpp
	g++ ${V8_CFLAGS} thread_test1.cpp -o thread_test1 ${V8_LFLAGS}

thread_test2: thread_test2.cpp
	g++ ${V8_CFLAGS} thread_test2.cpp -o thread_test2 ${V8_LFLAGS}

thread_test3: thread_test3.cpp
	g++ ${V8_CFLAGS} thread_test3.cpp -o thread_test3 ${V8_LFLAGS}