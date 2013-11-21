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
--- PASS: Test_ThreadSafe5 (0.01 seconds)
=== RUN Test_ThreadSafe6
--- PASS: Test_ThreadSafe6 (0.06 seconds)
PASS
Benchmark_NewContext      10000     696899 ns/op
Benchmark_NewInteger    1000000       2309 ns/op
Benchmark_NewString     1000000       3304 ns/op
Benchmark_NewObject     1000000       3547 ns/op
Benchmark_NewArray0     1000000       1655 ns/op
Benchmark_NewArray5     1000000       2331 ns/op
Benchmark_NewArray20    1000000       2238 ns/op
Benchmark_NewArray100   1000000       2700 ns/op
Benchmark_Compile        200000      14404 ns/op
Benchmark_PreCompile     200000      13302 ns/op
Benchmark_RunScript     1000000       2214 ns/op
Benchmark_JsFunction    1000000       2096 ns/op
Benchmark_GoFunction     500000       3587 ns/op
Benchmark_Getter         500000       3176 ns/op
Benchmark_Setter         500000       9532 ns/op
ok  	github.com/realint/go-v8	61.815s
done
```
