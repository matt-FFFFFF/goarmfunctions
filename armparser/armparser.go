package armparser

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/alecthomas/participle/v2"
	"github.com/matt-FFFFFF/goarmfunctions/armlexer"
	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

var (
	parserOnce   sync.Once
	sharedParser *participle.Parser[ArmValue]
)

// New returns a new ARM function parser.
func New() *participle.Parser[ArmValue] {
	return participle.MustBuild[ArmValue](
		participle.Lexer(armlexer.New()),
		participle.Unquote("String"),
		participle.CaseInsensitive("Ident"),
		participle.Elide("Whitespace"),
		participle.UseLookahead(5),
	)
}

// SharedParser returns a shared parser instance.
// It is safe for concurrent use.
func SharedParser() *participle.Parser[ArmValue] {
	parserOnce.Do(func() {
		sharedParser = New()
	})
	return sharedParser
}

func (a *ArmValue) Evaluate(ctx context.Context, evalCtx EvalContext, registry *FuncRegistry) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("ArmValue.Evaluate")
	defer lgr.Debug("ArmValue.Evaluate done")
	return a.ArmTemplateString.Evaluate(ctx, evalCtx, registry)
}

func (t *ArmTemplateString) Evaluate(ctx context.Context, evalCtx EvalContext, registry *FuncRegistry) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("ArmTemplateString.Evaluate")
	defer lgr.Debug("ArmTemplateString.Evaluate done")
	var result strings.Builder
	for _, part := range t.Parts {
		if part.Literal != nil {
			result.WriteString(*part.Literal)
		}
		if part.Expression != nil {
			value, err := part.Expression.Evaluate(ctx, evalCtx, registry)
			if err != nil {
				return "", err
			}
			// If there is only one part, then the return value might not be a string.
			if len(t.Parts) == 1 {
				return value, nil
			}
			result.WriteString(fmt.Sprintf("%v", value))
		}
	}
	return result.String(), nil
}

// Evaluate evaluates the ARM function AST.
func (f *FunctionCall) Evaluate(ctx context.Context, evalCtx EvalContext, registry *FuncRegistry) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("FunctionCall.Evaluate", slog.String("identifier", f.Name))
	defer lgr.Debug("FunctionCall.Evaluate done")

	if registry == nil {
		registry = DefaultRegistry()
	}

	fn, ok := registry.Lookup(f.Name)
	if !ok {
		lgr.Error("unknown function", slog.String("function", f.Name))
		return nil, fmt.Errorf("unknown function: %s", f.Name)
	}

	// Store registry in context so function implementations can access it
	// when they need to evaluate sub-expressions.
	ctx = ContextWithRegistry(ctx, registry)
	result, err := fn(ctx, f, evalCtx)
	if err != nil {
		return nil, err
	}

	// Apply member access (.dot and ['bracket']) on the function's return value.
	return resolveMemberAccess(result, f.MembersDot, f.MembersStr, ctx, evalCtx, registry)
}

func (e *Expression) Evaluate(ctx context.Context, evalCtx EvalContext, registry *FuncRegistry) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("Expression.Evaluate",
		slog.String("string", fmt.Sprintf("%v", e.String)),
		slog.String("number", fmt.Sprintf("%v", e.Number)),
		slog.String("boolean", fmt.Sprintf("%v", e.Boolean)),
	)
	defer lgr.Debug("Expression.Evaluate done")
	if e.String != nil {
		return *e.String, nil
	}
	if e.Number != nil {
		return *e.Number, nil
	}
	if e.Boolean != nil {
		return bool(*e.Boolean), nil
	}
	if e.FunctionCall != nil {
		return e.FunctionCall.Evaluate(ctx, evalCtx, registry)
	}
	lgr.Error(
		"unsupported expression type",
		slog.String("expression", fmt.Sprintf("%v", e)),
	)
	return nil, fmt.Errorf("unsupported expression type")
}

func (e *StringExpression) Evaluate(ctx context.Context, evalCtx EvalContext, registry *FuncRegistry) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("StringExpression.Evaluate",
		slog.String("string", fmt.Sprintf("%v", e.String)),
	)
	defer lgr.Debug("StringExpression.Evaluate done")
	if e.String != nil {
		return *e.String, nil
	}
	if e.FunctionCall != nil {
		return e.FunctionCall.Evaluate(ctx, evalCtx, registry)
	}
	lgr.Error(
		"unsupported expression type",
		slog.String("expression", fmt.Sprintf("%v", e)),
	)
	return nil, fmt.Errorf("unsupported expression type")
}
