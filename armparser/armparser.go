package armparser

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/matt-FFFFFF/goarmfunctions/armlexer"
	"github.com/matt-FFFFFF/goarmfunctions/logger"
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

func (a *ArmValue) Evaluate(ctx context.Context, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("ArmValue.Evaluate")
	defer lgr.Debug("ArmValue.Evaluate done")
	return a.ArmTemplateString.Evaluate(ctx, evalCtx)
}

func (t *ArmTemplateString) Evaluate(ctx context.Context, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("ArmTemplateString.Evaluate")
	defer lgr.Debug("ArmTemplateString.Evaluate done")
	var result strings.Builder
	for _, part := range t.Parts {
		if part.Literal != nil {
			result.WriteString(*part.Literal)
		}
		if part.Expression != nil {
			value, err := part.Expression.Evaluate(ctx, evalCtx)
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
func (f *FunctionCall) Evaluate(ctx context.Context, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("FunctionCall.Evaluate", slog.String("identifier", f.Name))
	defer lgr.Debug("FunctionCall.Evaluate done")
	switch f.Name {
	case "if":
		return If(ctx, f, evalCtx)
	case "equals":
		return Equals(ctx, f, evalCtx)
	case "parameters":
		return Parameters(ctx, f, evalCtx)
	case "format":
		return Format(ctx, f, evalCtx)
	case "replace":
		return Replace(ctx, f, evalCtx)
	case "toLower":
		return ToLower(ctx, f, evalCtx)
	case "concat":
		return Concat(ctx, f, evalCtx)
	}
	lgr.Error("unknown function", slog.String("function", f.Name))
	return nil, fmt.Errorf("unknown function: %s", f.Name)
}

func (e *Expression) Evaluate(ctx context.Context, evalCtx EvalContext) (any, error) {
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
		return e.FunctionCall.Evaluate(ctx, evalCtx)
	}
	lgr.Error(
		"unsupported expression type",
		slog.String("expression", fmt.Sprintf("%v", e)),
	)
	return nil, fmt.Errorf("unsupported expression type")
}

func (e *StringExpression) Evaluate(ctx context.Context, evalCtx EvalContext) (any, error) {
	lgr := logger.LoggerFromContext(ctx)
	lgr.Debug("StringExpression.Evaluate",
		slog.String("string", fmt.Sprintf("%v", e.String)),
	)
	defer lgr.Debug("StringExpression.Evaluate done")
	if e.String != nil {
		return *e.String, nil
	}
	if e.FunctionCall != nil {
		return e.FunctionCall.Evaluate(ctx, evalCtx)
	}
	lgr.Error(
		"unsupported expression type",
		slog.String("expression", fmt.Sprintf("%v", e)),
	)
	return nil, fmt.Errorf("unsupported expression type")
}
