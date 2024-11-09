package armparser

import "fmt"

// If evaluates the first argument and returns the second argument if it is true, otherwise it returns the third argument.
func If(f *FunctionCall, ctx EvalContext) (any, error) {
	if len(f.Args) != 3 {
		return nil, NewArgumentError("if", 3, len(f.Args))
	}
	condition, err := f.Args[0].Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	conditionB, ok := condition.(bool)
	if !ok {
		return nil, fmt.Errorf("if condition must be a boolean, got %T", condition)
	}
	if conditionB {
		return f.Args[1].Evaluate(ctx)
	}
	return f.Args[2].Evaluate(ctx)
}
