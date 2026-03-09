package armparser

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Trim trims whitespace from a string.
func Trim(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("Trim", slog.Any("args", f.Args))
	if len(f.Args) != 1 {
		return nil, NewArgumentError("trim", 1, len(f.Args))
	}
	val, err := f.Args[0].Evaluate(ctx, evalCtx, RegistryFromContext(ctx))
	if err != nil {
		return nil, err
	}
	s, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("trim: argument must be a string, got %T", val)
	}
	return strings.TrimSpace(s), nil
}
