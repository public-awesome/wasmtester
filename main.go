package main

import (
	"fmt"
	"log"
	"strings"

	"go.starlark.net/starlark"
)

// ExampleExecFile demonstrates a simple embedding
// of the Starlark interpreter into a Go program.
func ExampleExecFile() {
	const data = `
code_id = store_code("sg721_base.wasm", "sg721_base")
print(code_id)
print(greeting + ", world")
print(repeat("one"))
print(repeat("mur", 2))
squares = [x*x for x in range(10)]
`

	storeCode := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		contractName := "c"
		starlark.UnpackArgs(b.Name(), args, kwargs, contractName, &contractName)
		return starlark.MakeInt64(2), nil
	}
	// repeat(str, n=1) is a Go function called from Starlark.
	// It behaves like the 'string * int' operation.
	repeat := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var s string
		var n int = 1
		if err := starlark.UnpackArgs(b.Name(), args, kwargs, "s", &s, "n?", &n); err != nil {
			return nil, err
		}

		return starlark.String(strings.Repeat(s, n)), nil
	}

	// The Thread defines the behavior of the built-in 'print' function.
	thread := &starlark.Thread{
		Name:  "example",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}

	// This dictionary defines the pre-declared environment.
	predeclared := starlark.StringDict{
		"greeting":   starlark.String("hello"),
		"repeat":     starlark.NewBuiltin("repeat", repeat),
		"store_code": starlark.NewBuiltin("store_code", storeCode),
	}

	// Execute a program.
	globals, err := starlark.ExecFile(thread, "apparent/filename.star", data, predeclared)
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			log.Fatal(evalErr.Backtrace())
		}
		log.Fatal(err)
	}

	// Print the global environment.
	fmt.Println("\nGlobals:")
	for _, name := range globals.Keys() {
		v := globals[name]
		fmt.Printf("%s (%s) = %s\n", name, v.Type(), v.String())
	}

	// Output:
	// hello, world
	// one
	// murmur
	//
	// Globals:
	// squares (list) = [0, 1, 4, 9, 16, 25, 36, 49, 64, 81]
}

func main() {
	ExampleExecFile()
}
