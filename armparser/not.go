package armparser

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Not negates a boolean value.
func Not(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("Not", slog.Any("args", f.Args))
	if len(f.Args) != 1 {
		return nil, NewArgumentError("not", 1, len(f.Args))
	}
	val, err := f.Args[0].Evaluate(ctx, evalCtx, RegistryFromContext(ctx))
	if err != nil {
		return nil, err
	}
	b, ok := val.(bool)
	if !ok {
		return nil, fmt.Errorf("not: argument must be boolean, got %T", val)
	}
	return !b, nil
}
