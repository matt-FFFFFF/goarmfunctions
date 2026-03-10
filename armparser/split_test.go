package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestSplit(t *testing.T) {
	tcs := testCases{
		{
			desc:     "split by comma",
			in:       "[split('a,b,c', ',')]",
			ctx:      map[string]any{},
			expected: []any{"a", "b", "c"},
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
