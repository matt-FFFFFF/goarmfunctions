package armparser

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// If evaluates the first argument and returns the second argument if it is true, otherwise it returns the third argument.
func If(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.GetLogger(ctx)
	lgr.Debug("If", slog.Any("args", f.Args))
	defer lgr.Debug("If done")
	if len(f.Args) != 3 {
		lgr.Error("If - Invalid number of arguments", slog.Int("expected", 3), slog.Int("actual", len(f.Args)))
		return nil, NewArgumentError("if", 3, len(f.Args))
	}
	condition, err := f.Args[0].Evaluate(ctx, evalCtx)
	if err != nil {
		return nil, err
	}
	conditionB, ok := condition.(bool)
	if !ok {
		lgr.Error("If - Condition is not a boolean", slog.Any("condition", condition))
		return nil, fmt.Errorf("if condition must be a boolean, got %T", condition)
	}
	if conditionB {
		lgr.Debug("If - Returning true branch")
		return f.Args[1].Evaluate(ctx, evalCtx)
	}
	lgr.Debug("If - Returning false branch")
	return f.Args[2].Evaluate(ctx, evalCtx)
}
