package armparser

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Bool converts a value to a boolean.
func Bool(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("Bool", slog.Any("args", f.Args))
	if len(f.Args) != 1 {
		return nil, NewArgumentError("bool", 1, len(f.Args))
	}
	val, err := f.Args[0].Evaluate(ctx, evalCtx, RegistryFromContext(ctx))
	if err != nil {
		return nil, err
	}
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		switch v {
		case "true", "True", "TRUE", "1":
			return true, nil
		case "false", "False", "FALSE", "0":
			return false, nil
		default:
			return nil, fmt.Errorf("bool: cannot convert string %q to bool", v)
		}
	case int:
		return v != 0, nil
	case float64:
		return v != 0, nil
	default:
		return nil, fmt.Errorf("bool: unsupported type %T", val)
	}
}
