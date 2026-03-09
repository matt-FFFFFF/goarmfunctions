package armparser

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// StringFunc converts a value to its string representation.
func StringFunc(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("String", slog.Any("args", f.Args))
	if len(f.Args) != 1 {
		return nil, NewArgumentError("string", 1, len(f.Args))
	}
	val, err := f.Args[0].Evaluate(ctx, evalCtx, RegistryFromContext(ctx))
	if err != nil {
		return nil, err
	}
	switch v := val.(type) {
	case string:
		return v, nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}
