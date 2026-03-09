package armparser

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Contains checks if a container (string, array, or object) contains an item.
// For strings, the comparison is case-insensitive.
func Contains(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("Contains", slog.Any("args", f.Args))
	defer lgr.Debug("Contains done")
	if len(f.Args) != 2 {
		return nil, NewArgumentError("contains", 2, len(f.Args))
	}
	reg := RegistryFromContext(ctx)
	container, err := f.Args[0].Evaluate(ctx, evalCtx, reg)
	if err != nil {
		return nil, err
	}
	item, err := f.Args[1].Evaluate(ctx, evalCtx, reg)
	if err != nil {
		return nil, err
	}
	switch c := container.(type) {
	case string:
		s, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("contains: when container is a string, itemToFind must be a string, got %T", item)
		}
		return strings.Contains(strings.ToLower(c), strings.ToLower(s)), nil
	case []any:
		for _, v := range c {
			if v == item {
				return true, nil
			}
		}
		return false, nil
	case map[string]any:
		key, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("contains: when container is an object, itemToFind must be a string key, got %T", item)
		}
		_, found := c[key]
		return found, nil
	default:
		return nil, fmt.Errorf("contains: container must be a string, array, or object, got %T", container)
	}
}
