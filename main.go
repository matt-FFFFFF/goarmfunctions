package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "true"
	return nil
}

var lexerDef = lexer.Must(lexer.New(lexer.Rules{
	"Root": {
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
			Action:  nil,
		},
		{
			Name:    "Punct",
			Pattern: `[,()\[\].]`,
		},
	}}))

type ArmFunction struct {
	Expression *Expression `"[" @@ "]"`
}

type Expression struct {
	String       *string       `@String`
	Number       *int          `| @Number`
	Boolean      *Boolean      `| @Boolean`
	FunctionCall *FunctionCall `| @@`
}

type FunctionCall struct {
	Name string        `@Ident`
	Args []*Expression `"(" ( @@ ( "," @@ )* )? ")"`

	//Members []string `( "." @Ident )*` // Capture additional members in a slice
}

var parser = participle.MustBuild[ArmFunction](
	participle.Lexer(lexerDef),
	participle.Unquote("String"),
	participle.CaseInsensitive("Ident"),
)

func main() {
	//example := `[if(equals(parameters('dnsZoneSubscriptionId'), ''), parameters('azureAutomationWebhookPrivateDnsZoneId'), format('/subscriptions/{0}/resourceGroups/{1}/providers/{2}/{3}', parameters('dnsZoneSubscriptionId'), toLower(parameters('dnsZoneResourceGroupName')), parameters('dnsZoneResourceType'), replace(replace(parameters('dnsZoneNames').azureAutomationWebhookPrivateDnsZoneId, '{regionName}', parameters('dnsZoneRegion')), '{regionCode}', parameters('dnzZoneRegionShortNames')[parameters('dnsZoneRegion')])))]`
	//example := `[if(equals(parameters('dnsZoneSubscriptionId'), ''), parameters('test1'), parameters('test2'))]`
	example := `[if(true, 'test1', 'test2')]`

	reader := io.Reader(strings.NewReader(example))
	lexed, err := lexerDef.Lex("test", reader)
	for tok, err := lexed.Next(); err == nil && !tok.EOF(); tok, err = lexed.Next() {
		fmt.Printf("Type: %s, Value: %s\n", tokType2Str(tok.Type), tok.Value)
	}

	_ = lexed
	ast, err := parser.ParseString("test", example)
	if err != nil {
		fmt.Println(err.Error())
	}
	if ast.Expression != nil && ast.Expression.FunctionCall != nil {
		fmt.Printf("Function Name: %s\n", ast.Expression.FunctionCall.Name)
		for i, arg := range ast.Expression.FunctionCall.Args {
			fmt.Printf("Arg %d: %+v\n", i, arg)
		}
	}
	_ = ast
}

// Helper function to get the token type name
func tokType2Str(t lexer.TokenType) string {
	for k, v := range lexerDef.Symbols() {
		if v == t {
			return k
		}
	}
	return ""
}
