package armparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArmParser(t *testing.T) {
	testCases := []struct {
		desc     string
		ctx      EvalContext
		in       string
		expected any
		parseErr error
		evalErr  error
	}{
		{
			desc:     "String literal",
			ctx:      nil,
			in:       "Just a string",
			expected: "Just a string",
			parseErr: nil,
			evalErr:  nil,
		},
		{
			desc:     "Whitespace",
			ctx:      nil,
			in:       " ",
			expected: " ",
			parseErr: nil,
			evalErr:  nil,
		},
		{
			desc:     "Embedded function",
			ctx:      nil,
			in:       "foo[if(true, 1, 2)]bar",
			expected: "foo1bar",
			parseErr: nil,
			evalErr:  nil,
		},
		{
			desc: "Multiple embedded functions",
			ctx: EvalContext{
				"test": "testvalue",
			},
			in:       "foo [if(true, 1, 2)] bar [parameters('test')] baz",
			expected: "foo 1 bar testvalue baz",
			parseErr: nil,
			evalErr:  nil,
		},
		{
			desc:     "symbols",
			ctx:      nil,
			in:       "!@£$%^&*()",
			expected: "!@£$%^&*()",
			parseErr: nil,
			evalErr:  nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := New()
			f, err := parser.ParseString("test", tC.in)
			require.ErrorIsf(t, err, tC.parseErr, "parse error not equal: %v", err)
			got, err := f.Evaluate(tC.ctx)
			require.ErrorIs(t, err, tC.evalErr)
			if err != nil {
				return
			}
			assert.Equal(t, tC.expected, got)
		})
	}
}
