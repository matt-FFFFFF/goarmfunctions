package armparser

import (
	"context"
	"fmt"
	"time"
)

// UtcNow returns the current UTC time as an ISO 8601 string.
func UtcNow(ctx context.Context, f *FunctionCall, evalCtx EvalContext) (any, error) {
	if len(f.Args) != 0 {
		return nil, fmt.Errorf("utcNow function expects 0 arguments but got %d", len(f.Args))
	}
	return time.Now().UTC().Format("2006-01-02T15:04:05.0000000Z"), nil
}
