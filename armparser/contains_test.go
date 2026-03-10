package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestContains(t *testing.T) {
	tcs := testCases{
		{
			desc:     "string contains substring case-insensitive",
			in:       "[contains('Hello World', 'hello')]",
			ctx:      map[string]any{},
			expected: true,
		},
		{
			desc:     "string does not contain",
			in:       "[contains('Hello World', 'xyz')]",
			ctx:      map[string]any{},
			expected: false,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
