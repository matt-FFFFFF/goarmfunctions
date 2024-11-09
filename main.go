package main

import (
	"bytes"
	"fmt"

	"github.com/matt-FFFFFF/goarmfunctions/armlexer"
	"github.com/matt-FFFFFF/goarmfunctions/armparser"
)

func main() {
	armlex := armlexer.New()
	//example := `[if(equals(parameters('dnsZoneSubscriptionId'), ''), parameters('azureAutomationWebhookPrivateDnsZoneId'), format('/subscriptions/{0}/resourceGroups/{1}/providers/{2}/{3}', parameters('dnsZoneSubscriptionId'), toLower(parameters('dnsZoneResourceGroupName')), parameters('dnsZoneResourceType'), replace(replace(parameters('dnsZoneNames').azureAutomationWebhookPrivateDnsZoneId, '{regionName}', parameters('dnsZoneRegion')), '{regionCode}', parameters('dnzZoneRegionShortNames')[parameters('dnsZoneRegion')])))]`
	example := `[if(equals(parameters('dnsZoneSubscriptionId'), ''), parameters('test1'), parameters('test2'))]`
	//example := `[if(true, 'test1', 'test2')]`

	reader := bytes.NewReader([]byte(example))
	lexed, err := armlex.Lex("test", reader)
	t2s := armlexer.TokenType2Str(armlex.Symbols())
	for tok, err := lexed.Next(); err == nil && !tok.EOF(); tok, err = lexed.Next() {
		fmt.Printf("Type: %s, Value: %s\n", t2s[tok.Type], tok.Value)
	}
	parser := armparser.New()
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
	val, err := ast.Expression.Evaluate(armparser.EvalContext{
		"dnsZoneSubscriptionId": "test",
		"test1":                 "test1",
		"test2":                 "test2",
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Evaluated: %v\n", val)
}
