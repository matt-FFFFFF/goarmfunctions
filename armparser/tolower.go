package armparser

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// ToLower converts the input string to lower case.
// Function signature: toLower('string')
func ToLower(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.GetLogger(ctx)
	lgr.Debug("If", slog.Any("args", f.Args))
	defer lgr.Debug("If done")
	if len(f.Args) != 1 {
		lgr.Error("ToLower - Invalid number of arguments", slog.Int("expected", 1), slog.Int("actual", len(f.Args)))
		return nil, NewArgumentError("toLower", 1, len(f.Args))
	}
	str, err := f.Args[0].Evaluate(ctx, evalCtx)
	if err != nil {
		return nil, err
	}
	strS, ok := str.(string)
	if !ok {
		lgr.Error("ToLower - Argument is not a string", slog.Any("argument", str))
		return nil, fmt.Errorf("toLower argument must be a string, got %T", str)
	}
	return strings.ToLower(strS), nil
}
