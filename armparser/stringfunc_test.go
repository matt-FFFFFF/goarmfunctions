package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestStringFunc(t *testing.T) {
	tcs := testCases{
		{
			desc:     "int to string",
			in:       "[string(42)]",
			ctx:      map[string]any{},
			expected: "42",
		},
		{
			desc:     "bool to string",
			in:       "[string(true)]",
			ctx:      map[string]any{},
			expected: "true",
		},
		{
			desc:     "string passthrough",
			in:       "[string('hello')]",
			ctx:      map[string]any{},
			expected: "hello",
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
