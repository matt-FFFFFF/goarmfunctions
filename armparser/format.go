package armparser

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Format returns a formatted string.
// E.g. format("Hello {0}", "world") returns "Hello world".
// Note this function does not support String.Format composite formatting options like .NET (https://learn.microsoft.com/dotnet/standard/base-types/composite-formatting).
// It currently only supports strings.
func Format(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("Format", slog.Any("args", f.Args))
	defer lgr.Debug("Format done")
	if len(f.Args) < 2 {
		lgr.Error("Format - Invalid number of arguments", slog.Int("expected", 2), slog.Int("actual", len(f.Args)))
		return nil, NewArgumentError("format", 2, len(f.Args))
	}
	lgr.Debug("Evaluating format string")
	format, err := f.Args[0].Evaluate(ctx, evalCtx)
	if err != nil {
		lgr.Error("Format - Error evaluating format string", slog.String("error", err.Error()))
		return nil, err
	}
	for i, arg := range f.Args[1:] {
		lgr.Debug("Format - Evaluating argument", slog.Int("index", i+1))
		val, err := arg.Evaluate(ctx, evalCtx)
		if err != nil {
			lgr.Error("Format - Error evaluating argument", slog.Int("index", i+1), slog.String("error", err.Error()))
			return nil, err
		}
		if _, ok := val.(string); !ok {
			lgr.Error("Format - Argument is not a string", slog.Int("index", i+1), slog.Any("value", val))
			return nil, fmt.Errorf("format() only supports strings at this time. Format argument %d is not a string", i)
		}
		token := fmt.Sprintf("{%d}", i)
		if !strings.Contains(format.(string), token) {
			lgr.Error("Format - Token not found in format string", slog.String("token", token), slog.String("format", format.(string)))
			return nil, fmt.Errorf("token %s not found in format string", token)
		}
		format = strings.Replace(format.(string), fmt.Sprintf("{%d}", i), val.(string), 1)
		lgr.Debug("Format - Replaced token", slog.String("token", token), slog.String("value", val.(string)), slog.String("result", format.(string)))
		if strings.Contains(format.(string), token) {
			lgr.Error("Format - Token still found in format string after replacement", slog.String("token", token), slog.String("format", format.(string)))
			return nil, fmt.Errorf("token %s still found in format string after replacement", token)
		}
	}
	lgr.Debug("Format - Returning resultant string", slog.String("result", format.(string)))
	return format, nil
}
