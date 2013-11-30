go-v8
=====

v8 JavaScript engine binding for Go.

Features
=======

* Thread safe
* Thorough and careful testing
* Boolean, Number, String, Object, Array, Regexp, Function
* Compile JavaScript and run
* Save and load pre-compiled script data
* Create JavaScript context with global object template
* Operate JavaScript object properties and array elements in Go
* Define JavaScript object template in Go with property accessors and interceptors
* Define JavaScript function template in Go
* Catch JavaScript exception in Go
* Throw JavaScript exception by Go
* JSON parse and generate

Install
=======

Prepare you computer:

1. make sure you have Go version 1.2
2. make sure there has 'curl' or 'wget' command
3. make sure there has 'git' command

For 'curl' user. please run this shell command:

```
curl -O https://raw.github.com/realint/go-v8/master/get.sh && chmod +x get.sh && ./get.sh go-v8
```

For 'wget' user. Please run this shell command:

```
wget -O get.sh https://raw.github.com/realint/go-v8/master/get.sh && chmod +x get.sh && ./get.sh go-v8
```

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

Hello World
===========

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

