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
NewContext     249474 ns/op
NewInteger        984 ns/op
NewString         983 ns/op
NewObject        1036 ns/op
NewArray0        1130 ns/op
NewArray5        1314 ns/op
NewArray20       1666 ns/op
NewArray100      3124 ns/op
Compile         11059 ns/op
PreCompile      11697 ns/op
RunScript        1085 ns/op
JsFunction       1114 ns/op
GoFunction       1630 ns/op
Getter           2060 ns/op
Setter           2934 ns/op
TryCatch        43097 ns/op
```

Hello World
===========

This 'Hello World' program shows how to use go-v8 to compile and run JavaScript code then get the result.

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

Contexts
========

The description in V8 embedding guide:

> In V8, a context is an execution environment that allows separate, unrelated, JavaScript applications to run in a single instance of V8. You must explicitly specify the context in which you want any JavaScript code to be run.

In go-v8, you can create many contexts from a V8 engine instance. When you want to run some JavaScript in a context. You need to enter the context by calling Scope() and run the JavaScript in the callback.

```go
context.Scope(func(cs v8.ContextScope){
	script.Run()
})
```

Context in V8 is necessary. So in go-v8 you can do this:

```go
context.Scope(func(cs v8.ContextScope) {
	context2 := engine.NewContext()
	context2.Scope(func(cs2 v8.ContextScope) {

	})
})
```

But, please note. Don't take any JavaScript value out scope.

When Scope() return, all of the JavaScript value created in this scope will be destroyed.

More
====

Please read the 'v8_all_test.go'.
