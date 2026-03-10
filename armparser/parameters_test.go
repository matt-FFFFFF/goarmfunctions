package armparser

import (
	"context"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func TestParameters(t *testing.T) {
	tcs := testCases{
		{
			desc: "parameter found",
			in:   "[parameters('foo')]",
			ctx: map[string]any{
				"foo": "1",
			},
			expected: "1",
			err:      nil,
		},
		{
			desc: "parameter member dot access",
			in:   "[parameters('foo').bar]",
			ctx: map[string]any{
				"foo": map[string]any{
					"bar": "1",
				},
			},
			expected: "1",
			err:      nil,
		},
		{
			desc: "parameter multiple member dot access",
			in:   "[parameters('foo').bar.bat]",
			ctx: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"bat": "1",
					},
				},
			},
			expected: "1",
			err:      nil,
		},
		{
			desc: "parameter member square bracket access",
			in:   "[parameters('foo')['bar']]",
			ctx: map[string]any{
				"foo": map[string]any{
					"bar": "1",
				},
			},
			expected: "1",
			err:      nil,
		},
		{
			desc: "parameter member square bracket nested function",
			in:   "[parameters('foo')[if(true, 'bar', 'bat')]]",
			ctx: map[string]any{
				"foo": map[string]any{
					"bar": "1",
				},
			},
			expected: "1",
			err:      nil,
		},
		{
			desc: "parameter member multiple square bracket access",
			in:   "[parameters('foo')['bar']['bat']]",
			ctx: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"bat": "1",
					},
				},
			},
			expected: "1",
			err:      nil,
		},
		{
			desc: "parameter mixed multiple square bracket access",
			in:   "[parameters('foo')['bar'].bat]",
			ctx: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"bat": "1",
					},
				},
			},
			expected: "1",
			err:      nil,
		},
	}
	ctx := context.WithValue(context.Background(), logger.LoggerContextKey, logger.NewDebugLogger())
	runFunctionTest(ctx, t, tcs)
}
