package armparser

import "fmt"

// Parameters returns the value of a parameter from the context.
func Parameters(f *FunctionCall, ctx EvalContext) (any, error) {
	if len(f.Args) != 1 || f.Args[0].String == nil {
		return nil, NewArgumentError("parameters", 1, len(f.Args))
	}
	paramName := *f.Args[0].String
	value, ok := ctx[paramName]
	if !ok {
		return nil, fmt.Errorf("parameter %s not found", paramName)
	}
	if len(f.Members) == 0 {
		return value, nil
	}
	for _, member := range f.Members {
		if value, ok = value.(map[string]any)[member]; !ok {
			return nil, fmt.Errorf("member %s not found", member)
		}
	}
	return value, nil
}
