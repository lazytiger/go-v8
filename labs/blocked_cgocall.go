package main

/*
#include <stdio.h>

void slow_c_call() {
	FILE *big_file = fopen("bigfile", "w");

	if (big_file == NULL)
		return;

	for (int i = 0; i < 1024 * 1024 * 1024; i ++) {
		fputs("0", big_file);
	}

	fclose(big_file);
}
*/
import "C"
import "fmt"
import "time"
import "runtime"

func main() {
	runtime.GOMAXPROCS(1)

	go func() {
		for i := 0; true; i++ {
			time.Sleep(time.Second)
			fmt.Println("This is goroutine1")
		}
	}()

	t1 := time.Now()
	C.slow_c_call()
	t2 := time.Now()

	fmt.Printf("blocked %fs\n", t2.Sub(t1).Seconds())
}
