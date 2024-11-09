package armparser

import "fmt"

// Parameters returns the value of a parameter from the context.
// If the function results in an object with members, dot notation can be used to retrieve values.
// For example, `parameters('foo').bar` would return the value of `ctx['foo']“, then access the `.bar` property.
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
