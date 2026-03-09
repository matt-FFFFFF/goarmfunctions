package armparser

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// And returns true if all arguments are true.
func And(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("And", slog.Any("args", f.Args))
	if len(f.Args) < 2 {
		return nil, fmt.Errorf("and function expects at least 2 arguments but got %d", len(f.Args))
	}
	reg := RegistryFromContext(ctx)
	for _, a := range f.Args {
		val, err := a.Evaluate(ctx, evalCtx, reg)
		if err != nil {
			return nil, err
		}
		b, ok := val.(bool)
		if !ok {
			return nil, fmt.Errorf("and: all arguments must be boolean, got %T", val)
		}
		if !b {
			return false, nil
		}
	}
	return true, nil
}
