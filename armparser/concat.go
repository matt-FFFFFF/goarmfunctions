package armparser

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Concat concatenates arrays or strings using a variadic input.
// Function signature: concat('string', 'string', ...)
func Concat(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("Concat", slog.Any("args", f.Args))
	defer lgr.Debug("Concat done")
	if len(f.Args) < 2 {
		lgr.Error("Concat - Invalid number of arguments", slog.Int("expected", 2), slog.Int("actual", len(f.Args)))
		return nil, NewArgumentError("concat", 2, len(f.Args))
	}
	argZero, err := f.Args[0].Evaluate(ctx, evalCtx)
	if err != nil {
		return nil, err
	}

	switch a := argZero.(type) {
	case string:
		sb := strings.Builder{}
		sb.WriteString(a)
		for i := 1; i < len(f.Args); i++ {
			str, err := f.Args[i].Evaluate(ctx, evalCtx)
			if err != nil {
				return nil, err
			}
			strS, ok := str.(string)
			if !ok {
				lgr.Error("Concat - Argument is not a string", slog.Any("argument", str))
				return nil, fmt.Errorf("Concat all argument must of same type, got %T", str)
			}
			sb.WriteString(strS)
		}
		return sb.String(), nil
	case []any:
		for i := 1; i < len(f.Args); i++ {
			arr, err := f.Args[i].Evaluate(ctx, evalCtx)
			arrA, ok := arr.([]any)
			if !ok {
				lgr.Error("Concat - Argument is not an array", slog.Any("argument", arr))
				return nil, fmt.Errorf("Concat all argument must of same type, got %T", arr)
			}
			if err != nil {
				return nil, err
			}
			a = append(a, arrA...)
		}
		return a, nil
	}
	return nil, errors.New("Concat - Unsupported argument, please supply string or array")
}
