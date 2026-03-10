package armparser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_RegisterAndLookup(t *testing.T) {
	r := NewRegistry()
	called := false
	r.Register("myFunc", func(ctx context.Context, call *FunctionCall, evalCtx EvalContext) (any, error) {
		called = true
		return "ok", nil
	})

	fn, ok := r.Lookup("myFunc")
	require.True(t, ok)
	result, err := fn(context.Background(), &FunctionCall{}, nil)
	require.NoError(t, err)
	assert.Equal(t, "ok", result)
	assert.True(t, called)
}

func TestRegistry_LookupCaseInsensitive(t *testing.T) {
	r := NewRegistry()
	r.Register("MyFunc", func(ctx context.Context, call *FunctionCall, evalCtx EvalContext) (any, error) {
		return "found", nil
	})

	for _, name := range []string{"myfunc", "MYFUNC", "MyFunc", "myFunc", "MYFUNC"} {
		fn, ok := r.Lookup(name)
		require.True(t, ok, "expected to find function with name %q", name)
		result, err := fn(context.Background(), &FunctionCall{}, nil)
		require.NoError(t, err)
		assert.Equal(t, "found", result)
	}
}

func TestRegistry_LookupUnknown(t *testing.T) {
	r := NewRegistry()
	_, ok := r.Lookup("nonexistent")
	assert.False(t, ok)
}

func TestRegistry_Clone(t *testing.T) {
	r := NewRegistry()
	r.Register("original", func(ctx context.Context, call *FunctionCall, evalCtx EvalContext) (any, error) {
		return "original", nil
	})

	cloned := r.Clone()

	// Clone has the original function
	fn, ok := cloned.Lookup("original")
	require.True(t, ok)
	result, err := fn(context.Background(), &FunctionCall{}, nil)
	require.NoError(t, err)
	assert.Equal(t, "original", result)

	// Adding to clone doesn't affect original
	cloned.Register("extra", func(ctx context.Context, call *FunctionCall, evalCtx EvalContext) (any, error) {
		return "extra", nil
	})
	_, ok = r.Lookup("extra")
	assert.False(t, ok, "original registry should not have 'extra'")
	_, ok = cloned.Lookup("extra")
	assert.True(t, ok, "cloned registry should have 'extra'")
}

func TestRegistry_UnknownFunctionError(t *testing.T) {
	r := NewRegistry()
	parser := New()
	f, err := parser.ParseString("test", "[nonexistent()]")
	require.NoError(t, err)
	_, err = f.Evaluate(context.Background(), nil, r)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown function")
}

func TestDefaultRegistry_AllFunctions(t *testing.T) {
	r := DefaultRegistry()
	for _, name := range []string{"if", "equals", "parameters", "format", "replace", "tolower", "concat", "empty"} {
		_, ok := r.Lookup(name)
		assert.True(t, ok, "DefaultRegistry should contain %q", name)
	}
}
