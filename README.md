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

Below is the benchmark result on my iMac:

```
Benchmark_NewContext        5000            249474 ns/op
Benchmark_NewInteger     2000000               984 ns/op
Benchmark_NewString      2000000               983 ns/op
Benchmark_NewObject      1000000              1036 ns/op
Benchmark_NewArray0      1000000              1130 ns/op
Benchmark_NewArray5      1000000              1314 ns/op
Benchmark_NewArray20     1000000              1666 ns/op
Benchmark_NewArray100    1000000              3124 ns/op
Benchmark_Compile         200000             11059 ns/op
Benchmark_PreCompile      200000             11697 ns/op
Benchmark_RunScript      1000000              1085 ns/op
Benchmark_JsFunction     1000000              1114 ns/op
Benchmark_GoFunction     1000000              1630 ns/op
Benchmark_Getter         1000000              2060 ns/op
Benchmark_Setter         1000000              2934 ns/op
Benchmark_TryCatch         50000             43097 ns/op
```

How to use
==========

Hello World
-----------

Let's write a Hello World program to learn how to use go-v8.

At the begining, we need import go-v8 package and create a V8 engine instance.

```go
package main

import "github.com/realint/go-v8"

func main() {
	engine := v8.NewEngine()
}
```

NOTE: You can create many V8 engine but don't share any data (like Value, Object, Function and Context etc.) between engine instances.

And then, we need to compile the JavaScript code that we want to run.

```go
...
	script := engine.Compile([]byte("'Hello ' + 'World!'"), nil, nil)
...
```

NOTE: Script can compile one time and run many times.

Now, we need a context scope to run the script.

```go
...
	context := engine.NewContext(nil)

	context.Scope(func(cs v8.ContextScope) {
		script.Run()
	})
...
```

Last we can get the result and print it.

```go
package main

import "github.com/realint/go-v8"

func main() {
	engine := v8.NewEngine()
	script := engine.Compile([]byte("'Hello ' + 'World!'"), nil, nil)
	context := engine.NewContext(nil)

	context.Scope(func(cs v8.ContextScope) {
		result := script.Run()
		println(result.ToString())
	})
}

```

