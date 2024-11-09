package armparser

import "fmt"

func Parameters(f *FunctionCall, ctx EvalContext) (interface{}, error) {
	if len(f.Args) != 1 || f.Args[0].String == nil {
		return nil, fmt.Errorf("parameters function requires a single string argument")
	}
	paramName := *f.Args[0].String
	value, ok := ctx[paramName]
	if !ok {
		return nil, fmt.Errorf("parameter %s not found", paramName)
	}
	return value, nil
}
