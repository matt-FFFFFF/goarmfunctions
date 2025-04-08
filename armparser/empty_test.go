package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestEmpty(t *testing.T) {
	tcs := testCases{
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
