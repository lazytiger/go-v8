go-v8
=====

v8 javascript engine binding for golang.

How to install
==============

Need 'Go 1.2' for auto compile C++ file.

Use the command below to auto download and compile go-v8:

```
mkdir go-v8
cd go-v8 && mkdir -p bin pkg src/github.com/realint/go-v8
export GOPATH=`pwd`
cd src/github.com/realint/go-v8
git init
git remote add origin https://github.com/realint/go-v8
git pull origin master
./install.sh
```

Tools
=====

There is a shell script named 'install.sh'. It can auto download and compile v8 engine.

For example:

```
./install.sh
```

And there is a shell script named 'test.sh'. It can auto configure cgo environment variable and run go-v8 test.

For example:

```
./test.sh . .
```

The first argument of test.sh is test name pattern, second argument is benchmark name pattern.

For example:

```
./test.sh .
```

The command above will run test without benchmark.


Stability and Performance
=========================

```
=== RUN Test_HelloWorld
--- PASS: Test_HelloWorld (0.00 seconds)
=== RUN Test_PreCompile
--- PASS: Test_PreCompile (0.00 seconds)
=== RUN Test_Values
--- PASS: Test_Values (0.02 seconds)
=== RUN Test_Object
--- PASS: Test_Object (0.00 seconds)
=== RUN Test_Array
--- PASS: Test_Array (0.00 seconds)
=== RUN Test_Function
--- PASS: Test_Function (0.00 seconds)
=== RUN Test_Context
--- PASS: Test_Context (0.00 seconds)
=== RUN Test_ObjectTemplate
--- PASS: Test_ObjectTemplate (0.00 seconds)
=== RUN Test_UnderscoreJS
--- PASS: Test_UnderscoreJS (0.01 seconds)
=== RUN Test_ThreadSafe1
--- PASS: Test_ThreadSafe1 (0.05 seconds)
=== RUN Test_ThreadSafe2
--- PASS: Test_ThreadSafe2 (0.03 seconds)
=== RUN Test_ThreadSafe3
--- PASS: Test_ThreadSafe3 (0.05 seconds)
=== RUN Test_ThreadSafe4
--- PASS: Test_ThreadSafe4 (0.02 seconds)
=== RUN Test_ThreadSafe5
--- PASS: Test_ThreadSafe5 (0.02 seconds)
=== RUN Test_ThreadSafe6
--- PASS: Test_ThreadSafe6 (0.06 seconds)
PASS
Benchmark_NewContext	   10000	    702612 ns/op
Benchmark_NewInteger	 1000000	      3326 ns/op
Benchmark_NewString	 1000000	      2964 ns/op
Benchmark_NewObject	 1000000	      2774 ns/op
Benchmark_NewArray0	 1000000	      2502 ns/op
Benchmark_NewArray5	 1000000	      1782 ns/op
Benchmark_NewArray20	 1000000	      2140 ns/op
Benchmark_NewArray100	 1000000	      2740 ns/op
Benchmark_Compile	  100000	     14692 ns/op
Benchmark_PreCompile	  100000	     14749 ns/op
Benchmark_RunScript	 1000000	      2671 ns/op
Benchmark_JsFunction	 1000000	      1682 ns/op
Benchmark_GoFunction	  500000	      3949 ns/op
Benchmark_Getter	  500000	      4127 ns/op
Benchmark_Setter	  500000	      6413 ns/op
```
