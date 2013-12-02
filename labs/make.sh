#!/bin/bash

v8_version="3.23.0"
v8_path="../v8-$v8_version"

# check again
libv8_base="`find $v8_path/out/native/ -name 'libv8_base.*.a' | head -1`"
libv8_snapshot="`find $v8_path/out/native/ -name 'libv8_snapshot.a' | head -1`"
if [ libv8_base == '' ] || [ libv8_snapshot == '' ]; then
	echo >&2 "V8 build failed?"
	exit
fi

LDFLAGS="$libv8_base $libv8_snapshot $librt"
CFLAGS="-I $v8_path/include"

if [ ! -d bin ]; then
	mkdir bin
fi

g++ $CFLAGS hello_world.cpp -o bin/hello_world $LDFLAGS

g++ $CFLAGS thread_test1.cpp -o bin/thread_test1 $LDFLAGS
g++ $CFLAGS thread_test2.cpp -o bin/thread_test2 $LDFLAGS
g++ $CFLAGS thread_test3.cpp -o bin/thread_test3 $LDFLAGS

g++ $CFLAGS instance_template.cpp -o bin/instance_template $LDFLAGS
