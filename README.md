go-v8
=====

v8 javascript engine binding for golang.

Features
=======

* thread safe
* thorough and careful testing
* number, string, object, array, regexp, function types
* access javascript object properties and array elements in Go
* define javascript object template in Go with property getter and setter callback
* define javascript function template in Go with callback
* compile or pre-compile script and run
* JSON parse and generate
* try catch and throw exception

How to install
==============

The easy way
------------

Prepare you computer:

1. make sure your 'go' is version 1.2
2. make sure there has 'curl' or 'wget' installed
3. make sure there has 'git' installed

There has some shell script in the project root directory for help you download and install v8 engine and go-v8.

Install v8 engine and go-v8 only need one line of shell command.

Use 'curl':

```
curl -O https://raw.github.com/realint/go-v8/master/get.sh && chmod +x get.sh && ./get.sh go-v8
```

Use 'wget':

```
wget -O get.sh https://raw.github.com/realint/go-v8/master/get.sh && chmod +x get.sh && ./get.sh go-v8
```

The hard way
------------

You can install go-v8 by manual:

1. download or clone go-v8 from github
2. make sure go-v8 package can be searched by your GOPATH setting
3. cd to go-v8 root directory
4. run install.sh to auto download and compile v8 engine
5. the install.sh will install and test go-v8 after v8 engine compiled

Read the get.sh and install.sh if you want to know the detail.

Stability and Performance
=========================

I write many test and benchmark to make sure go-v8 stable and efficient.

There is a shell script named 'test.sh' in the project. 

It can auto configure cgo environment variables and run go-v8 test.

For example:

```
./test.sh . .
```

The above command will run all of test and benchmark.

The first argument of test.sh is test name pattern, second argument is benchmark name pattern.

For example:

```
./test.sh ThreadSafe Array
```

The above command will run all of thread safe test and all of benchmark about Array type.

Below is the test and benchmark result on my iMac:

```
=== RUN Test_GetVersion
--- PASS: Test_GetVersion (0.00 seconds)
        v8_all_test.go:58: 3.23.0
=== RUN Test_HelloWorld
--- PASS: Test_HelloWorld (0.00 seconds)
=== RUN Test_TryCatch
--- PASS: Test_TryCatch (0.01 seconds)
=== RUN Test_PreCompile
--- PASS: Test_PreCompile (0.00 seconds)
=== RUN Test_Values
--- PASS: Test_Values (0.00 seconds)
=== RUN Test_Object
--- PASS: Test_Object (0.00 seconds)
=== RUN Test_Array
--- PASS: Test_Array (0.00 seconds)
=== RUN Test_Function
--- PASS: Test_Function (0.00 seconds)
=== RUN Test_ObjectTemplate
--- PASS: Test_ObjectTemplate (0.00 seconds)
=== RUN Test_Context
--- PASS: Test_Context (0.00 seconds)
        v8_all_test.go:552: Hello World!
=== RUN Test_UnderscoreJS
--- PASS: Test_UnderscoreJS (0.00 seconds)
=== RUN Test_JSON
--- PASS: Test_JSON (0.00 seconds)
=== RUN Test_ThreadSafe1
--- PASS: Test_ThreadSafe1 (0.06 seconds)
=== RUN Test_ThreadSafe2
--- PASS: Test_ThreadSafe2 (0.06 seconds)
=== RUN Test_ThreadSafe3
--- PASS: Test_ThreadSafe3 (0.09 seconds)
=== RUN Test_ThreadSafe4
--- PASS: Test_ThreadSafe4 (0.07 seconds)
=== RUN Test_ThreadSafe5
--- PASS: Test_ThreadSafe5 (0.02 seconds)
=== RUN Test_ThreadSafe6
--- PASS: Test_ThreadSafe6 (0.06 seconds)
PASS
Benchmark_NewContext       10000            688826 ns/op
Benchmark_NewInteger     1000000              2249 ns/op
Benchmark_NewString      1000000              3869 ns/op
Benchmark_NewObject      1000000              4258 ns/op
Benchmark_NewArray0      1000000              1884 ns/op
Benchmark_NewArray5      1000000              1558 ns/op
Benchmark_NewArray20     1000000              4677 ns/op
Benchmark_NewArray100    1000000              2290 ns/op
Benchmark_Compile         200000             13752 ns/op
Benchmark_PreCompile      200000             13879 ns/op
Benchmark_RunScript      5000000              1271 ns/op
Benchmark_JsFunction     5000000               719 ns/op
Benchmark_GoFunction      500000              4434 ns/op
Benchmark_Getter         1000000              2718 ns/op
Benchmark_Setter          500000              4523 ns/op
Benchmark_TryCatch        100000             26957 ns/op
ok      github.com/realint/v8   97.863s
```

Getting Started
===============

Start The Engine
----------------

At the begining. You need create a V8 engine instance for your program.

```
import "github.com/realint/go-v8"

var engine = v8.NewEngine()
```

NOTE: You can create many V8 engine but don't share any data (like Value, Object, Function and Context etc.) between engine instances.

Do Works In Context
-------------------

```
var context = engine.NewContext()

context.Scope(func(c *Context){
	// put your v8 operation code here
})
```

Compile And Run
---------------

You need compile the JavaScript code, before you run it.

```
engine.NewContext().Scope(func(c *Context){
	var script = engine.Compile("'Hello ' + 'World!'")
	script.Run()
})
```

A compiled script can run in different context.

```
var script = engine.Compile("x = 1")

engine.NewContext().Scope(func(c *Context){
	script.Run()
	// now there is a variable named 'x' in this context
})


engine.NewContext().Scope(func(c *Context){
	script.Run()
	// now there is a variable named 'x' in this context
})
```
