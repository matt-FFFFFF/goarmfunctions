package armparser

import "fmt"

func Equals(f *FunctionCall, ctx EvalContext) (interface{}, error) {
	if len(f.Args) != 2 {
		return nil, fmt.Errorf("equals function requires 2 arguments")
	}
	arg1, err := f.Args[0].Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	arg2, err := f.Args[1].Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return arg1 == arg2, nil
}