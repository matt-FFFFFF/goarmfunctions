package armparser

import (
	"context"
	"errors"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestEmpty(t *testing.T) {
	tcs := testCases{
		{
			desc:     "Wrong number of arguments",
			in:       "[empty('', '')]",
			ctx:      nil,
			expected: true,
			err:      NewArgumentError("empty", 1, 2),
		},
		{
			desc:     "Cannot evaluate arg",
			in:       "[empty(parameters(''))]",
			ctx:      nil,
			expected: true,
			err:      errors.New("parameter  not found"),
		},
		{
			desc:     "Empty string",
			in:       "[empty('')]",
			ctx:      nil,
			expected: true,
			err:      nil,
		},
		{
			desc:     "Not empty string",
			in:       "[empty('Hello')]",
			ctx:      nil,
			expected: false,
			err:      nil,
		},
		{
			desc:     "Whitespace string",
			in:       "[empty(' ')]",
			ctx:      nil,
			expected: false,
			err:      nil,
		},
		{
			desc: "Empty array",
			in:   "[empty(parameters('array'))]",
			ctx: EvalContext{
				"array": []any{},
			},
			expected: true,
			err:      nil,
		},
		{
			desc: "Not empty array",
			in:   "[empty(parameters('array'))]",
			ctx: EvalContext{
				"array": []any{1, 2, 3},
			},
			expected: false,
			err:      nil,
		},
		{
			desc: "Empty object",
			in:   "[empty(parameters('object'))]",
			ctx: EvalContext{
				"object": map[string]any{},
			},
			expected: true,
			err:      nil,
		},
		{
			desc: "Not empty object",
			in:   "[empty(parameters('object'))]",
			ctx: EvalContext{
				"object": map[string]any{"key": "value"},
			},
			expected: false,
			err:      nil,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.LoggerFromContext(context.Background()))
	runFunctionTest(ctx, t, tcs)
}
