package armparser

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Replace replaces all occurrences of a specified substring with another substring.
// The function signature is:
// replace(original, old, new)
func Replace(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.GetLogger(ctx)
	lgr.Debug("Replace", slog.Any("args", f.Args))
	defer lgr.Debug("Replace done")
	if len(f.Args) != 3 {
		lgr.Error("Replace - Invalid number of arguments", slog.Int("expected", 3), slog.Int("actual", len(f.Args)))
		return nil, NewArgumentError("replace", 3, len(f.Args))
	}
	original, err := f.Args[0].Evaluate(ctx, evalCtx)
	if err != nil {
		lgr.Error("Replace - Error evaluating original string", slog.String("error", err.Error()))
		return nil, err
	}
	old, err := f.Args[1].Evaluate(ctx, evalCtx)
	if err != nil {
		lgr.Error("Replace - Error evaluating old string", slog.String("error", err.Error()))
		return nil, err
	}
	new, err := f.Args[2].Evaluate(ctx, evalCtx)
	if err != nil {
		lgr.Error("Replace - Error evaluating new string", slog.String("error", err.Error()))
		return nil, err
	}

	if _, ok := original.(string); !ok {
		lgr.Error("Replace - Original is not a string", slog.Any("value", original))
		return nil, fmt.Errorf("replace() only supports strings at this time. Original argument is not a string: %v", original)
	}
	if _, ok := old.(string); !ok {
		lgr.Error("Replace - Old is not a string", slog.Any("value", old))
		return nil, fmt.Errorf("replace() only supports strings at this time. Old argument is not a string: %v", old)
	}
	if _, ok := new.(string); !ok {
		lgr.Error("Replace - New is not a string", slog.Any("value", new))
		return nil, fmt.Errorf("replace() only supports strings at this time. New argument is not a string: %v", new)
	}

	res := strings.Replace(original.(string), old.(string), new.(string), -1)
	lgr.Debug("Replace - Returning resultant string", slog.String("result", res))
	return res, nil
}
