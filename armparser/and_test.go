package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestAnd(t *testing.T) {
	tcs := testCases{
		{
			desc:     "both true",
			in:       "[and(true, true)]",
			ctx:      map[string]any{},
			expected: true,
		},
		{
			desc:     "one false",
			in:       "[and(true, false)]",
			ctx:      map[string]any{},
			expected: false,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
