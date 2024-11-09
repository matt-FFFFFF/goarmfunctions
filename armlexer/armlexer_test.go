package armlexer

import (
	"testing"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/assert"
)

func TestBooleanCapture(t *testing.T) {
	tests := []struct {
		input    []string
		expected Boolean
	}{
		{input: []string{"true"}, expected: true},
		{input: []string{"false"}, expected: false},
		{input: []string{"TRUE"}, expected: false},  // invalid input, should default to false
		{input: []string{"FALSE"}, expected: false}, // invalid input, should default to false
	}

	for _, test := range tests {
		var b Boolean
		err := b.Capture(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, b)
	}
}

func TestTokenType2Str(t *testing.T) {
	input := map[string]lexer.TokenType{
		"UnquotedLiteral": 1,
		"Boolean":         2,
		"Ident":           3,
		"String":          4,
		"Number":          5,
		"Punct":           6,
		"Whitespace":      7,
		"ArmFunctionEnd":  8,
	}

	expected := map[lexer.TokenType]string{
		1: "UnquotedLiteral",
		2: "Boolean",
		3: "Ident",
		4: "String",
		5: "Number",
		6: "Punct",
		7: "Whitespace",
		8: "ArmFunctionEnd",
	}

	result := TokenType2Str(input)
	assert.Equal(t, expected, result)
}
