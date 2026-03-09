package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestLessOrEquals(t *testing.T) {
	tcs := testCases{
		{
			desc:     "equal ints",
			in:       "[lessOrEquals(5, 5)]",
			ctx:      map[string]any{},
			expected: true,
		},
		{
			desc:     "less int",
			in:       "[lessOrEquals(4, 5)]",
			ctx:      map[string]any{},
			expected: true,
		},
		{
			desc:     "greater int",
			in:       "[lessOrEquals(6, 5)]",
			ctx:      map[string]any{},
			expected: false,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
