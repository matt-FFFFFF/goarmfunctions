package armparser

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Greater returns true if arg1 > arg2. Works on ints, floats, and strings.
func Greater(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("Greater", slog.Any("args", f.Args))
	if len(f.Args) != 2 {
		return nil, NewArgumentError("greater", 2, len(f.Args))
	}
	reg := RegistryFromContext(ctx)
	a, err := f.Args[0].Evaluate(ctx, evalCtx, reg)
	if err != nil {
		return nil, err
	}
	b, err := f.Args[1].Evaluate(ctx, evalCtx, reg)
	if err != nil {
		return nil, err
	}
	return compareValues(a, b, "greater")
}

func compareValues(a, b any, fname string) (any, error) {
	af, aIsNum := toFloat(a)
	bf, bIsNum := toFloat(b)
	if aIsNum && bIsNum {
		switch fname {
		case "greater":
			return af > bf, nil
		case "less":
			return af < bf, nil
		case "greaterOrEquals":
			return af >= bf, nil
		case "lessOrEquals":
			return af <= bf, nil
		}
	}
	as, aStr := a.(string)
	bs, bStr := b.(string)
	if aStr && bStr {
		switch fname {
		case "greater":
			return as > bs, nil
		case "less":
			return as < bs, nil
		case "greaterOrEquals":
			return as >= bs, nil
		case "lessOrEquals":
			return as <= bs, nil
		}
	}
	return nil, fmt.Errorf("%s: arguments must be both numbers or both strings, got %T and %T", fname, a, b)
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case float64:
		return n, true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}
