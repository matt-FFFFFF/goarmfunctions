package armparser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIf(t *testing.T) {
	testCases := []struct {
		desc     string
		in       string
		ctx      EvalContext
		expected any
		err      error
	}{
		{
			desc:     "bool literal",
			in:       "[if(true, 1, 2)]",
			ctx:      EvalContext{},
			expected: "1",
			err:      nil,
		},
		{
			desc:     "bool literal",
			in:       "[if(false, 1, 2)]",
			ctx:      EvalContext{},
			expected: "2",
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
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := New()
			f, err := parser.ParseString("test", tC.in)
			require.NoError(t, err)
			result, err := f.Evaluate(tC.ctx)
			require.Equalf(t, tC.err, err, "unexpected evaluate error: %v", err)
			if err != nil {
				return
			}
			assert.Equalf(t, tC.expected, result, "unexpected result: %v", result)
		})
	}
}
