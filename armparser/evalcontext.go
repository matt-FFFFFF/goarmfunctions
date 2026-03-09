package armparser

// EvalContext is the interface for hierarchical evaluation contexts.
// Implementations form a linked chain of named scopes.
type EvalContext interface {
	// Get retrieves a value by key from this scope or any parent scope.
	Get(key string) (any, bool)

	// GetLocal retrieves a value from THIS scope only (no parent walk).
	GetLocal(key string) (any, bool)

	// Parent returns the parent scope, or nil if this is the root.
	Parent() EvalContext

	// Scope returns the scope name (e.g., "parameters", "subscription", "resource").
	Scope() string

	// Child creates a new child scope with the given name and initial values.
	Child(name string, values map[string]any) EvalContext
}

// MapScope is the standard implementation of EvalContext.
type MapScope struct {
	name   string
	values map[string]any
	parent EvalContext
}

// NewScope creates a new MapScope with the given name, values, and optional parent.
func NewScope(name string, values map[string]any, parent EvalContext) *MapScope {
	if values == nil {
		values = make(map[string]any)
	}
	return &MapScope{name: name, values: values, parent: parent}
}

func (s *MapScope) Get(key string) (any, bool) {
	if v, ok := s.values[key]; ok {
		return v, true
	}
	if s.parent != nil {
		return s.parent.Get(key)
	}
	return nil, false
}

func (s *MapScope) GetLocal(key string) (any, bool) {
	v, ok := s.values[key]
	return v, ok
}

func (s *MapScope) Parent() EvalContext {
	return s.parent
}

func (s *MapScope) Scope() string {
	return s.name
}

func (s *MapScope) Child(name string, values map[string]any) EvalContext {
	return NewScope(name, values, s)
}

// FindScope walks up the scope chain to find a scope with the given name.
// Returns nil if no matching scope is found.
func FindScope(ctx EvalContext, name string) EvalContext {
	for c := ctx; c != nil; c = c.Parent() {
		if c.Scope() == name {
			return c
		}
	}
	return nil
}

// FromMap wraps a legacy map[string]any in a "parameters" scope for backwards compatibility.
// The map values are placed directly in the scope so that GetLocal(paramName) works.
func FromMap(m map[string]any) EvalContext {
	return NewScope("parameters", m, nil)
}
