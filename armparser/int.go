package armparser

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strconv"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Int converts a value to an integer.
func Int(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("Int", slog.Any("args", f.Args))
	if len(f.Args) != 1 {
		return nil, NewArgumentError("int", 1, len(f.Args))
	}
	val, err := f.Args[0].Evaluate(ctx, evalCtx, RegistryFromContext(ctx))
	if err != nil {
		return nil, err
	}
	switch v := val.(type) {
	case int:
		return v, nil
	case float64:
		return int(math.Trunc(v)), nil
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("int: cannot convert string %q to int: %w", v, err)
		}
		return i, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return nil, fmt.Errorf("int: unsupported type %T", val)
	}
}
