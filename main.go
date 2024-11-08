package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var lexerRules = lexer.MustSimple([]lexer.SimpleRule{
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
		Name:    "Punct",
		Pattern: `[,(){}]`,
	},
	{
		Name:    "SquareBracket",
		Pattern: `[\[\]]`,
	},
	{
		Name:    "Whitespace",
		Pattern: `[ \t\n\r]+`,
	},
})

type Root struct {
	Expression *Expression `"[" @@ "]"`
}

type Expression struct {
	FunctionCall *FunctionCall `@@`
	Identifier   *Identifier   `| @Ident ( "." @Ident )*`
	String       *string       `| @String`
	Number       *int          `| @Number`
}

type FunctionCall struct {
	Name string      `@Ident`
	Args *Expression `"(" (@@ ( "," @@ )* )? ")"`
}

type Identifier struct {
	Base    string   `@Ident`
	Members []string `("." @Ident)*` // Capture additional members in a slice
}

var parser = participle.MustBuild[*Root](
	participle.Lexer(lexerRules),
	participle.Unquote("String"),
	participle.CaseInsensitive("Ident"),
)

func main() {
	example := `[if(equals(parameters('dnsZoneSubscriptionId'), ''), parameters('azureAutomationWebhookPrivateDnsZoneId'), format('/subscriptions/{0}/resourceGroups/{1}/providers/{2}/{3}', parameters('dnsZoneSubscriptionId'), toLower(parameters('dnsZoneResourceGroupName')), parameters('dnsZoneResourceType'), replace(replace(parameters('dnsZoneNames').azureAutomationWebhookPrivateDnsZoneId, '{regionName}', parameters('dnsZoneRegion')), '{regionCode}', parameters('dnzZoneRegionShortNames')[parameters('dnsZoneRegion')])))]`
	ast, err := parser.ParseString("test", example)
	if err != nil {
		panic(err)
	}
	_ = ast
}
