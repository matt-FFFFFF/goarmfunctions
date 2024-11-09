package armparser

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/matt-FFFFFF/goarmfunctions/armlexer"
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

func (a *ArmValue) Evaluate(ctx EvalContext) (string, error) {
	return a.ArmTemplateString.Evaluate(ctx)
}

func (t *ArmTemplateString) Evaluate(ctx EvalContext) (string, error) {
	var result strings.Builder
	for _, part := range t.Parts {
		if part.Literal != nil {
			result.WriteString(*part.Literal)
		}
		if part.Expression != nil {
			value, err := part.Expression.Evaluate(ctx)
			if err != nil {
				return "", err
			}
			result.WriteString(fmt.Sprintf("%v", value))
		}
	}
	return result.String(), nil
}

// Evaluate evaluates the ARM function AST.
func (f *FunctionCall) Evaluate(ctx EvalContext) (any, error) {
	switch f.Name {
	case "if":
		return If(f, ctx)
	case "equals":
		return Equals(f, ctx)
	case "parameters":
		return Parameters(f, ctx)
	}
	return nil, fmt.Errorf("unknown function: %s", f.Name)
}

func (e *Expression) Evaluate(ctx EvalContext) (any, error) {
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
		return e.FunctionCall.Evaluate(ctx)
	}
	return nil, fmt.Errorf("unsupported expression type")
}
