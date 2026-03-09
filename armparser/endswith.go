package armparser

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// EndsWith returns true if the first string ends with the second (case-insensitive).
func EndsWith(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("EndsWith", slog.Any("args", f.Args))
	if len(f.Args) != 2 {
		return nil, NewArgumentError("endsWith", 2, len(f.Args))
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
	s1, ok1 := a.(string)
	s2, ok2 := b.(string)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("endsWith: both arguments must be strings, got %T and %T", a, b)
	}
	return strings.HasSuffix(strings.ToLower(s1), strings.ToLower(s2)), nil
}
