package armparser

import (
	"fmt"
)

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
