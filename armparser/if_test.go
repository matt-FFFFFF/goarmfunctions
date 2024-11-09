package armparser

import "testing"

func TestIf(t *testing.T) {
	testCases := []struct {
		desc     string
		in       string
		ctx      EvalContext
		expected any
		err      error
	}{
		{
			desc:     "true",
			in:       "[if(true, 1, 2)]",
			ctx:      EvalContext{},
			expected: "1",
			err:      nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			parser := New()
			f, err := parser.ParseString("test", tC.in)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			result, err := f.Evaluate(tC.ctx)
			if result != tC.expected {
				t.Errorf("expected %v, got %v", tC.expected, result)
			}
			if err != tC.err {
				t.Errorf("expected %v, got %v", tC.err, err)
			}
		})
	}
}
