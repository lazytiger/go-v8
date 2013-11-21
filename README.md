go-v8
=====

v8 javascript engine binding for golang.

How to install
==============

Need Go 1.2 for auto compile C++ file.

Use install.sh to auto download and compile v8 engine:

```
./install.sh
```

Use test.sh to run test. It will auto configure cgo environment variable:

```
./test.sh . .
```

The first argument of test.sh is test name pattern, second argument is benchmark name pattern.

For example:

```
./test.sh .
```

The command above will run test without benchmark.


