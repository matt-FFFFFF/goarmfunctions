package armparser

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Length returns the length of a string, array, or object.
func Length(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("Length", slog.Any("args", f.Args))
	if len(f.Args) != 1 {
		return nil, NewArgumentError("length", 1, len(f.Args))
	}
	arg, err := f.Args[0].Evaluate(ctx, evalCtx, RegistryFromContext(ctx))
	if err != nil {
		return nil, err
	}
	switch v := arg.(type) {
	case string:
		return len(v), nil
	case []any:
		return len(v), nil
	case map[string]any:
		return len(v), nil
	default:
		return nil, fmt.Errorf("length: unsupported type %T", arg)
	}
}
