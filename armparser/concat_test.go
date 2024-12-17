package armparser

import (
	"context"
	"errors"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestConcat(t *testing.T) {
	tcs := testCases{
		{
			desc:     "Concat strings",
			in:       "[concat('Hello', ' ', 'world')]",
			ctx:      nil,
			expected: "Hello world",
			err:      nil,
		},
		{
			desc:     "Not enough arguments",
			in:       "[concat('Hello')]",
			ctx:      nil,
			expected: "",
			err:      NewArgumentError("concat", 2, 1),
		},
		{
			desc: "Concat arrays",
			in:   "[concat(parameters('firstArray'), parameters('secondArray'))]",
			ctx: EvalContext{
				"firstArray":  []any{"1-1", "1-2", "1-3"},
				"secondArray": []any{"2-1", "2-2", "2-3"},
			},
			expected: []any{"1-1", "1-2", "1-3", "2-1", "2-2", "2-3"},
			err:      nil,
		},
		{
			desc: "Concat arrays - not same type",
			in:   "[concat(parameters('firstArray'), parameters('secondArray'))]",
			ctx: EvalContext{
				"firstArray":  []any{"1-1", "1-2", "1-3"},
				"secondArray": []any{"2-1", "2-2", 1},
			},
			expected: []any{"1-1", "1-2", "1-3", "2-1", "2-2", 1},
			err:      nil,
		},
		{
			desc: "Concat arrays with string, error",
			in:   "[concat(parameters('firstArray'), parameters('secondArray'))]",
			ctx: EvalContext{
				"firstArray":  []any{"1-1", "1-2", "1-3"},
				"secondArray": "2-1",
			},
			expected: nil,
			err:      errors.New("Concat all argument must of same type, got string"),
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.LoggerFromContext(context.Background()))
	runFunctionTest(ctx, t, tcs)
}
