package armparser

import "github.com/matt-FFFFFF/goarmfunctions/armlexer"

// EvalContext is used to pass parameter names and values to the ARM function evaluator.
type EvalContext map[string]any

// ArmValue is the root node of the ARM function AST.
// It represents any JSON value in an ARM resource.
type ArmValue struct {
	ArmTemplateString *ArmTemplateString `@@`
}

// ArmTemplateString is a node in the ARM function AST.
// It represents a string that may contain either an ARM function or a literal.
type ArmTemplateString struct {
	Parts []*ArmTemplatePart `( @@ )+`
}

// ArmTemplatePart is a node in the ARM function AST.
// It represents a part of an ARM template string, which can be either a literal or an expression.
type ArmTemplatePart struct {
	Literal    *string     `@UnquotedLiteral`
	Expression *Expression `| "[" @@ "]"`
}

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
	Name       string              `@Ident`
	Args       []*Expression       `"(" ( @@ ( "," @@ )* )? ")"`
	MembersStr []*StringExpression `  ( "[" @@ "]" )*` // Capture additional member access by `['foo']` in a slice
	MembersDot []string            `  ( "." @Ident )*` // Capture additional members in a slice
}

type StringExpression struct {
	String       *string       `@String`
	FunctionCall *FunctionCall `| @@`
}
