package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestGreater(t *testing.T) {
	tcs := testCases{
		{
			desc:     "int greater true",
			in:       "[greater(5, 3)]",
			ctx:      map[string]any{},
			expected: true,
		},
		{
			desc:     "int greater false",
			in:       "[greater(3, 5)]",
			ctx:      map[string]any{},
			expected: false,
		},
		{
			desc:     "string greater",
			in:       "[greater('b', 'a')]",
			ctx:      map[string]any{},
			expected: true,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
