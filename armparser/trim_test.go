package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestTrim(t *testing.T) {
	tcs := testCases{
		{
			desc:     "trim whitespace",
			in:       "[trim('  hello  ')]",
			ctx:      map[string]any{},
			expected: "hello",
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
