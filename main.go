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
	example := `[parameters('testobject').key1.key2]`
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

	val, err := ast.Evaluate(armparser.EvalContext{
		"testobject": map[string]any{
			"key1": map[string]any{
				"key2": "value2",
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Evaluated: %v\n", val)
}
