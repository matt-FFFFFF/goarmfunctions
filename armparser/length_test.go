package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestLength(t *testing.T) {
	tcs := testCases{
		{
			desc:     "string length",
			in:       "[length('hello')]",
			ctx:      map[string]any{},
			expected: 5,
		},
		{
			desc:     "empty string",
			in:       "[length('')]",
			ctx:      map[string]any{},
			expected: 0,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
