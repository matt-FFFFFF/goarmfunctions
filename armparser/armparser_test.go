package armparser

import (
	"context"
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
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
			desc:     "string literal",
			ctx:      nil,
			in:       "Just a string",
			expected: "Just a string",
			parseErr: nil,
			evalErr:  nil,
		},
		{
			desc:     "whitespace",
			ctx:      nil,
			in:       " ",
			expected: " ",
			parseErr: nil,
			evalErr:  nil,
		},
		{
			desc:     "embedded function",
			ctx:      nil,
			in:       "foo[if(true, 1, 2)]bar",
			expected: "foo1bar",
			parseErr: nil,
			evalErr:  nil,
		},
		{
			desc: "multiple embedded functions",
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
		{
			desc:     "nested function",
			ctx:      nil,
			in:       "[if(true, [if(true, 1, 2)], 3)]",
			expected: nil,
			parseErr: &participle.ParseError{
				Msg: "sub-expression ArmTemplatePart+ must match at least once",
				Pos: lexer.Position{
					Filename: "test",
					Offset:   0,
					Line:     1,
					Column:   1,
				},
			},
			evalErr: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := New()
			f, err := parser.ParseString("test", tC.in)
			require.Equalf(t, tC.parseErr, err, "parse error not equal: %v", err)
			if err != nil {
				return
			}
			got, err := f.Evaluate(context.Background(), tC.ctx)
			require.Equalf(t, tC.evalErr, err, "eval error not equal: %v", err)
			if err != nil {
				return
			}
			assert.Equal(t, tC.expected, got)
		})
	}
}
