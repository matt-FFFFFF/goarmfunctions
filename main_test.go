package goarmfunctions

import (
	"context"
	"fmt"
)

func ExampleLexAndParse() {
	result, err := LexAndParse(context.Background(), "[if(equals('a', 'b'), 'a is equal to b', 'a is not equal to b')]", nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output: a is not equal to b
}
