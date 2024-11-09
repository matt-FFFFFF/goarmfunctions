package armparser

import (
	"fmt"
)

// If evaluates the first argument and returns the second argument if it is true, otherwise it returns the third argument.
func If(f *FunctionCall, ctx EvalContext) (interface{}, error) {
	if len(f.Args) != 3 {
		return nil, fmt.Errorf("if function requires 3 arguments")
	}
	condition, err := f.Args[0].Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	if condition.(bool) {
		return f.Args[1].Evaluate(ctx)
	}
	return f.Args[2].Evaluate(ctx)
}
