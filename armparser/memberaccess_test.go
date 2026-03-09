package armparser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessMember(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		key      string
		expected any
		wantErr  bool
	}{
		{
			name:     "direct key lookup",
			value:    map[string]any{"foo": "bar"},
			key:      "foo",
			expected: "bar",
		},
		{
			name:     "case insensitive fallback",
			value:    map[string]any{"Foo": "bar"},
			key:      "foo",
			expected: "bar",
		},
		{
			name:     "case insensitive fallback uppercase key",
			value:    map[string]any{"foo": "bar"},
			key:      "FOO",
			expected: "bar",
		},
		{
			name:     "missing member returns nil",
			value:    map[string]any{"foo": "bar"},
			key:      "baz",
			expected: nil,
		},
		{
			name:     "nil value returns nil",
			value:    nil,
			key:      "foo",
			expected: nil,
		},
		{
			name:    "non-map type returns error",
			value:   "a string",
			key:     "foo",
			wantErr: true,
		},
		{
			name:    "int type returns error",
			value:   42,
			key:     "foo",
			wantErr: true,
		},
		{
			name: "nested map access",
			value: map[string]any{
				"inner": map[string]any{"deep": "value"},
			},
			key:      "inner",
			expected: map[string]any{"deep": "value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := accessMember(tt.value, tt.key)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResolveMemberAccess(t *testing.T) {
	ctx := context.Background()
	registry := DefaultRegistry()
	evalCtx := FromMap(map[string]any{})

	tests := []struct {
		name           string
		value          any
		dotMembers     []string
		bracketMembers []*StringExpression
		expected       any
		wantErr        bool
	}{
		{
			name:     "no members returns value unchanged",
			value:    "hello",
			expected: "hello",
		},
		{
			name:       "single dot member",
			value:      map[string]any{"foo": "bar"},
			dotMembers: []string{"foo"},
			expected:   "bar",
		},
		{
			name:       "nested dot members",
			value:      map[string]any{"a": map[string]any{"b": map[string]any{"c": "deep"}}},
			dotMembers: []string{"a", "b", "c"},
			expected:   "deep",
		},
		{
			name:  "single bracket member",
			value: map[string]any{"foo": "bar"},
			bracketMembers: []*StringExpression{
				{String: strPtr("foo")},
			},
			expected: "bar",
		},
		{
			name:  "mixed bracket then dot",
			value: map[string]any{"a": map[string]any{"b": "val"}},
			bracketMembers: []*StringExpression{
				{String: strPtr("a")},
			},
			dotMembers: []string{"b"},
			expected:   "val",
		},
		{
			name:       "dot on non-map errors",
			value:      "not a map",
			dotMembers: []string{"foo"},
			wantErr:    true,
		},
		{
			name:       "missing nested dot returns nil",
			value:      map[string]any{"a": map[string]any{}},
			dotMembers: []string{"a", "missing"},
			expected:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveMemberAccess(tt.value, tt.dotMembers, tt.bracketMembers, ctx, evalCtx, registry)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func strPtr(s string) *string { return &s }
