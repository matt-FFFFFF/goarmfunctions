package armparser

import (
	"context"
	"fmt"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestIf(t *testing.T) {
	tcs := testCases{
		{
			desc:     "bool literal true",
			in:       "[if(true, 1, 2)]",
			ctx:      EvalContext{},
			expected: 1,
			err:      nil,
		},
		{
			desc:     "bool literal false",
			in:       "[if(false, 1, 2)]",
			ctx:      EvalContext{},
			expected: 2,
			err:      nil,
		},
		{
			desc:     "string instead of bool",
			in:       "[if('a', 1, 2)]",
			ctx:      EvalContext{},
			expected: "",
			err:      fmt.Errorf("if condition must be a boolean, got %T", "a"),
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
