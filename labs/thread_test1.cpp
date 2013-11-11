#include <v8.h>
#include <iostream>
#include <pthread.h>
#include <unistd.h>

#include "thread_test.h"

using namespace v8;

//
// multi-isolate in multi-thread
//
int main(int argc, char* argv[])
{
	pthread_t pids[256];
	int num = 3;
	for (int i = 0; i < num; i++)
	{
		Args* args = new Args();
		args->isolate = Isolate::New();
		int n = sprintf(args->message, "\"thread\" + %d", i);
		args->message[n] = '\0';
		int ret = pthread_create(&pids[i], NULL, test, args);
		if (ret != 0)
		{
			std::cout << "create pthread" << i << " failed" << std::endl;
			return 1;
		}
	}

	for (int i = 0; i < num; i++)
	{
		pthread_join(pids[i], NULL);
	}

	std::cerr << "test done" << std::endl;

	return 0;
}
