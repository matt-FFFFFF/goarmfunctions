package armparser

import (
	"context"
	"strings"
)

// ArmFunc is the signature for ARM template functions.
type ArmFunc func(ctx context.Context, call *FunctionCall, evalCtx EvalContext) (any, error)

// FuncRegistry holds a set of named ARM template functions.
type FuncRegistry struct {
	funcs map[string]ArmFunc
}

type registryContextKey struct{}

// NewRegistry creates an empty FuncRegistry.
func NewRegistry() *FuncRegistry {
	return &FuncRegistry{funcs: make(map[string]ArmFunc)}
}

// Register adds a function to the registry under its lowercase name.
func (r *FuncRegistry) Register(name string, fn ArmFunc) {
	r.funcs[strings.ToLower(name)] = fn
}

// Lookup retrieves a function by name (case-insensitive).
func (r *FuncRegistry) Lookup(name string) (ArmFunc, bool) {
	fn, ok := r.funcs[strings.ToLower(name)]
	return fn, ok
}

// Clone returns a deep copy of the registry for consumer extension.
func (r *FuncRegistry) Clone() *FuncRegistry {
	c := NewRegistry()
	for k, v := range r.funcs {
		c.funcs[k] = v
	}
	return c
}

// ContextWithRegistry stores the registry in a context for use by function implementations.
func ContextWithRegistry(ctx context.Context, r *FuncRegistry) context.Context {
	return context.WithValue(ctx, registryContextKey{}, r)
}

// RegistryFromContext retrieves the registry from context.
func RegistryFromContext(ctx context.Context) *FuncRegistry {
	if r, ok := ctx.Value(registryContextKey{}).(*FuncRegistry); ok {
		return r
	}
	return nil
}

// DefaultRegistry returns a FuncRegistry pre-loaded with all built-in ARM functions.
func DefaultRegistry() *FuncRegistry {
	r := NewRegistry()
	r.Register("if", If)
	r.Register("equals", Equals)
	r.Register("parameters", Parameters)
	r.Register("format", Format)
	r.Register("replace", Replace)
	r.Register("toLower", ToLower)
	r.Register("concat", Concat)
	r.Register("empty", Empty)
	r.Register("contains", Contains)
	r.Register("length", Length)
	r.Register("and", And)
	r.Register("or", Or)
	r.Register("not", Not)
	r.Register("greater", Greater)
	r.Register("less", Less)
	r.Register("greaterOrEquals", GreaterOrEquals)
	r.Register("lessOrEquals", LessOrEquals)
	r.Register("startsWith", StartsWith)
	r.Register("endsWith", EndsWith)
	r.Register("split", Split)
	r.Register("trim", Trim)
	r.Register("int", Int)
	r.Register("string", StringFunc)
	r.Register("bool", Bool)
	r.Register("utcNow", UtcNow)
	return r
}
