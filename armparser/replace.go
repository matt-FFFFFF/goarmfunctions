package armparser

import "strings"

// Replace replaces all occurrences of a specified substring with another substring.
// The function signature is:
// replace(original, old, new)
func Replace(f *FunctionCall, ctx EvalContext) (any, error) {
	if len(f.Args) != 3 {
		return nil, NewArgumentError("replace", 3, len(f.Args))
	}
	original, err := f.Args[0].Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	old, err := f.Args[1].Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	new, err := f.Args[2].Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return strings.Replace(original.(string), old.(string), new.(string), -1), nil
}
