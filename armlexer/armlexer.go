package armlexer

import "github.com/alecthomas/participle/v2/lexer"

// Boolean is a custom type for parsing boolean values represented as `true|false`.
type Boolean bool

// Capture implements the participle.Capture interface for the Boolean type.
func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "true"
	return nil
}

// New returns a new lexer definition for ARM functions.
func New() lexer.Definition {
	return lexer.MustSimple([]lexer.SimpleRule{
		{
			Name:    "Boolean",
			Pattern: `true|false`,
		},
		{
			Name:    "Ident",
			Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`,
		},
		{
			Name:    "String",
			Pattern: `'(?:\\.|[^'])*'`,
		},
		{
			Name:    "Number",
			Pattern: `[0-9]+`,
		},
		{
			Name:    "Whitespace",
			Pattern: `[ \t\n\r]+`,
		},
		{
			Name:    "Punct",
			Pattern: `[,()\[\].]`,
		},
	})
}

// TokenType2Str is a helper function to get the token type names from the lexer
func TokenType2Str(in map[string]lexer.TokenType) map[lexer.TokenType]string {
	res := make(map[lexer.TokenType]string)
	for k, v := range in {
		res[v] = k
	}
	return res
}