package armparser

import (
	"fmt"
	"strings"
)

// Format returns a formatted string.
// E.g. format("Hello {0}", "world") returns "Hello world".
// Note this function does not support String.Format composite formatting options like .NET (https://learn.microsoft.com/dotnet/standard/base-types/composite-formatting).
// It currently only supports strings.
func Format(f *FunctionCall, ctx EvalContext) (any, error) {
	if len(f.Args) < 2 {
		return nil, NewArgumentError("format", 2, len(f.Args))
	}
	format, err := f.Args[0].Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	for i, arg := range f.Args[1:] {
		val, err := arg.Evaluate(ctx)
		if _, ok := val.(string); !ok {
			return nil, fmt.Errorf("format() only supports strings at this time. Format argument %d is not a string", i)
		}
		if err != nil {
			return nil, err
		}
		token := fmt.Sprintf("{%d}", i)
		if !strings.Contains(format.(string), token) {
			return nil, fmt.Errorf("token %s not found in format string", token)
		}
		format = strings.Replace(format.(string), fmt.Sprintf("{%d}", i), val.(string), 1)
	}
	return format, nil
}
