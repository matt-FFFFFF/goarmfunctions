package armparser

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/matt-FFFFFF/goarmfunctions/armlexer"
)

type EvalContext map[string]interface{}

// ArmFunction is the root node of the ARM function AST.
// It contains an expression node that identifies ARM functions as being enclosed in square brackets.
type ArmFunction struct {
	Expression *Expression `"[" @@ "]"`
}

// Expression is a node in the ARM function AST.
// It can be a string, number, boolean, or function call.
type Expression struct {
	String       *string           `@String`
	Number       *int              `| @Number`
	Boolean      *armlexer.Boolean `| @Boolean`
	FunctionCall *FunctionCall     `| @@`
}

// FunctionCall is a node in the ARM function AST.
// It represents a function call with a name and optional arguments.
type FunctionCall struct {
	Name    string        `@Ident`
	Args    []*Expression `"(" ( @@ ( "," @@ )* )? ")"`
	Members []string      `( "." @Ident )*` // Capture additional members in a slice
}

// New returns a new ARM function parser.
func New() *participle.Parser[ArmFunction] {
	return participle.MustBuild[ArmFunction](
		participle.Lexer(armlexer.New()),
		participle.Unquote("String"),
		participle.CaseInsensitive("Ident"),
		participle.Elide("Whitespace"),
	)
}

// Evaluate evaluates the ARM function AST.
func (f *FunctionCall) Evaluate(ctx EvalContext) (interface{}, error) {
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

func (e *Expression) Evaluate(ctx EvalContext) (interface{}, error) {
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
