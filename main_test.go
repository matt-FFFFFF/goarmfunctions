package main

import (
	"fmt"

	"github.com/matt-FFFFFF/goarmfunctions/armparser"
)

func ExampleBasicUsage() {
	parser := armparser.New()
	f, err := parser.ParseString("example", "[if(equals('a', 'b'), 'a is equal to b', 'a is not equal to b')]")
	if err != nil {
		panic(err)
	}
	result, err := f.Evaluate(nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output: a is not equal to b
}
