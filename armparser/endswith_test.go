package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestEndsWith(t *testing.T) {
	tcs := testCases{
		{
			desc:     "case-insensitive match",
			in:       "[endsWith('Hello World', 'WORLD')]",
			ctx:      map[string]any{},
			expected: true,
		},
		{
			desc:     "no match",
			in:       "[endsWith('Hello World', 'Hello')]",
			ctx:      map[string]any{},
			expected: false,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
