package armparser

import (
	"context"
	"log/slog"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// GreaterOrEquals returns true if arg1 >= arg2.
func GreaterOrEquals(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("GreaterOrEquals", slog.Any("args", f.Args))
	if len(f.Args) != 2 {
		return nil, NewArgumentError("greaterOrEquals", 2, len(f.Args))
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
	return compareValues(a, b, "greaterOrEquals")
}
