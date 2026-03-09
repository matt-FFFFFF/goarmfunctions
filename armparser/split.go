package armparser

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Split splits a string by a delimiter and returns an array.
func Split(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("Split", slog.Any("args", f.Args))
	if len(f.Args) != 2 {
		return nil, NewArgumentError("split", 2, len(f.Args))
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
		return nil, fmt.Errorf("split: both arguments must be strings, got %T and %T", a, b)
	}
	parts := strings.Split(s1, s2)
	result := make([]any, len(parts))
	for i, p := range parts {
		result[i] = p
	}
	return result, nil
}
