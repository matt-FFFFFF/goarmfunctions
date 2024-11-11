package armparser

import (
	"context"

	"log/slog"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Equals returns true if the first argument is equal to the second argument.
func Equals(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.GetLogger(ctx)
	lgr.Debug("Equals",
		slog.Any("args", slog.Any("args", f.Args)),
	)
	defer lgr.Debug("Equals done")
	if len(f.Args) != 2 {
		lgr.Error("Equals - Invalid number of arguments", slog.Int("expected", 2), slog.Int("actual", len(f.Args)))
		return nil, NewArgumentError("equals", 2, len(f.Args))
	}
	arg1, err := f.Args[0].Evaluate(ctx, evalCtx)
	if err != nil {
		lgr.Error("Equals - Error evaluating first argument", slog.String("error", err.Error()))
		return nil, err
	}
	arg2, err := f.Args[1].Evaluate(ctx, evalCtx)
	if err != nil {
		lgr.Error("Equals - Error evaluating second argument", slog.String("error", err.Error()))
		return nil, err
	}
	return arg1 == arg2, nil
}
