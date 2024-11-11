package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestReplace(t *testing.T) {
	tcs := testCases{
		{
			desc:     "nothing to replace",
			in:       "[replace('Hello world', 'foo', 'bar')]",
			ctx:      nil,
			expected: "Hello world",
			err:      nil,
		},
		{
			desc:     "replace one",
			in:       "[replace('Goodbye world', 'Goodbye', 'Hello')]",
			ctx:      nil,
			expected: "Hello world",
			err:      nil,
		},
		{
			desc:     "replace one",
			in:       "[replace('Goodbye world Goodbye', 'Goodbye', 'Hello')]",
			ctx:      nil,
			expected: "Hello world Hello",
			err:      nil,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
