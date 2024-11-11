package armparser

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Parameters returns the value of a parameter from the context.
// If the function results in an object with members, dot notation can be used to retrieve values.
// For example, `parameters('foo').bar` would return the value of `ctx['foo']“, then access the `.bar` property.
// If you want to use a function as a member, you can use square brackets.
// For example, `parameters('foo')[if(true, 'bar', 'bat')]` would return the value of `ctx['foo']“, then access the property with the key returned by the `if` function.
func Parameters(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.GetLogger(ctx)
	lgr.Debug("Parameters", slog.Any("args", f.Args))
	defer lgr.Debug("Parameters done")
	if len(f.Args) != 1 || f.Args[0].String == nil {
		lgr.Error("Parameters - Invalid number of arguments", slog.Int("expected", 1), slog.Int("actual", len(f.Args)))
		return nil, NewArgumentError("parameters", 1, len(f.Args))
	}
	paramName := *f.Args[0].String
	value, ok := evalCtx[paramName]
	if !ok {
		lgr.Error("Parameters - Parameter not found", slog.String("parameter", paramName))
		return nil, fmt.Errorf("parameter %s not found", paramName)
	}
	if len(f.MembersDot) == 0 && len(f.MembersStr) == 0 {
		lgr.Debug("Parameters - No members - returning parameter value", slog.String("parameter", paramName), slog.Any("value", value))
		return value, nil
	}

	// First evaluate square bracket members
	for i, member := range f.MembersStr {
		lgr.Debug("Parameters - Evaluating square bracket member", slog.Any("member", *member), slog.Int("index", i))
		memberValue, err := member.Evaluate(ctx, evalCtx)
		if err != nil {
			lgr.Error("Parameters - Error evaluating square bracket member", slog.String("error", err.Error()))
			return nil, err
		}
		if value, ok = value.(map[string]any)[memberValue.(string)]; !ok {
			return nil, fmt.Errorf("member %s not found", memberValue)
		}
	}

	// Next evaluate dot members
	lgr.Debug("Parameters - Evaluating dot members", slog.Any("membersDot", f.MembersDot), slog.Any("membersStr", f.MembersStr))
	for i, member := range f.MembersDot {
		lgr.Debug("Parameters - Evaluating dot member", slog.String("member", member), slog.Int("index", i))
		if value, ok = value.(map[string]any)[member]; !ok {
			lgr.Error("Parameters - Member not found", slog.String("member", member))
			return nil, fmt.Errorf("member %s not found", member)
		}
	}
	lgr.Debug("Parameters - Returning member value", slog.Any("value", value))
	return value, nil
}
