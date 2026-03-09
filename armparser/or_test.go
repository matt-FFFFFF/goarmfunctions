package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestOr(t *testing.T) {
	tcs := testCases{
		{
			desc:     "one true",
			in:       "[or(false, true)]",
			ctx:      map[string]any{},
			expected: true,
		},
		{
			desc:     "both false",
			in:       "[or(false, false)]",
			ctx:      map[string]any{},
			expected: false,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
