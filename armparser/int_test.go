package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestInt(t *testing.T) {
	tcs := testCases{
		{
			desc:     "string to int",
			in:       "[int('42')]",
			ctx:      map[string]any{},
			expected: 42,
		},
		{
			desc:     "bool true to int",
			in:       "[int(true)]",
			ctx:      map[string]any{},
			expected: 1,
		},
		{
			desc:     "int passthrough",
			in:       "[int(7)]",
			ctx:      map[string]any{},
			expected: 7,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
