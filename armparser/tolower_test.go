package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestToLower(t *testing.T) {
	tcs := testCases{
		{
			desc:     "no function",
			in:       "hello",
			expected: "hello",
			ctx:      nil,
			err:      nil,
		},
		{
			desc:     "lowercase",
			in:       "[toLower('hello')]",
			expected: "hello",
			ctx:      nil,
			err:      nil,
		},
		{
			desc:     "mixed case",
			in:       "[toLower('HELlo')]",
			expected: "hello",
			ctx:      nil,
			err:      nil,
		},
		{
			desc:     "nested function",
			in:       "[toLower(if(true, 'HELLO', 'WORLD'))]",
			expected: "hello",
			ctx:      nil,
			err:      nil,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
