package armparser

import (
	"context"
	"errors"
	"log/slog"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Empty returns a bool depending on if the input array, object or string is empty.
// Function signature: empty(”)
func Empty(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("Empty", slog.Any("args", f.Args))
	defer lgr.Debug("Empty done")
	if len(f.Args) != 1 {
		lgr.Error("Empty - Invalid number of arguments", slog.Int("expected", 1), slog.Int("actual", len(f.Args)))
		return nil, NewArgumentError("empty", 1, len(f.Args))
	}
	arg, err := f.Args[0].Evaluate(ctx, evalCtx)
	if err != nil {
		return nil, err
	}

	switch a := arg.(type) {
	case string:
		if a != "" {
			return false, nil
		}
		return true, nil
	case []any:
		if len(a) != 0 {
			return false, nil
		}
		return true, nil
	case map[string]any:
		if len(a) != 0 {
			return false, nil
		}
		return true, nil
	}
	return nil, errors.New("Empty - Unsupported argument, please supply an object, string or an array")
}
