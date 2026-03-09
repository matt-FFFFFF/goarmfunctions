package armparser

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Parameters returns the value of a parameter from the context.
// Member access (.dot and ['bracket']) is handled automatically by FunctionCall.Evaluate().
func Parameters(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("Parameters", slog.Any("args", f.Args))
	defer lgr.Debug("Parameters done")
	if len(f.Args) != 1 || f.Args[0].String == nil {
		lgr.Error("Parameters - Invalid number of arguments", slog.Int("expected", 1), slog.Int("actual", len(f.Args)))
		return nil, NewArgumentError("parameters", 1, len(f.Args))
	}
	paramName := *f.Args[0].String

	// Find the parameters scope and resolve the value.
	// Strategy: find "parameters" scope, check _parameters map first, then fall back to direct lookup.
	var value any
	var ok bool

	paramScope := FindScope(evalCtx, "parameters")
	if paramScope != nil {
		// Check for _parameters key (rich context with nested parameter definitions)
		if params, found := paramScope.GetLocal("_parameters"); found {
			if pm, isMap := params.(map[string]any); isMap {
				value, ok = pm[paramName]
			}
		}
		// Fall back to direct lookup on the scope (FromMap compatibility)
		if !ok {
			value, ok = paramScope.GetLocal(paramName)
		}
	}

	if !ok {
		lgr.Error("Parameters - Parameter not found", slog.String("parameter", paramName))
		return nil, fmt.Errorf("parameter %s not found", paramName)
	}

	lgr.Debug("Parameters - returning parameter value", slog.String("parameter", paramName), slog.Any("value", value))
	return value, nil
}
