#!/bin/bash

v8_version="3.23.0"
v8_path="v8-$v8_version"

# check again
libv8_base="`find $v8_path/out/native/ -name 'libv8_base.*.a' | head -1`"
libv8_snapshot="`find $v8_path/out/native/ -name 'libv8_snapshot.a' | head -1`"
if [ libv8_base == '' ] || [ libv8_snapshot == '' ]; then
	echo >&2 "V8 build failed?"
	exit
fi

# build
CGO_LDFLAGS="$libv8_base $libv8_snapshot -lrt" \
CGO_CFLAGS="-I $v8_path/include" \
CGO_CXXFLAGS="-I $v8_path/include" \
go test -run="$1" -bench="$2" -v
