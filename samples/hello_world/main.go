package main

import "github.com/realint/go-v8"

func main() {
	engine := v8.NewEngine()
	script := engine.Compile([]byte("'Hello ' + 'World!'"), nil, nil)
	context := engine.NewContext(nil)

	context.Scope(func(c *v8.Context) {
		result := script.Run()
		println(result.ToString())
	})
}
