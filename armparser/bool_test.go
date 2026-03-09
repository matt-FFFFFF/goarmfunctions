package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestBool(t *testing.T) {
	tcs := testCases{
		{
			desc:     "string true",
			in:       "[bool('true')]",
			ctx:      map[string]any{},
			expected: true,
		},
		{
			desc:     "string false",
			in:       "[bool('false')]",
			ctx:      map[string]any{},
			expected: false,
		},
		{
			desc:     "int 1",
			in:       "[bool(1)]",
			ctx:      map[string]any{},
			expected: true,
		},
		{
			desc:     "int 0",
			in:       "[bool(0)]",
			ctx:      map[string]any{},
			expected: false,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
