package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestStartsWith(t *testing.T) {
	tcs := testCases{
		{
			desc:     "case-insensitive match",
			in:       "[startsWith('Hello World', 'hello')]",
			ctx:      map[string]any{},
			expected: true,
		},
		{
			desc:     "no match",
			in:       "[startsWith('Hello World', 'world')]",
			ctx:      map[string]any{},
			expected: false,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
