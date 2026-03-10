package armparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapScope_GetLocal(t *testing.T) {
	s := NewScope("test", map[string]any{"a": 1, "b": "two"}, nil)
	v, ok := s.GetLocal("a")
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	_, ok = s.GetLocal("missing")
	assert.False(t, ok)
}

func TestMapScope_Get_WalksParent(t *testing.T) {
	parent := NewScope("parent", map[string]any{"x": 10}, nil)
	child := NewScope("child", map[string]any{"y": 20}, parent)

	// Child can see its own values
	v, ok := child.Get("y")
	assert.True(t, ok)
	assert.Equal(t, 20, v)

	// Child can see parent values
	v, ok = child.Get("x")
	assert.True(t, ok)
	assert.Equal(t, 10, v)

	// Parent cannot see child values
	_, ok = parent.Get("y")
	assert.False(t, ok)
}

func TestMapScope_Get_ChildShadowsParent(t *testing.T) {
	parent := NewScope("parent", map[string]any{"k": "parent-val"}, nil)
	child := NewScope("child", map[string]any{"k": "child-val"}, parent)

	v, ok := child.Get("k")
	assert.True(t, ok)
	assert.Equal(t, "child-val", v)
}

func TestMapScope_Child(t *testing.T) {
	root := NewScope("root", map[string]any{"a": 1}, nil)
	child := root.Child("level1", map[string]any{"b": 2})

	assert.Equal(t, "level1", child.Scope())
	v, ok := child.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	v, ok = child.Get("b")
	assert.True(t, ok)
	assert.Equal(t, 2, v)
}

func TestMapScope_Parent(t *testing.T) {
	root := NewScope("root", nil, nil)
	assert.Nil(t, root.Parent())

	child := NewScope("child", nil, root)
	assert.Equal(t, root, child.Parent())
}

func TestFindScope(t *testing.T) {
	request := NewScope("request", map[string]any{"apiVersion": "2021"}, nil)
	params := NewScope("parameters", map[string]any{"foo": "bar"}, request)
	resource := NewScope("resource", map[string]any{"type": "vm"}, params)

	found := FindScope(resource, "parameters")
	require.NotNil(t, found)
	assert.Equal(t, "parameters", found.Scope())
	v, ok := found.GetLocal("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", v)

	found = FindScope(resource, "request")
	require.NotNil(t, found)
	assert.Equal(t, "request", found.Scope())

	found = FindScope(resource, "nonexistent")
	assert.Nil(t, found)
}

func TestFindScope_Nil(t *testing.T) {
	assert.Nil(t, FindScope(nil, "anything"))
}

func TestFromMap(t *testing.T) {
	ctx := FromMap(map[string]any{"key": "value"})
	assert.Equal(t, "parameters", ctx.Scope())
	v, ok := ctx.GetLocal("key")
	assert.True(t, ok)
	assert.Equal(t, "value", v)
	assert.Nil(t, ctx.Parent())
}

func TestFromMap_Nil(t *testing.T) {
	ctx := FromMap(nil)
	assert.Equal(t, "parameters", ctx.Scope())
	_, ok := ctx.GetLocal("anything")
	assert.False(t, ok)
}

func TestDeepScopeChain(t *testing.T) {
	s1 := NewScope("s1", map[string]any{"a": 1}, nil)
	s2 := s1.Child("s2", map[string]any{"b": 2})
	s3 := s2.Child("s3", map[string]any{"c": 3})
	s4 := s3.Child("s4", map[string]any{"d": 4})

	// s4 can reach all the way to s1
	v, ok := s4.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	// FindScope from s4 finds s2
	found := FindScope(s4, "s2")
	require.NotNil(t, found)
	assert.Equal(t, "s2", found.Scope())
}

func TestNewScope_NilValues(t *testing.T) {
	s := NewScope("test", nil, nil)
	_, ok := s.GetLocal("anything")
	assert.False(t, ok)
}
