package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestLess(t *testing.T) {
	tcs := testCases{
		{
			desc:     "int less true",
			in:       "[less(3, 5)]",
			ctx:      map[string]any{},
			expected: true,
		},
		{
			desc:     "int less false",
			in:       "[less(5, 3)]",
			ctx:      map[string]any{},
			expected: false,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
