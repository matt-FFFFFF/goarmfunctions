package armparser

import (
	"context"
	"errors"
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
		{
			desc:     "wrong args",
			in:       "[toLower()]",
			expected: "",
			ctx:      nil,
			err: &ArgumentError{
				"toLower",
				1,
				0,
			},
		},
		{
			desc:     "not string",
			in:       "[toLower(1)]",
			expected: "",
			ctx:      nil,
			err:      errors.New("toLower argument must be a string, got int"),
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
