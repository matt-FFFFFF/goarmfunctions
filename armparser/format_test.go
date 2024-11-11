package armparser

import (
	"context"
	"errors"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestFormat(t *testing.T) {
	tcs := testCases{
		{
			desc:     "no tokens",
			in:       "Hello world",
			ctx:      nil,
			expected: "Hello world",
			err:      nil,
		},
		{
			desc:     "one token",
			in:       "[format('Hello {0}', 'world')]",
			ctx:      nil,
			expected: "Hello world",
			err:      nil,
		},
		{
			desc:     "unsupported number token",
			in:       "[format('Hello {0}', 1)]",
			ctx:      nil,
			expected: "",
			err:      errors.New("format() only supports strings at this time. Format argument 0 is not a string"),
		},
		{
			desc:     "two tokens",
			in:       "[format('{0} {1}', 'Hello', 'world')]",
			ctx:      nil,
			expected: "Hello world",
			err:      nil,
		},
		{
			desc:     "one token nested function",
			in:       "[format('Hello {0}', parameters('foo'))]",
			ctx:      EvalContext{"foo": "world"},
			expected: "Hello world",
			err:      nil,
		},
		{
			desc:     "not enough arguments",
			in:       "[format('Hello {0}')]",
			ctx:      EvalContext{"foo": "world"},
			expected: "",
			err: &ArgumentError{
				function: "format",
				expected: 2,
				got:      1,
			},
		},
		{
			desc:     "repeated token",
			in:       "[format('Hello {0} {0}', 'world')]",
			ctx:      EvalContext{"foo": "world"},
			expected: "",
			err:      errors.New("token {0} still found in format string after replacement"),
		},
		{
			desc:     "not enough tokens in string",
			in:       "[format('{0} {1}', 'Hello', 'world', 'missing')]",
			ctx:      EvalContext{"foo": "world"},
			expected: "",
			err:      errors.New("token {2} not found in format string"),
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.GetLogger(context.Background()))
	runFunctionTest(ctx, t, tcs)
}
