package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestNot(t *testing.T) {
	tcs := testCases{
		{
			desc:     "not true",
			in:       "[not(true)]",
			ctx:      map[string]any{},
			expected: false,
		},
		{
			desc:     "not false",
			in:       "[not(false)]",
			ctx:      map[string]any{},
			expected: true,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
