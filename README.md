# goarmfunctions

Go library for parsing and evaluating ARM template expressions such as `[concat('a', parameters('b'))]`.

## Installation

```bash
go get github.com/matt-FFFFFF/goarmfunctions
```

## Quick Start

The simplest way to evaluate an expression — pass a `map[string]any` of parameters:

```go
package main

import (
	"context"
	"fmt"

	goarmfunctions "github.com/matt-FFFFFF/goarmfunctions"
)

func main() {
	params := map[string]any{
		"env": "prod",
	}
	result, err := goarmfunctions.LexAndParse(context.Background(), "[concat('hello-', parameters('env'))]", params, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(result) // hello-prod
}
```

## Advanced Usage

For full control, use `Evaluate()` with an `EvalContext` and `FuncRegistry`:

```go
package main

import (
	"context"
	"fmt"

	goarmfunctions "github.com/matt-FFFFFF/goarmfunctions"
	"github.com/matt-FFFFFF/goarmfunctions/armparser"
)

func main() {
	// Build a tiered evaluation context
	params := armparser.NewScope("parameters", map[string]any{
		"location": "uksouth",
	}, nil)
	resource := params.Child("resource", map[string]any{
		"type": "Microsoft.Storage/storageAccounts",
	})

	// Clone the default registry and add custom functions
	registry := armparser.DefaultRegistry().Clone()

	result, err := goarmfunctions.Evaluate(
		context.Background(),
		"[concat(parameters('location'), '-storage')]",
		resource, // evaluation walks up the scope chain
		registry,
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(result) // uksouth-storage
}
```

## Tiered Evaluation Context

`EvalContext` is an interface that forms a linked chain of named scopes. When a value is requested, the chain walks upward from the current scope to the root until the key is found.

```go
// Root scope — subscription-level data
sub := armparser.NewScope("subscription", map[string]any{
    "subscriptionId": "00000000-0000-0000-0000-000000000000",
}, nil)

// Child scope — resource group
rg := sub.Child("resourceGroup", map[string]any{
    "name":     "my-rg",
    "location": "uksouth",
})

// Child scope — parameters
params := rg.Child("parameters", map[string]any{
    "env": "prod",
})

// Get walks up: finds "env" in params, "subscriptionId" in sub
v, _ := params.Get("env")             // "prod"
v, _ = params.Get("subscriptionId")   // "00000000-..."

// GetLocal only checks the current scope
_, ok := params.GetLocal("subscriptionId") // false

// FindScope locates a named scope in the chain
rgScope := armparser.FindScope(params, "resourceGroup")

// FromMap wraps a plain map for backwards compatibility
legacy := armparser.FromMap(map[string]any{"key": "value"})
```

This lets consumers build hierarchical contexts for policy engines, deployment engines, or any system with layered configuration.

## Function Registry

`FuncRegistry` holds named functions that the parser invokes during evaluation.

```go
// Get the default registry with all 25 built-in functions
registry := armparser.DefaultRegistry()

// Clone it to add your own without mutating the default
custom := registry.Clone()

// Register a new function
custom.Register("myFunc", func(ctx context.Context, call *armparser.FunctionCall, evalCtx armparser.EvalContext) (any, error) {
    return "custom-value", nil
})

// Lookup a function
fn, ok := custom.Lookup("concat")
```

## Built-in Functions

The `DefaultRegistry()` includes 25 functions:

| Function | Description |
|---|---|
| `if` | Conditional — returns one of two values |
| `equals` | Equality comparison |
| `parameters` | Retrieve a parameter value from context |
| `format` | String formatting |
| `replace` | String replacement |
| `toLower` | Lowercase a string |
| `concat` | Concatenate strings or arrays |
| `empty` | Check if a value is empty |
| `contains` | Check if a string/array contains a value |
| `length` | Length of a string, array, or object |
| `and` | Logical AND |
| `or` | Logical OR |
| `not` | Logical NOT |
| `greater` | Greater-than comparison |
| `less` | Less-than comparison |
| `greaterOrEquals` | Greater-than-or-equal comparison |
| `lessOrEquals` | Less-than-or-equal comparison |
| `startsWith` | Check if a string starts with a value |
| `endsWith` | Check if a string ends with a value |
| `split` | Split a string by delimiter |
| `trim` | Trim whitespace |
| `int` | Convert to integer |
| `string` | Convert to string |
| `bool` | Convert to boolean |
| `utcNow` | Current UTC time as ISO 8601 string |

## Extending with Custom Functions

Register domain-specific functions for your use case. For example, a `field()` function for Azure Policy evaluation:

```go
registry := armparser.DefaultRegistry().Clone()

registry.Register("field", func(ctx context.Context, call *armparser.FunctionCall, evalCtx armparser.EvalContext) (any, error) {
    // Resolve the field name from the first argument
    args, err := call.EvalArgs(ctx, evalCtx)
    if err != nil {
        return nil, err
    }
    fieldName, ok := args[0].(string)
    if !ok {
        return nil, fmt.Errorf("field() expects a string argument")
    }

    // Look up in a "resource" scope
    if scope := armparser.FindScope(evalCtx, "resource"); scope != nil {
        if v, found := scope.GetLocal(fieldName); found {
            return v, nil
        }
    }
    return nil, nil
})

// Now you can evaluate: [equals(field('type'), 'Microsoft.Storage/storageAccounts')]
result, _ := goarmfunctions.Evaluate(ctx, expr, evalCtx, registry, nil)
```

## Member Access

All functions automatically support `.dot` and `['bracket']` member access on return values:

```
[parameters('storage').name]
[parameters('tags')['environment']]
[concat(parameters('prefix'), '-vm').name]
```

No special handling is needed — the parser resolves member access chains on any value returned by a function.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a PR

## License

MIT — see [LICENSE](LICENSE) for details.
